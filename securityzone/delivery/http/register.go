package http

import (
	"github.com/gin-gonic/gin"
	"github.com/gradusp/crispy/ctrl/securityzone"
)

func RegisterHTTPEndpoint(router *gin.RouterGroup, uc securityzone.Usecase) {
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
