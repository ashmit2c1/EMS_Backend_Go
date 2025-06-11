package routes

import (
	"ems_backend_go/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func signUp(cntxt *gin.Context) {
	var user models.User
	err := cntxt.ShouldBindJSON(&user)

	if err != nil {
		cntxt.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse data", "error": err})
	}
	err = user.Save()
	if err != nil {
		cntxt.JSON(http.StatusInternalServerError, gin.H{"message": "There was some error", "error": err.Error()})
		return
	}
	cntxt.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user": user})

}
