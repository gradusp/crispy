// TODO: https://github.com/caarlos0/env
// TODO: https://pkg.go.dev/go.uber.org/zap

package server

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gradusp/crispy/swagger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/hashicorp/consul/api"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/gradusp/crispy/balancingservice"
	"github.com/gradusp/crispy/cluster"
	"github.com/gradusp/crispy/securityzone"

	szhttp "github.com/gradusp/crispy/securityzone/delivery/http"
	szpg "github.com/gradusp/crispy/securityzone/repository/pgsql"
	szusecase "github.com/gradusp/crispy/securityzone/usecase"

	chttp "github.com/gradusp/crispy/cluster/delivery/http"
	cpg "github.com/gradusp/crispy/cluster/repository/pgsql"
	cuc "github.com/gradusp/crispy/cluster/usecase"

	//ohttp "github.com/gradusp/crispy/order/delivery/http"
	//opg "github.com/gradusp/crispy/order/repository/pgsql"
	//ouc "github.com/gradusp/crispy/order/usecase"

	bspg "github.com/gradusp/crispy/balancingservice/repository/pgsql"
	bsuc "github.com/gradusp/crispy/balancingservice/usecase"
)

type App struct {
	httpServer *http.Server

	securityZoneUC securityzone.Usecase
	clusterUC      cluster.Usecase
	//orderUC            order.Usecase
	balancingserviceUC balancingservice.Usecase
}

func NewApp() *App {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found", err)
	}

	db := initDB()
	kv := initConsul()
	l, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	securityZoneRepo := szpg.NewSecurityzoneRepo(db, kv, l.Sugar())
	clusterRepo := cpg.NewClusterRepo(db, kv, l.Sugar())
	//orderRepo := opg.NewOrderRepo(db, kv, l.Sugar())
	balancingserviceRepo := bspg.NewBalancingserviceRepo(db, kv, l.Sugar())

	return &App{
		clusterUC:      cuc.NewClusterUsecase(clusterRepo),
		securityZoneUC: szusecase.NewSecurityZoneUseCase(securityZoneRepo),
		//orderUC:            ouc.NewOrderUsecase(orderRepo),
		balancingserviceUC: bsuc.NewBalancingserviceUsecase(balancingserviceRepo),
	}
}

func (a *App) Run(port string) error {
	// logger, _ := zap.NewProduction()

	// Init gin handler
	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())
	// router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	// router.Use(ginzap.RecoveryWithZap(logger, true))

	// Set up gin CORS official middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost"} // FIXME: eliminate hardcode
	config.AllowMethods = []string{"GET", "HEAD", "OPTIONS", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"LBOS_API_KEY", "Access-Control-Allow-Headers", "Origin", "Accept", "X-Requested-With", "Content-Type", "Access-Control-Request-Method", "Access-Control-Request-Headers"} //nolint:lll
	config.AllowCredentials = true
	router.Use(cors.New(config))

	// Set up http handlers
	// SignIn endpoints
	// ...

	// embedding of static assets of SwaggerUI
	f, _ := fs.Sub(swagger.UI, "ui")
	router.StaticFS("/swagger", http.FS(f))
	router.GET("openapi.yml", func(c *gin.Context) {
		file, _ := swagger.OpenAPI.ReadFile("openapi.yml")
		c.Data(http.StatusOK, "application/yaml", file)
	})
	// TODO: implement favicon path
	// https://github.com/gin-gonic/examples/blob/master/assets-in-binary/example02/main.go#L34
	// TODO: implement gzip
	// https://github.com/gin-contrib/gzip

	// API Endpoints
	rapi := router.Group("/api/v1")
	rapi.Use(szhttp.AuthAPIKey("LBOS_API_KEY")) // FIXME: current LBOS_API_KEY flow is wrong

	szhttp.RegisterHTTPEndpoint(rapi, a.securityZoneUC)
	chttp.RegisterHTTPEndpoint(rapi, a.clusterUC)
	//ohttp.RegisterHTTPEndpoint(rapi, a.orderUC, a.balancingserviceUC) // FIXME: figure out why two UC here?

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
