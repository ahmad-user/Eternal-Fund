package middleware

import (
	"eternal-fund/usecase/service"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	CheckToken(roles ...string) gin.HandlerFunc
}
type authMiddleware struct {
	jwtService service.JwtService
}
type AuthHeader struct {
	Autheader string `header:"Authorization" required:"true"`
}

func (a *authMiddleware) CheckToken(roles ...string) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		var header AuthHeader
		if err := ctx.ShouldBindHeader(&header); err != nil {
			log.Println("Authorization header missing or malformed")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		token := strings.Replace(header.Autheader, "Bearer ", "", -1)
		claims, err := a.jwtService.ValidateToken(token)
		if err != nil {
			log.Println("Token validation error:", err)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userIdStr, ok := claims["userId"].(string)
		if !ok {
			log.Println("User ID not found in token claims")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			log.Println("Error converting user ID to int:", err)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("userID", userId) // Convert float64 ke int
		log.Printf("User ID set in context: %d", userId)
		validRole := false
		for _, role := range roles {
			if role == claims["role"] {
				validRole = true
				break
			}
		}
		if !validRole {
			log.Println(" invalid role ")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Next()

	}

}

func NewAuthMiddleware(jwtService service.JwtService) AuthMiddleware {
	return &authMiddleware{jwtService: jwtService}
}
