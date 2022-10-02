package middlewares

import (
	"net/http"
	"strings"
	"task-5-vix-btpns-HaiqalRamanizarAlFajri/app"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.JSON(401, gin.H{"error": "request does not contain an access token"})
			ctx.Abort()
			return
		}

		err := app.ValidateToken(strings.Split(tokenString, "Bearer ")[1])
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"status": "F", "message": err.Error(), "data": nil})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
