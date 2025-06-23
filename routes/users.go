package routes

import (
	"ems_backend_go/models"
	"ems_backend_go/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func loginUser(cntxt *gin.Context) {
	var user models.User
	err := cntxt.ShouldBindJSON(&user)

	if err != nil {
		cntxt.JSON(http.StatusBadRequest, gin.H{"message": "There was some error", "error": err.Error()})
		return
	}
	err = user.ValidateCredentials()
	if err != nil {
		cntxt.JSON(http.StatusUnauthorized, gin.H{"message": "There was some error", "error": err.Error()})
		return
	}
	userID, err := models.FetchIDByEmail(user.Email)
	if err != nil {
		cntxt.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve user ID",
			"error":   err.Error(),
		})
		return
	}
	token, err := utils.GenerateToken(user.Email, userID)
	if err != nil {
		cntxt.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to generate token",
			"error":   err.Error(),
		})
		return
	}
	cntxt.JSON(http.StatusOK, gin.H{
		"message": "User logged in successfully",
		"token":   token,
	})
}
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

func getAllUsers(cntxt *gin.Context) {
	users, err := models.GetUsers()
	if err != nil {
		cntxt.JSON(http.StatusInternalServerError, gin.H{"message": "There was some error", "error": err.Error()})
		return
	}
	cntxt.JSON(http.StatusOK, gin.H{"message": "GET Request", "users": users})
}
