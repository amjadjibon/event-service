package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type makeReservationInput struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
}

// MakeReservation godoc
// @Summary Make reservation for a workshop
// @Description Make reservation for a workshop by providing its ID, name and email
// @Tags events
// @Accept json
// @Produce json
// @Param id path int true "Workshop ID"
// @Param input body makeReservationInput true "Reservation details"
// @Success 200 {object} model.MakeReservationResponse
// @Router /events/{id}/reservation [post]
func (e EventHandler) MakeReservation(c *gin.Context) {
	var input makeReservationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workshopID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reservation, err := e.repo.MakeReservation(workshopID, input.Name, input.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reservation)
	return
}
