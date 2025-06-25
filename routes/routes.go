package routes

import (
	"ems_backend_go/middleware"

	"github.com/gin-gonic/gin"
)

func GetRoutes(server *gin.Engine) {
	// EVENT METHODS
	server.GET("/events", getEvents)
	server.GET("/events/:event_id", getEventByID)
	// PROTECTED METHODS
	server.POST("/events", middleware.Authenticate, createEvent)
	server.DELETE("/events", middleware.Authenticate, deleteAllEvents)
	server.DELETE("/events/:event_id", middleware.Authenticate, deleteEventByID)
	server.PUT("/events/:event_id", middleware.Authenticate, updateEvent)
	// USER METHODS
	server.POST("/signup", signUp)
	server.POST("/login", loginUser)
	server.GET("/users", getAllUsers)
	server.Run(":8080")
	// REGISTRATIONS
	server.POST("/events/:id/register", middleware.Authenticate, registerForEvent)
	server.DELETE("/events/:id/register", middleware.Authenticate, cancelRegistrationForEvent)

}
