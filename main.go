package main

import (
	"ems_backend_go/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	server.GET("/events", getEvents)
	server.POST("/events", createEvent)
	server.Run(":8080")

}

func getEvents(cntxt *gin.Context) {
	events := models.GetAllEvents()
	cntxt.JSON(http.StatusOK, gin.H{"message": "GET Request", "events": events})
}

func createEvent(cntxt *gin.Context) {
	var event models.Event
	err := cntxt.ShouldBindJSON(&event)

	if err != nil {
		cntxt.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse data", "error": err})
	}
	event.ID = 1
	event.CreatedBy = 1
	event.Save()
	cntxt.JSON(http.StatusCreated, gin.H{"message": "POST Request", "event": event})
}
