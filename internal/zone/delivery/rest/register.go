package rest

import (
	"github.com/gin-gonic/gin"

	"github.com/gradusp/crispy/internal/audit"
	"github.com/gradusp/crispy/internal/zone"
)

func RegisterHTTPEndpoint(router *gin.RouterGroup, zuc zone.Usecase, auc audit.Usecase) {
	h := NewHandler(zuc, auc)

	securityZones := router.Group("/zones")
	{
		securityZones.POST("", h.Create)
		securityZones.GET("", h.Get)
		securityZones.GET("/:id", h.GetByID)
		securityZones.PUT("/:id", h.Update)
		securityZones.DELETE("/:id", h.Delete)
	}
}
