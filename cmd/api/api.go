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
	"github.com/gradusp/crispy/internal/node"
	"github.com/gradusp/crispy/internal/real"
	"github.com/gradusp/crispy/internal/service"
	"github.com/gradusp/crispy/internal/zone"

	zrest "github.com/gradusp/crispy/internal/zone/delivery/rest"
	zpg "github.com/gradusp/crispy/internal/zone/repository/pgsql"
	zuc "github.com/gradusp/crispy/internal/zone/usecase"

	crest "github.com/gradusp/crispy/internal/cluster/delivery/rest"
	cpg "github.com/gradusp/crispy/internal/cluster/repository/pgsql"
	cuc "github.com/gradusp/crispy/internal/cluster/usecase"

	srest "github.com/gradusp/crispy/internal/service/delivery/rest"
	spg "github.com/gradusp/crispy/internal/service/repository/pgsql"
	suc "github.com/gradusp/crispy/internal/service/usecase"

	hcpg "github.com/gradusp/crispy/internal/healthcheck/repository/pgsql"
	hcuc "github.com/gradusp/crispy/internal/healthcheck/usecase"

	apg "github.com/gradusp/crispy/internal/audit/repository/pgsql"
	auc "github.com/gradusp/crispy/internal/audit/usecase"

	rrest "github.com/gradusp/crispy/internal/real/delivery/rest"
	rpg "github.com/gradusp/crispy/internal/real/repository/pgsql"
	ruc "github.com/gradusp/crispy/internal/real/usecase"

	nrest "github.com/gradusp/crispy/internal/node/delivery/rest"
	npg "github.com/gradusp/crispy/internal/node/repository/pgsql"
	nuc "github.com/gradusp/crispy/internal/node/usecase"
)

// TODO: https://github.com/caarlos0/env
// TODO: https://pkg.go.dev/go.uber.org/zap

type App struct {
	httpServer *http.Server
	logger     *zap.Logger

	zoneUC        zone.Usecase
	clusterUC     cluster.Usecase
	nodeUC        node.Usecase
	serviceUC     service.Usecase
	realUC        real.Usecase
	healthcheckUC healthcheck.Usecase
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

	zoneRepo := zpg.NewPgRepo(pool, logger.Sugar())
	clusterRepo := cpg.NewPgRepo(pool, kv, logger.Sugar())
	nodeRepo := npg.NewPgRepo(pool, logger.Sugar())
	serviceRepo := spg.NewPgRepo(pool, logger.Sugar())
	realRepo := rpg.NewPgRepo(pool, logger.Sugar())
	healthcheckRepo := hcpg.NewPgRepo(pool, logger.Sugar())
	auditRepo := apg.NewPgRepo(pool, logger.Sugar())

	return &App{
		logger:        logger,
		zoneUC:        zuc.NewUsecase(zoneRepo),
		clusterUC:     cuc.NewUsecase(clusterRepo),
		nodeUC:        nuc.NewUsecase(nodeRepo),
		serviceUC:     suc.NewUsecase(serviceRepo),
		realUC:        ruc.NewUsecase(realRepo),
		healthcheckUC: hcuc.NewUsecase(healthcheckRepo),
		auditUC:       auc.NewUsecase(auditRepo),
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
	// TODO: clean it so no API_KEY used
	// config.AllowHeaders = []string{"LBOS_API_KEY", "Access-Control-Allow-Headers", "Origin", "Accept", "X-Requested-With", "Content-Type", "Access-Control-Request-Method", "Access-Control-Request-Headers"} //nolint:lll
	config.AllowHeaders = []string{"Access-Control-Allow-Headers", "Origin", "Accept", "X-Requested-With", "Content-Type", "Access-Control-Request-Method", "Access-Control-Request-Headers"} //nolint:lll
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
	// rapi.Use(zrest.AuthAPIKey("CRISPY_API_KEY")) // FIXME: current CRISPY_API_KEY flow needs refactor

	zrest.RegisterHTTPEndpoint(rapi, a.zoneUC, a.auditUC)
	crest.RegisterHTTPEndpoint(rapi, a.clusterUC, a.auditUC)
	srest.RegisterHTTPEndpoint(rapi, a.serviceUC, a.auditUC)
	rrest.RegisterHTTPEndpoint(rapi, a.realUC, a.auditUC)
	nrest.RegisterHTTPEndpoint(rapi, a.nodeUC, a.auditUC)

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
