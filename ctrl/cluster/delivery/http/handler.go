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
	Capacity       int64  `json:"capacity" binding:"required"`
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

func (h *Handler) Update(c *gin.Context) {
	var req request
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": http.StatusText(http.StatusBadRequest),
		})
		return
	}

	err := h.usecase.Update(c.Request.Context(), req.SecurityZoneID, c.Param("id"), req.Name, req.Capacity)
	if err != nil {
		if errors.Is(err, cluster.ErrRequestedSecZoneNotFound) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": http.StatusText(http.StatusBadRequest),
				"note":    cluster.ErrRequestedSecZoneNotFound.Error(),
			})
			return
		}
		if errors.Is(err, cluster.ErrClusterAlreadyExist) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": http.StatusText(http.StatusBadRequest),
				"note":    cluster.ErrClusterAlreadyExist.Error(),
			})
			return
		}
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
		fmt.Println("CLUSTER_HANDLER:118:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": http.StatusText(http.StatusInternalServerError),
		})
		return
	}
	c.Status(http.StatusNoContent)
}
