package http

import (
	"github.com/gin-gonic/gin"
	"github.com/gradusp/crispy/ctrl/security_zone"
)

func RegisterHTTPEndpoint(router *gin.RouterGroup, uc security_zone.Usecase) {
	h := NewHandler(uc)

	securityZones := router.Group("/security-zones")
	{
		securityZones.POST("", h.Create)
		securityZones.GET("", h.Get)
		securityZones.GET("/:id", h.GetByID)
		securityZones.PUT("/:id", h.Update)
		securityZones.DELETE("/:id", h.Delete)
	}
}
