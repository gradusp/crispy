package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/gradusp/crispy/internal/audit"
	"github.com/gradusp/crispy/internal/node"
)

func RegisterHTTPEndpoint(router *gin.RouterGroup, nuc node.Usecase, auc audit.Usecase) {
	h := NewHandler(nuc, auc)

	reals := router.Group("/nodes")
	{
		reals.POST("", h.Create)
		reals.GET("", h.Get)
		reals.GET("/:id", h.GetByID)
		reals.DELETE("/:id", h.Delete)
	}
}
