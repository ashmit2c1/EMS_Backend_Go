package main

import (
	"ems_backend_go/db"
	"ems_backend_go/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	server := gin.Default()
	routes.GetRoutes(server)

}
