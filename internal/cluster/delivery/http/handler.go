package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/gradusp/crispy/internal/cluster"
)

type Handler struct {
	usecase cluster.Usecase
}

func NewHandler(uc cluster.Usecase) *Handler {
	return &Handler{
		usecase: uc,
	}
}

type request struct {
	Name     string `json:"name" binding:"required"`
	ZoneID   string `json:"zoneId" binding:"required"`
	Capacity int64  `json:"capacity" binding:"required"`
}

func (h *Handler) Create(c *gin.Context) {
	var req request
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": http.StatusText(http.StatusBadRequest),
		})
		return
	}

	res, err := h.usecase.Create(c.Request.Context(), req.ZoneID, req.Name, req.Capacity)
	if err != nil {
		if errors.Is(err, cluster.ErrAlreadyExist) {
			// TODO: 303 status is not good here since there is 3 params
			//loc := fmt.Sprintf("%s/%s", c.FullPath(), res.ID)
			//c.Header("Location", loc)
			//c.AbortWithStatus(http.StatusSeeOther)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": http.StatusText(http.StatusBadRequest),
				"note":    cluster.ErrAlreadyExist.Error(),
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	c.JSON(http.StatusCreated, res)
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

type updateRequest struct {
	Name     string `json:"name"`
	ZoneID   string `json:"zoneId"`
	Capacity int64  `json:"capacity"`
}

func (h *Handler) Update(c *gin.Context) {
	var req updateRequest
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": http.StatusText(http.StatusBadRequest),
		})
		return
	}

	err := h.usecase.Update(c.Request.Context(), c.Param("id"), req.Name, req.Capacity)
	if err != nil {
		switch {
		case errors.Is(err, cluster.ErrNotFound):
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": http.StatusText(http.StatusBadRequest),
				"note":    cluster.ErrNotFound.Error(),
			})
			return
		case errors.Is(err, cluster.ErrAlreadyExist):
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": http.StatusText(http.StatusBadRequest),
				"note":    cluster.ErrAlreadyExist.Error(),
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
		case errors.Is(err, cluster.ErrHaveServices):
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": http.StatusText(http.StatusBadRequest),
				"note":    cluster.ErrHaveServices.Error(),
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
