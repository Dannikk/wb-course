package api

import (
	"net/http"
	entity "ordermngmt/pkg/entity"

	"github.com/gin-gonic/gin"
)

func (h *Handler) AddRoutes(eng *gin.Engine) {
	eng.GET("/ping", h.Ping)
	eng.POST("/orders", h.GetOrder)
}

func (h Handler) Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func (h Handler) GetOrder(c *gin.Context) {
	var ctx = c.Request.Context()
	orderID := entity.OrderID{}

	if err := c.BindJSON(&orderID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.uc.GetOrder(ctx, orderID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"id": res})
	}
}
