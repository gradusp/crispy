package rest

import (
	"github.com/gin-gonic/gin"

	"github.com/gradusp/crispy/internal/healthcheck"
	"github.com/gradusp/crispy/internal/real"
	"github.com/gradusp/crispy/internal/service"
)

func RegisterHTTPEndpoint(router *gin.RouterGroup, huc healthcheck.Usecase, ruc real.Usecase, suc service.Usecase) {
	h := NewHandler(huc, ruc, suc)

	services := router.Group("/services")
	{
		services.POST("", h.Create)
		services.GET("", h.Get)
		services.GET("/:id", h.GetByID)
		services.DELETE("/:id", h.Delete)
	}
}
