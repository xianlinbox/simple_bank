package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xianlinbox/simple_bank/api/security"
)

const (
	authorizationHeader = "Authorization"
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
		authType := authFields[0]
		if authType != authorizationType {
			err := errors.New("authorization type is not supported")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		access_token := authFields[1]
		payload,err :=tokenMaker.VerifyToken(access_token)
		if err != nil {
			err := errors.New("access_token is not valid")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.Set(AUTH_KEY, payload)
		c.Next()
	}

}