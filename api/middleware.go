package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xianlinbox/simple_bank/api/security"
)

const (
	authorizationHeader = "authorization"
	authorizationType = "Bearer"
	AUTH_KEY = "auth"
)
func authMiddleware(tokenMaker security.Maker) gin.HandlerFunc{
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader(authorizationHeader)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		authFields := strings.Fields(authorizationHeader)
		if len(authFields) != 2 {
			err := errors.New("authorization header is not valid")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		// authType := authFields[1].
		c.Next()
	}

}