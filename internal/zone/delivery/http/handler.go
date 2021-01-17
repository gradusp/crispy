package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/gradusp/crispy/internal/model"
	"github.com/gradusp/crispy/internal/zone"
)

type Handler struct {
	usecase zone.Usecase
}

func NewHandler(useCase zone.Usecase) *Handler {
	return &Handler{
		usecase: useCase,
	}
}

type szInput struct {
	Name string `json:"name" binding:"required"`
}

func (h *Handler) Create(c *gin.Context) {
	var req szInput                          // TODO: test case should test garbage json on input
	if err := c.BindJSON(&req); err != nil { // FIXME: produces error event @ gin.logger which is not obvious
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": http.StatusText(http.StatusBadRequest),
		})
		return
	}

	res, err := h.usecase.Create(c.Request.Context(), req.Name)
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

	c.JSON(http.StatusCreated, &model.Zone{
		ID:   res.ID,
		Name: res.Name,
	})
}

func (h *Handler) Get(c *gin.Context) {
	res, err := h.usecase.Get(c.Request.Context())
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
	res, err := h.usecase.GetByID(c.Request.Context(), c.Param("id"))
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
	var req szInput
	if err := c.BindJSON(&req); err != nil { // FIXME: produces error event @ gin.logger which is not obvious
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": http.StatusText(http.StatusBadRequest),
		})
		return
	}

	err := h.usecase.Update(c.Request.Context(), c.Param("id"), req.Name)
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
	c.Status(http.StatusOK)
}

func (h *Handler) Delete(c *gin.Context) {
	err := h.usecase.Delete(c.Request.Context(), c.Param("id"))
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
	c.Status(http.StatusNoContent)
}
