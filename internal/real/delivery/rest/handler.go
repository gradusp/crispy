package rest

import (
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/gradusp/crispy/internal/audit"
	"github.com/gradusp/crispy/internal/real"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	ruc real.Usecase
	auc audit.Usecase
}

func NewHandler(ruc real.Usecase, auc audit.Usecase) *Handler {
	return &Handler{
		ruc: ruc,
		auc: auc,
	}
}

type realRequest struct {
	ServiceID string `form:"serviceId" json:"serviceId" binding:"required"`
	Addr      net.IP `form:"addr" json:"addr" binding:"required"`
	Port      int    `json:"port" binding:"required"`
}

func (h *Handler) Create(c *gin.Context) {
	var req realRequest
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": http.StatusText(http.StatusBadRequest),
		})
		return
	}

	res, err := h.ruc.Create(c.Request.Context(), req.ServiceID, req.Addr, req.Port)
	if err != nil {
		switch {
		case errors.Is(err, real.ErrAlreadyExist):
			loc := fmt.Sprintf("%s/%s", c.FullPath(), res.ID)
			c.Header("Location", loc)
			c.Status(http.StatusSeeOther)
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": http.StatusText(http.StatusInternalServerError),
			})
			return
		}
	}

	// TODO: audit

	c.JSON(http.StatusCreated, res)
}

// getRealRequest presents struct for binding query parameters into strings
type getRealRequest struct {
	ServiceID string `form:"serviceId" json:"serviceId"`
	Addr      string `form:"addr" json:"addr"`
}

func (h *Handler) Get(c *gin.Context) {
	var req getRealRequest
	err := c.BindQuery(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	res, err := h.ruc.Get(c.Request.Context(), req.ServiceID, req.Addr)
	if err != nil {
		switch {
		case errors.Is(err, real.ErrWrongQuery):
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": http.StatusText(http.StatusBadRequest),
			})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": http.StatusText(http.StatusInternalServerError),
			})
			return
		}
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetByID(c *gin.Context) {
	res, err := h.ruc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": http.StatusText(http.StatusInternalServerError),
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) Delete(c *gin.Context) {
	err := h.ruc.Delete(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	// TODO audit
	//res := &model.Cluster{
	//	Observable: model.Observable{Subs: new(list.List)},
	//}
	//a := &model.Audit{
	//	Entity: "cluster",
	//	Action: "delete",
	//	Who:    c.Request.RemoteAddr + " -- " + c.Request.UserAgent(),
	//	What:   fmt.Sprintf(`{"id":"%s"}`, c.Param("id")),
	//}
	//res.Subscribe(h.auc)
	//res.Fire(c.Request.Context(), a)
	//res.Unsubscribe(h.auc)

	c.Status(http.StatusNoContent)
}
