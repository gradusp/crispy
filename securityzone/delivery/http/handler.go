package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/gradusp/crispy/model"
	"github.com/gradusp/crispy/securityzone"
)

type Handler struct {
	usecase securityzone.Usecase
}

func NewHandler(useCase securityzone.Usecase) *Handler {
	return &Handler{
		usecase: useCase,
	}
}

type szInput struct {
	Name string `json:"name" binding:"required"`
}

func (h *Handler) Create(c *gin.Context) {
	var req szInput // TODO: test case should test garbage json on input
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": http.StatusText(http.StatusBadRequest),
		})
		return
	}

	res, err := h.usecase.Create(c.Request.Context(), req.Name)
	if err != nil {
		// This is related to 303 suggestion of RFC 7231
		// https://tools.ietf.org/html/rfc7231#section-4.3.3
		if errors.Is(err, securityzone.ErrSecurityzoneAlreadyExist) {
			loc := fmt.Sprintf("%s/%s", c.FullPath(), res.ID)
			c.Header("Location", loc)
			c.AbortWithStatus(http.StatusSeeOther)
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	c.JSON(http.StatusCreated, &model.SecurityZone{
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
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": http.StatusText(http.StatusInternalServerError),
		})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) Update(c *gin.Context) {
	var req szInput
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": http.StatusText(http.StatusBadRequest),
		})
		return
	}

	err := h.usecase.Update(c.Request.Context(), c.Param("id"), req.Name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) Delete(c *gin.Context) {
	err := h.usecase.Delete(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
