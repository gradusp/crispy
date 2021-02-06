package rest

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/gradusp/crispy/internal/audit"
	"github.com/gradusp/crispy/internal/healthcheck"
	"github.com/gradusp/crispy/internal/model"
	"github.com/gradusp/crispy/internal/real"
	"github.com/gradusp/crispy/internal/service"
)

type Handler struct {
	huc healthcheck.Usecase
	ruc real.Usecase
	suc service.Usecase
	auc audit.Usecase
}

func NewHandler(huc healthcheck.Usecase, ruc real.Usecase, suc service.Usecase, auc audit.Usecase) *Handler {
	return &Handler{
		huc: huc,
		ruc: ruc,
		suc: suc,
		auc: auc,
	}
}

type request struct {
	ClusterID     string               `json:"clusterId" binding:"required"`
	RoutingType   string               `json:"routingType" binding:"required"`
	BalancingType string               `json:"balancingType" binding:"required"`
	Proto         string               `json:"proto" binding:"required"`
	Addr          string               `json:"addr" binding:"required"`
	Port          int                  `json:"port" binding:"required"`
	Bandwidth     int                  `json:"bandwidth" binding:"required"`
	Healthchecks  []*model.Healthcheck `json:"healthchecks" binding:"required"`
	Reals         []*model.Real        `json:"reals" binding:"required"`
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

	// creating service
	svc, err := h.suc.Create(c.Request.Context(),
		req.ClusterID,
		req.RoutingType, req.BalancingType, req.Proto, net.ParseIP(req.Addr), req.Bandwidth, req.Port)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": http.StatusText(http.StatusInternalServerError),
			"err":     err,
		})
		return
	}

	// creating reals
	for _, r := range req.Reals {
		if _, err = h.ruc.Create(c.Request.Context(), svc.ID, r.Addr.To4(), r.HealthcheckAddr.To4(), r.Port, r.HealthcheckPort); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": http.StatusText(http.StatusInternalServerError),
				"err":     err,
			})
			return
		}
	}

	// creating healthchecks
	for _, hc := range req.Healthchecks {
		if _, err := h.huc.Create(c.Request.Context(),
			svc.ID,
			hc.HelloTimer, hc.ResponseTimer, hc.AliveThreshold, hc.DeadThreshold, hc.Quorum, hc.Hysteresis); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": http.StatusText(http.StatusInternalServerError),
				"err":     err,
			})
			return
		}
	}

	// FIXME: refactor -- handle error or review interface
	res, _ := h.suc.GetByID(c.Request.Context(), svc.ID)

	// TODO: should be refactored for DRY reason
	who := c.Request.RemoteAddr + " -- " + c.Request.UserAgent()
	j, err := json.Marshal(&res)
	if err != nil {
		panic(err)
	}
	what := `{"op":"create","obj":"service","dsc":` + string(j) + `}`
	h.auc.Create(c.Request.Context(), who, what)

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

	// TODO: should be refactored for DRY reason
	who := c.Request.RemoteAddr + " -- " + c.Request.UserAgent()
	what := fmt.Sprintf(`{"op":"delete","obj":"service","dsc":{"id":"%s"}}`, c.Param("id"))
	h.auc.Create(c.Request.Context(), who, what)

	c.Status(http.StatusNoContent)
}
