package rest

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gradusp/crispy/internal/audit"
	"github.com/gradusp/crispy/internal/node"
)

type Handler struct {
	nuc node.Usecase
	auc audit.Usecase
}

func NewHandler(nuc node.Usecase, auc audit.Usecase) *Handler {
	return &Handler{
		nuc: nuc,
		auc: auc,
	}
}

type nodeRequest struct {
	ClusterID string `json:"clusterId" binding:"required"`
	Addr      net.IP `json:"addr" binding:"required"`
	Hostname  string `json:"hostname"`
}

func (h *Handler) Create(c *gin.Context) {
	var req nodeRequest
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": http.StatusText(http.StatusBadRequest),
		})
		return
	}

	res, err := h.nuc.Create(c.Request.Context(), req.ClusterID, req.Hostname, req.Addr)
	if err != nil {
		switch {
		case errors.Is(err, node.ErrAlreadyExist):
			loc := fmt.Sprintf("%s/%d", c.FullPath(), res.ID)
			c.Header("Location", loc)
			c.Status(http.StatusSeeOther) // FIXME: rework to c.Status()
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
	//j, err := json.Marshal(&res)
	//if err != nil {
	//	panic(err)
	//}
	//a := &model.Audit{
	//	Entity: "node",
	//	Action: "create",
	//	Who:    c.Request.RemoteAddr + " -- " + c.Request.UserAgent(),
	//	What:   string(j),
	//}
	//res.Observable = model.Observable{Subs: new(list.List)}
	//res.Subscribe(h.auc)
	//res.Fire(c.Request.Context(), a)
	//res.Unsubscribe(h.auc)

	c.JSON(http.StatusCreated, res)
}

// getNodeRequest presents struct for binding query parameters into strings
type getNodeRequest struct {
	ClusterID string `form:"clusterId" json:"clusterId"`
	Addr      string `form:"addr" json:"addr"`
}

func (h *Handler) Get(c *gin.Context) {
	var req getNodeRequest
	err := c.BindQuery(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": http.StatusText(http.StatusBadRequest),
		})
		return
	}

	res, err := h.nuc.Get(c.Request.Context(), req.ClusterID, req.Addr)
	if err != nil {
		switch {
		case errors.Is(err, node.ErrWrongQuery):
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
	if id, err := strconv.Atoi(c.Param("id")); err == nil {
		res, err := h.nuc.GetByID(c.Request.Context(), id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": http.StatusText(http.StatusInternalServerError),
				"error":   err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, res)
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": http.StatusText(http.StatusBadRequest),
		})
		return
	}
}

func (h *Handler) Delete(c *gin.Context) {
	if id, err := strconv.Atoi(c.Param("id")); err == nil {
		err := h.nuc.Delete(c.Request.Context(), id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": http.StatusText(http.StatusInternalServerError),
			})
			return
		}
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": http.StatusText(http.StatusBadRequest),
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
