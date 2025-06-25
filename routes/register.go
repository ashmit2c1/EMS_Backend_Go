package routes

import (
	"ems_backend_go/models"
	"ems_backend_go/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func registerForEvent(cntxt *gin.Context) {
	token := cntxt.Request.Header.Get("Authorisation")
	usID, err := utils.GetUserIDFromToken(token)
	if err != nil {
		cntxt.JSON(http.StatusInternalServerError, gin.H{"message": "There was some error in parsing"})
		return
	}
	eventId, err := strconv.ParseInt(cntxt.Param("id"), 10, 64)
	if err != nil {
		cntxt.JSON(http.StatusBadRequest, gin.H{"message": "There was some error in fetching the event ID", "error": err.Error()})
		return
	}
	event, err := models.GetEventByID(eventId)
	if err != nil {
		cntxt.JSON(http.StatusInternalServerError, gin.H{"message": "There was some error in parsing"})
		return
	}
	err = event.Register(usID)
	if err != nil {
		cntxt.JSON(http.StatusInternalServerError, gin.H{"message": "There was some error in registering", "error": err.Error()})
		return
	}
	cntxt.JSON(http.StatusOK, gin.H{"message": "Congratulations you have registered for the event"})

}

func cancelRegistrationForEvent(cntxt *gin.Context) {

}
