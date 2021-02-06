package rest

import (
	"github.com/gin-gonic/gin"

	"github.com/gradusp/crispy/internal/audit"
	"github.com/gradusp/crispy/internal/cluster"
)

func RegisterHTTPEndpoint(router *gin.RouterGroup, cuc cluster.Usecase, auc audit.Usecase) {
	h := NewHandler(cuc, auc)

	clusters := router.Group("/clusters")
	{
		clusters.POST("", h.Create)
		clusters.GET("", h.Get)
		clusters.GET("/:id", h.GetByID)
		clusters.PUT("/:id", h.Update)
		clusters.DELETE("/:id", h.Delete)
	}
}
