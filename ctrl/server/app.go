// TODO: https://github.com/caarlos0/env
// TODO: https://pkg.go.dev/go.uber.org/zap

package server

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/gradusp/crispy/ctrl/cluster"
	"github.com/gradusp/crispy/ctrl/security_zone"
	"github.com/hashicorp/consul/api"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	chttp "github.com/gradusp/crispy/ctrl/cluster/delivery/http"
	cpg "github.com/gradusp/crispy/ctrl/cluster/repository/pgsql"
	cuc "github.com/gradusp/crispy/ctrl/cluster/usecase"

	szhttp "github.com/gradusp/crispy/ctrl/security_zone/delivery/http"
	szpg "github.com/gradusp/crispy/ctrl/security_zone/repository/pgsql"
	szusecase "github.com/gradusp/crispy/ctrl/security_zone/usecase"
)

type App struct {
	httpServer *http.Server

	clusterUC      cluster.Usecase
	securityZoneUC security_zone.Usecase
}

func NewApp() *App {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found", err)
	}

	db := initDB()
	kv := initConsul()

	clusterRepo := cpg.NewClusterRepo(db, kv)
	securityZoneRepo := szpg.NewSecurityZoneRepo(db, kv)

	return &App{
		clusterUC:      cuc.NewClusterUsecase(clusterRepo),
		securityZoneUC: szusecase.NewSecurityZoneUseCase(securityZoneRepo),
	}
}

func (a *App) Run(port string) error {
	// Init gin handler
	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())

	// Set up gin CORS official middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost"} // FIXME: eliminate hardcode
	config.AllowMethods = []string{"GET", "HEAD", "OPTIONS", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"LBOS_API_KEY", "Access-Control-Allow-Headers", "Origin", "Accept", "X-Requested-With", "Content-Type", "Access-Control-Request-Method", "Access-Control-Request-Headers"}
	config.AllowCredentials = true
	router.Use(cors.New(config))

	// Set up http handlers
	// SignIn endpoints
	// ...

	// API Endpoints
	api := router.Group("/api/v1")
	api.Use(szhttp.AuthAPIKey("LBOS_API_KEY")) // FIXME: current LBOS_API_KEY flow is wrong

	chttp.RegisterHTTPEndpoint(api, a.clusterUC)
	szhttp.RegisterHTTPEndpoint(api, a.securityZoneUC)

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

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}

func initDB() *pg.DB {
	addr := fmt.Sprintf("%s:%s", os.Getenv("LBOS_DB_HOST"), os.Getenv("LBOS_DB_PORT"))
	db := pg.Connect(&pg.Options{
		ApplicationName: "lbosCtrl",
		Addr:            addr,
		Database:        os.Getenv("LBOS_DB_USER"),
		User:            os.Getenv("LBOS_DB_USER"),
		Password:        os.Getenv("LBOS_DB_PASS"),
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
