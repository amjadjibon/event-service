package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// EventDetail godoc
// @Summary Get event details by ID
// @Description Get event details by providing its ID
// @Tags events
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {object} model.Event
// @Router /events/{id} [get]
func (e EventHandler) EventDetail(c *gin.Context) {
	id := c.Param("id")

	idx, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event, err := e.repo.GetEventByID(idx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, event)
}
