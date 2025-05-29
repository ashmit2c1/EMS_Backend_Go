package routes

import (
	"ems_backend_go/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetRoutes(server *gin.Engine) {

	server.GET("/events", getEvents)
	server.GET("/events/:event_id", getEventByID)
	server.POST("/events", createEvent)
	server.DELETE("/events", deleteAllEvents)
	server.Run(":8080")
}

func getEvents(cntxt *gin.Context) {
	events, err := models.GetAllEvents()
	if err != nil {
		cntxt.JSON(http.StatusInternalServerError, gin.H{"message": "There was some error", "error": err.Error()})
	}
	cntxt.JSON(http.StatusOK, gin.H{"message": "GET Request", "events": events})
}

func getEventByID(cntxt *gin.Context) {
	id, err := strconv.ParseInt(cntxt.Param("event_id"), 10, 64)
	if err != nil {
		cntxt.JSON(http.StatusBadRequest, gin.H{"message": "There was some error in fetching the event ID", "error": err.Error()})
		return
	}
	event, err := models.GetEventByID(id)
	if err != nil {
		cntxt.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch the event", "error": err.Error()})
	}
	cntxt.JSON(http.StatusOK, gin.H{"message": "GET Request", "event": event})

}
func createEvent(cntxt *gin.Context) {
	var event models.Event
	err := cntxt.ShouldBindJSON(&event)

	if err != nil {
		cntxt.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse data", "error": err})
	}
	event.CreatedBy = 1
	err = event.Save()
	if err != nil {
		cntxt.JSON(http.StatusBadRequest, gin.H{"message": "Could not proceed with request", "error": err.Error()})
	}
	cntxt.JSON(http.StatusCreated, gin.H{"message": "POST Request", "event": event})
}

func deleteAllEvents(cntxt *gin.Context) {
	err := models.DeleteAllEvents()
	if err != nil {
		cntxt.JSON(http.StatusInternalServerError, gin.H{"message": "There was some error", "error": err.Error()})
	}
	cntxt.JSON(http.StatusOK, gin.H{"message": "DELETE Request"})
}
