package middleware

import (
	"ems_backend_go/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authenticate(cntxt *gin.Context) {
	token := cntxt.Request.Header.Get("Authorisation")
	if token == "" {
		cntxt.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not Authorised"})
		return
	}
	err := utils.VerifyToken(token)
	if err != nil {
		cntxt.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "There was some error", "error": err.Error()})
		return
	}
	cntxt.Next()
}
