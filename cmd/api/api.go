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
	"github.com/go-pg/pg/v10"
	"github.com/hashicorp/consul/api"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"

	swagger "github.com/gradusp/crispy/api"
	"github.com/gradusp/crispy/assets"
	"github.com/gradusp/crispy/internal/cluster"
	"github.com/gradusp/crispy/internal/service"
	"github.com/gradusp/crispy/internal/zone"

	zhttp "github.com/gradusp/crispy/internal/zone/delivery/http"
	zpg "github.com/gradusp/crispy/internal/zone/repository/pgsql"
	szusecase "github.com/gradusp/crispy/internal/zone/usecase"

	chttp "github.com/gradusp/crispy/internal/cluster/delivery/http"
	cpg "github.com/gradusp/crispy/internal/cluster/repository/pgsql"
	cuc "github.com/gradusp/crispy/internal/cluster/usecase"

	spg "github.com/gradusp/crispy/internal/service/repository/pgsql"
	bsuc "github.com/gradusp/crispy/internal/service/usecase"
	//
	//ohttp "github.com/gradusp/crispy/order/delivery/http"
	//opg "github.com/gradusp/crispy/order/repository/pgsql"
	//ouc "github.com/gradusp/crispy/order/usecase"
)

// TODO: https://github.com/caarlos0/env
// TODO: https://pkg.go.dev/go.uber.org/zap

type App struct {
	httpServer *http.Server
	logger     *zap.Logger

	zoneUC    zone.Usecase
	clusterUC cluster.Usecase
	serviceUC service.Usecase
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
	db := initDB()
	kv := initConsul()

	zoneRepo := zpg.NewZonePostgresRepo(pool, kv, logger.Sugar())
	clusterRepo := cpg.NewClusterRepo(pool, kv, logger.Sugar())
	serviceRepo := spg.NewServiceRepo(db, kv, logger.Sugar())
	//orderRepo := opg.NewOrderRepo(db, kv, logger.Sugar())

	return &App{
		logger:    logger,
		clusterUC: cuc.NewClusterUsecase(clusterRepo),
		serviceUC: bsuc.NewServiceUsecase(serviceRepo),
		zoneUC:    szusecase.NewZoneUseCase(zoneRepo),
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
	// TODO: implement favicon path
	// https://github.com/gin-gonic/examples/blob/master/assets-in-binary/example02/main.go#L34
	// TODO: implement gzip
	// https://github.com/gin-contrib/gzip

	// API Endpoints
	rapi := router.Group("/api/v1")
	rapi.Use(zhttp.AuthAPIKey("CRISPY_API_KEY")) // FIXME: current CRISPY_API_KEY flow needs refactor

	zhttp.RegisterHTTPEndpoint(rapi, a.zoneUC)
	chttp.RegisterHTTPEndpoint(rapi, a.clusterUC)

	// HTTP Server
	a.httpServer = &http.Server{
		Addr:           ":8080",
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

func initDB() *pg.DB {
	addr := fmt.Sprintf("%s:%s", os.Getenv("CRISPY_DB_HOST"), os.Getenv("CRISPY_DB_PORT"))
	db := pg.Connect(&pg.Options{
		ApplicationName: "crispy",
		Addr:            addr,
		Database:        os.Getenv("CRISPY_DB_NAME"),
		User:            os.Getenv("CRISPY_DB_USER"),
		Password:        os.Getenv("CRISPY_DB_PASS"),
	})

	return db
}

func initConsul() *api.KV {
	cfg := &api.Config{
		Address: "10.56.204.194:18500",
	}
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	return client.KV()
}
