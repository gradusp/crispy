package http

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gradusp/crispy/ctrl/cluster"
	"net/http"
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
	Name           string `json:"name" binding:"required"`
	SecurityZoneID string `json:"securityZoneId" binding:"required"`
	Capacity       int64  `json:"capacity"`
	//SecurityZone *model.SecurityZone `json:"securityZone"`
}

type response struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	SecurityZoneID string `json:"securityZoneId"`
	Capacity       int64  `json:"capacity"`
	Usage          int64  `json:"usage"`
}

func (h *Handler) Create(c *gin.Context) {
	var req request
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": http.StatusText(http.StatusBadRequest),
		})
	}

	res, err := h.usecase.Create(c.Request.Context(), req.SecurityZoneID, req.Name, req.Capacity)
	if err != nil {
		if errors.Is(err, cluster.ErrClusterAlreadyExist) {
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

	c.JSON(http.StatusCreated, &response{
		ID:             res.ID,
		Name:           res.Name,
		SecurityZoneID: res.SecurityZoneID,
		Capacity:       res.Capacity,
		Usage:          res.Usage,
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
