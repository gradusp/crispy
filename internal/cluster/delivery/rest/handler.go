package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/gradusp/crispy/internal/audit"
	"github.com/gradusp/crispy/internal/cluster"
)

type Handler struct {
	cuc cluster.Usecase
	auc audit.Usecase
}

func NewHandler(uc cluster.Usecase, auc audit.Usecase) *Handler {
	return &Handler{
		cuc: uc,
		auc: auc,
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

	res, err := h.cuc.Create(c.Request.Context(), req.ZoneID, req.Name, req.Capacity)
	if err != nil {
		if errors.Is(err, cluster.ErrAlreadyExist) {
			// TODO: 303 status is not good here since there is 3 params
			// loc := fmt.Sprintf("%s/%s", c.FullPath(), res.ID)
			// c.Header("Location", loc)
			// c.AbortWithStatus(http.StatusSeeOther)
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

	// TODO: should be refactored for DRY reason
	who := c.Request.RemoteAddr + " -- " + c.Request.UserAgent()
	j, err := json.Marshal(&res)
	if err != nil {
		panic(err)
	}
	what := `{"op":"create","obj":"cluster","dsc":` + string(j) + `}`
	h.auc.Create(c.Request.Context(), who, what)

	c.JSON(http.StatusCreated, res)
}

func (h *Handler) Get(c *gin.Context) {
	res, err := h.cuc.Get(c.Request.Context())
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
	res, err := h.cuc.GetByID(c.Request.Context(), c.Param("id"))
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

	err := h.cuc.Update(c.Request.Context(), c.Param("id"), req.Name, req.Capacity)
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

	// TODO: should be refactored for DRY reason
	who := c.Request.RemoteAddr + " -- " + c.Request.UserAgent()
	what := fmt.Sprintf(`{"op":"update","obj":"cluster","dsc":{"id":"%s","name":"%s","capacity":%d}}`, c.Param("id"), req.Name, req.Capacity)
	h.auc.Create(c.Request.Context(), who, what)

	c.Status(http.StatusOK)
}

func (h *Handler) Delete(c *gin.Context) {
	err := h.cuc.Delete(c.Request.Context(), c.Param("id"))
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

	// TODO: should be refactored for DRY reason
	who := c.Request.RemoteAddr + " -- " + c.Request.UserAgent()
	what := fmt.Sprintf(`{"op":"delete","obj":"cluster","dsc":{"id":"%s"}}`, c.Param("id"))
	h.auc.Create(c.Request.Context(), who, what)

	c.Status(http.StatusNoContent)
}
