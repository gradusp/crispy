package rest

import (
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/gradusp/crispy/internal/audit"
	"github.com/gradusp/crispy/internal/service"
)

type Handler struct {
	suc service.Usecase
	auc audit.Usecase
}

func NewHandler(suc service.Usecase, auc audit.Usecase) *Handler {
	return &Handler{
		suc: suc,
		auc: auc,
	}
}

type request struct {
	ClusterID     string `json:"clusterId" binding:"required"`
	RoutingType   string `json:"routingType" binding:"required"`
	BalancingType string `json:"balancingType" binding:"required"`
	Bandwidth     int    `json:"bandwidth" binding:"required"`
	Proto         string `json:"proto" binding:"required"`
	Addr          string `json:"addr" binding:"required"`
	Port          int    `json:"port" binding:"required"`
}

func (h *Handler) Create(c *gin.Context) {
	var req request
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": http.StatusText(http.StatusBadRequest),
			"note":    "input malformed",
		})
		return
	}

	res, err := h.suc.Create(c.Request.Context(),
		req.ClusterID,
		req.RoutingType, req.BalancingType, req.Proto, net.ParseIP(req.Addr), req.Bandwidth, req.Port)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrAlreadyExist):
			loc := fmt.Sprintf("%s/%s", c.FullPath(), res.ID)
			c.Header("Location", loc)
			c.AbortWithStatus(http.StatusSeeOther) // FIXME: rework to c.Status()
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": http.StatusText(http.StatusInternalServerError),
				"err":     err,
			})
			return
		}
	}

	// TODO: audit

	c.JSON(http.StatusCreated, res)
}

func (h *Handler) Get(c *gin.Context) {
	res, err := h.suc.Get(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": http.StatusText(http.StatusInternalServerError),
			"err":     err,
		})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetByID(c *gin.Context) {
	res, err := h.suc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": http.StatusText(http.StatusInternalServerError),
			"err":     err,
		})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) Delete(c *gin.Context) {
	err := h.suc.Delete(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": http.StatusText(http.StatusInternalServerError),
			"err":     err,
		})
		return
	}

	// TODO: audit

	c.Status(http.StatusNoContent)
}
