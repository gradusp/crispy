// TODO: https://github.com/caarlos0/env
// TODO: https://pkg.go.dev/go.uber.org/zap

package server

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/gradusp/crispy/ctrl/security_zone"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	szhttp "github.com/gradusp/crispy/ctrl/security_zone/delivery/http"
	szpg "github.com/gradusp/crispy/ctrl/security_zone/repository/pgsql"
	szusecase "github.com/gradusp/crispy/ctrl/security_zone/usecase"
)

type App struct {
	httpServer *http.Server

	securityZoneUC security_zone.Usecase
}

func NewApp() *App {
	db := initDB()

	securityZoneRepo := szpg.NewSecurityZoneRepo(db)

	return &App{
		securityZoneUC: szusecase.NewSecurityZoneUseCase(securityZoneRepo),
	}
}

func (a *App) Run(port string) error {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	// Init gin handler
	router := gin.Default()
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
	db := pg.Connect(&pg.Options{
		ApplicationName: "lbosCtrl",
		Database:        "postgres",
		User:            "postgres",
		Password:        "secret",
	})
	//defer db.Close()

	return db
}
