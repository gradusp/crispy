package rest

import (
	"github.com/gin-gonic/gin"

	"github.com/gradusp/crispy/internal/audit"
	"github.com/gradusp/crispy/internal/service"
)

func RegisterHTTPEndpoint(router *gin.RouterGroup, suc service.Usecase, auc audit.Usecase) {
	h := NewHandler(suc, auc)

	services := router.Group("/services")
	{
		services.POST("", h.Create)
		services.GET("", h.Get)
		services.GET("/:id", h.GetByID)
		services.DELETE("/:id", h.Delete)
	}
}
