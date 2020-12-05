package http

import (
	"github.com/gin-gonic/gin"
	"github.com/gradusp/crispy/ctrl/cluster"
)

func RegisterHTTPEndpoint(router *gin.RouterGroup, uc cluster.Usecase) {
	h := NewHandler(uc)

	clusters := router.Group("/clusters")
	{
		clusters.POST("", h.Create)
		clusters.GET("", h.Get)
		clusters.GET("/:id", h.GetByID)
		clusters.PUT("/:id", h.Update)
		clusters.DELETE("/:id", h.Delete)
	}
}
