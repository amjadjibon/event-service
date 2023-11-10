package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// EventWorkshopList godoc
// @Summary Get workshop list by event ID
// @Description Get workshop list by providing event ID
// @Tags events
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {object} model.WorkshopList
// @Router /events/{id}/workshops [get]
func (e EventHandler) EventWorkshopList(c *gin.Context) {
	eventID := c.Param("id")
	eventIDX, err := strconv.Atoi(eventID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workshops, err := e.repo.GetWorkshopList(eventIDX)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, workshops)
	return
}
