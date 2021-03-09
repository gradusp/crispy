package rest

import (
	"github.com/gradusp/crispy/internal/audit"
	"github.com/gradusp/crispy/internal/real"

	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndpoint(router *gin.RouterGroup, ruc real.Usecase, auc audit.Usecase) {
	h := NewHandler(ruc, auc)

	reals := router.Group("/reals")
	{
		reals.POST("", h.Create)
		reals.GET("", h.Get)
		reals.GET("/:id", h.GetByID)
		reals.DELETE("/:id", h.Delete)
	}
}
