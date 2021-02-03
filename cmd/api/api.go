package api

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"

	swagger "github.com/gradusp/crispy/api"
	"github.com/gradusp/crispy/assets"
	"github.com/gradusp/crispy/internal/audit"
	"github.com/gradusp/crispy/internal/cluster"
	"github.com/gradusp/crispy/internal/healthcheck"
	"github.com/gradusp/crispy/internal/real"
	"github.com/gradusp/crispy/internal/service"
	"github.com/gradusp/crispy/internal/zone"

	zhttp "github.com/gradusp/crispy/internal/zone/delivery/rest"
	zpg "github.com/gradusp/crispy/internal/zone/repository/pgsql"
	zuc "github.com/gradusp/crispy/internal/zone/usecase"

	chttp "github.com/gradusp/crispy/internal/cluster/delivery/rest"
	cpg "github.com/gradusp/crispy/internal/cluster/repository/pgsql"
	cuc "github.com/gradusp/crispy/internal/cluster/usecase"

	srest "github.com/gradusp/crispy/internal/service/delivery/rest"
	spg "github.com/gradusp/crispy/internal/service/repository/pgsql"
	suc "github.com/gradusp/crispy/internal/service/usecase"

	rpg "github.com/gradusp/crispy/internal/real/repository/pgsql"
	ruc "github.com/gradusp/crispy/internal/real/usecase"

	hcpg "github.com/gradusp/crispy/internal/healthcheck/repository/pgsql"
	hcuc "github.com/gradusp/crispy/internal/healthcheck/usecase"

	apg "github.com/gradusp/crispy/internal/audit/repository/pgsql"
	auc "github.com/gradusp/crispy/internal/audit/usecase"
)

// TODO: https://github.com/caarlos0/env
// TODO: https://pkg.go.dev/go.uber.org/zap

type App struct {
	httpServer *http.Server
	logger     *zap.Logger

	zoneUC        zone.Usecase
	clusterUC     cluster.Usecase
	healthcheckUC healthcheck.Usecase
	realUC        real.Usecase
	serviceUC     service.Usecase
	auditUC       audit.Usecase
}

func NewApp() *App {
	var logger *zap.Logger
	if gin.Mode() == gin.ReleaseMode {
		l, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		logger = l
	} else {
		l, err := zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
		logger = l
	}

	pool := initPGX()
	kv := initConsul()

	zoneRepo := zpg.NewZonePostgresRepo(pool, logger.Sugar())
	clusterRepo := cpg.NewClusterRepo(pool, kv, logger.Sugar())
	serviceRepo := spg.NewServiceRepo(pool, logger.Sugar())
	realRepo := rpg.NewRealPostgresRepo(pool, logger.Sugar())
	healthcheckRepo := hcpg.NewHealthcheckPostgresRepo(pool, logger.Sugar())
	auditRepo := apg.NewAuditRepo(pool, logger.Sugar())

	return &App{
		logger:        logger,
		clusterUC:     cuc.NewClusterUsecase(clusterRepo),
		serviceUC:     suc.NewServiceUsecase(serviceRepo),
		realUC:        ruc.NewRealUsecase(realRepo),
		healthcheckUC: hcuc.NewHealthcheckUsecase(healthcheckRepo),
		zoneUC:        zuc.NewZoneUseCase(zoneRepo),
		auditUC:       auc.NewAuditUsecase(auditRepo),
	}
}

func (a *App) Run(port string) error {
	// Init gin handler
	router := gin.New()
	router.Use(ginzap.Ginzap(a.logger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(a.logger, true))

	// Set up gin CORS official middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost"} // FIXME: eliminate hardcode
	config.AllowMethods = []string{"GET", "HEAD", "OPTIONS", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"LBOS_API_KEY", "Access-Control-Allow-Headers", "Origin", "Accept", "X-Requested-With", "Content-Type", "Access-Control-Request-Method", "Access-Control-Request-Headers"} //nolint:lll
	config.AllowCredentials = true
	router.Use(cors.New(config))

	// embedding of static assets of SwaggerUI
	f, _ := fs.Sub(assets.UI, "ui")
	router.StaticFS("/swagger", http.FS(f))
	router.GET("/api/v1/openapi.yml", func(c *gin.Context) {
		file, _ := swagger.OpenAPI.ReadFile("openapi.yml")
		c.Data(http.StatusOK, "application/yaml", file)
	})
	// TODO: favicon path
	// https://github.com/gin-gonic/examples/blob/master/assets-in-binary/example02/main.go#L34
	// TODO: gzip
	// https://github.com/gin-contrib/gzip

	// API Endpoints
	rapi := router.Group("/api/v1")
	rapi.Use(zhttp.AuthAPIKey("CRISPY_API_KEY")) // FIXME: current CRISPY_API_KEY flow needs refactor

	zhttp.RegisterHTTPEndpoint(rapi, a.zoneUC, a.auditUC)
	chttp.RegisterHTTPEndpoint(rapi, a.clusterUC)
	srest.RegisterHTTPEndpoint(rapi, a.healthcheckUC, a.realUC, a.serviceUC)

	// HTTP Server
	a.httpServer = &http.Server{
		Addr:           ":" + os.Getenv("CRISPY_PORT"),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second) //nolint:gomnd FIXME
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}

func initPGX() *pgxpool.Pool {
	s := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("CRISPY_DB_USER"),
		os.Getenv("CRISPY_DB_PASS"),
		os.Getenv("CRISPY_DB_HOST"),
		os.Getenv("CRISPY_DB_PORT"),
		os.Getenv("CRISPY_DB_NAME"),
	)

	pool, err := pgxpool.Connect(context.Background(), s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return pool
}

func initConsul() *api.KV {
	cfg := &api.Config{
		Address: os.Getenv("CRISPY_CONSUL_ADDR"),
	}
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	return client.KV()
}
