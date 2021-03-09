package rest

import (
	"github.com/gin-gonic/gin"

	"github.com/gradusp/crispy/internal/audit"
	"github.com/gradusp/crispy/internal/zone"
)

func RegisterHTTPEndpoint(router *gin.RouterGroup, zuc zone.Usecase, auc audit.Usecase) {
	h := NewHandler(zuc, auc)

	zones := router.Group("/zones")
	{
		zones.POST("", h.Create)
		zones.GET("", h.Get)
		zones.GET("/:id", h.GetByID)
		zones.PUT("/:id", h.Update)
		zones.DELETE("/:id", h.Delete)
	}
}
