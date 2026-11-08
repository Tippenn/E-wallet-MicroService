package middleware

import (
	"net/http"
	"strings"

	"e-wallet-system-api/shared/utils"

	"e-wallet-system-api/domains/JWTVerification"
	"e-wallet-system-api/domains/JWTVerification/models/request"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(uc JWTVerification.ServerJWTUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid token"})
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		//signature + format check
		userId, err := utils.ParseToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"auth error1": err.Error()})
			return
		}

		//server side check
		_, err = uc.CheckJWTServer(c.Request.Context(), &request.CheckJWTRequest{UserId: userId})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"auth error2": err.Error()})
			return
		}

		c.Set("user_Id", userId)
		c.Next()
	}

}
