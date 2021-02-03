package rest

import (
	"github.com/gin-gonic/gin"

	"github.com/gradusp/crispy/internal/zone"
)

func RegisterHTTPEndpoint(router *gin.RouterGroup, uc zone.Usecase) {
	h := NewHandler(uc)

	securityZones := router.Group("/zones")
	{
		securityZones.POST("", h.Create)
		securityZones.GET("", h.Get)
		securityZones.GET("/:id", h.GetByID)
		securityZones.PUT("/:id", h.Update)
		securityZones.DELETE("/:id", h.Delete)
	}
}
