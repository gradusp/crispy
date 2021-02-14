package rest

import (
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/gradusp/crispy/internal/audit"
	"github.com/gradusp/crispy/internal/model"
	"github.com/gradusp/crispy/internal/zone"
)

type Handler struct {
	zuc zone.Usecase
	auc audit.Usecase
}

func NewHandler(zuc zone.Usecase, auc audit.Usecase) *Handler {
	return &Handler{
		zuc: zuc,
		auc: auc,
	}
}

type zoneInput struct {
	Name string `json:"name" binding:"required"`
}

func (h *Handler) Create(c *gin.Context) {
	var req zoneInput                        // TODO: test case should test garbage json on input
	if err := c.BindJSON(&req); err != nil { // FIXME: produces error event @ gin.logger which is not obvious
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": http.StatusText(http.StatusBadRequest),
		})
		return
	}

	res, err := h.zuc.Create(c.Request.Context(), req.Name)
	if err != nil {
		switch {
		case errors.Is(err, zone.ErrZoneAlreadyExist):
			// The solution is related to 303 suggestion of RFC 7231
			// https://tools.ietf.org/html/rfc7231#section-4.3.3
			loc := fmt.Sprintf("%s/%s", c.FullPath(), res.ID)
			c.Header("Location", loc)
			c.AbortWithStatus(http.StatusSeeOther)
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": http.StatusText(http.StatusInternalServerError),
			})
			return
		}
	}

	j, err := json.Marshal(&res)
	if err != nil {
		panic(err)
	}
	a := &model.Audit{
		Entity: "zone",
		Action: "create",
		Who:    c.Request.RemoteAddr + " -- " + c.Request.UserAgent(),
		What:   string(j),
	}
	res.Observable = model.Observable{Subs: new(list.List)}
	res.Subscribe(h.auc)
	res.Fire(c.Request.Context(), a)
	res.Unsubscribe(h.auc)

	c.JSON(http.StatusCreated, res)
}

func (h *Handler) Get(c *gin.Context) {
	res, err := h.zuc.Get(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": http.StatusText(http.StatusInternalServerError),
		})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetByID(c *gin.Context) {
	res, err := h.zuc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		switch {
		case errors.Is(err, zone.ErrZoneNotFound):
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

func (h *Handler) Update(c *gin.Context) {
	var req zoneInput
	if err := c.BindJSON(&req); err != nil { // FIXME: produces error event @ gin.logger which is not obvious
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": http.StatusText(http.StatusBadRequest),
		})
		return
	}

	err := h.zuc.Update(c.Request.Context(), c.Param("id"), req.Name)
	if err != nil {
		switch {
		case errors.Is(err, zone.ErrZoneAlreadyExist) ||
			errors.Is(err, zone.ErrZoneNotFound):
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

	res := &model.Zone{
		Observable: model.Observable{Subs: new(list.List)},
	}
	a := &model.Audit{
		Entity: "zone",
		Action: "update",
		Who:    c.Request.RemoteAddr + " -- " + c.Request.UserAgent(),
		What:   fmt.Sprintf(`{"id":"%s","name":"%s"}`, c.Param("id"), req.Name),
	}
	res.Subscribe(h.auc)
	res.Fire(c.Request.Context(), a)
	res.Unsubscribe(h.auc)

	c.Status(http.StatusOK)
}

func (h *Handler) Delete(c *gin.Context) {
	err := h.zuc.Delete(c.Request.Context(), c.Param("id"))
	if err != nil {
		switch {
		case errors.Is(err, zone.ErrZoneHaveClusters):
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": http.StatusText(http.StatusBadRequest),
				"note":    zone.ErrZoneHaveClusters.Error(),
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

	res := &model.Zone{
		Observable: model.Observable{Subs: new(list.List)},
	}
	a := &model.Audit{
		Entity: "zone",
		Action: "delete",
		Who:    c.Request.RemoteAddr + " -- " + c.Request.UserAgent(),
		What:   fmt.Sprintf(`{"id":"%s"}`, c.Param("id")),
	}
	res.Subscribe(h.auc)
	res.Fire(c.Request.Context(), a)
	res.Unsubscribe(h.auc)

	c.Status(http.StatusNoContent)
}
