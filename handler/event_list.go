package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"event-service/model"
)

// EventList godoc
// @Summary Get list of events
// @Description Get list of events with pagination
// @Tags events
// @Accept json
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} map[string]interface{}
// @Router /events [get]
func (e EventHandler) EventList(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the total number of events
	totalEvents, err := e.repo.GetActiveEventsCount()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get pagination event list
	events, err := e.repo.GetActiveEventsWithPagination(limit, offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	totalPage := 0
	if totalEvents%limit == 0 {
		totalPage = totalEvents / limit
	} else {
		totalPage = (totalEvents / limit) + 1
	}

	pagination := model.Pagination{
		Total:       totalEvents,
		PerPage:     limit,
		TotalPages:  totalPage,
		CurrentPage: offset + 1,
	}

	// Construct the response JSON
	response := gin.H{
		"events":     events,
		"pagination": pagination,
	}

	c.JSON(http.StatusOK, response)
	return
}
