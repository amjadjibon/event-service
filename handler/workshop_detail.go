package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// WorkshopDetail godoc
// @Summary Get workshop details by ID
// @Description Get workshop details by providing its ID
// @Tags events
// @Accept json
// @Produce json
// @Param id path int true "Workshop ID"
// @Success 200 {object} model.WorkshopDetail
// @Router /workshops/{id} [get]
func (e EventHandler) WorkshopDetail(c *gin.Context) {
	workshopID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workshop, err := e.repo.GetWorkshopDetail(workshopID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, workshop)
	return
}
