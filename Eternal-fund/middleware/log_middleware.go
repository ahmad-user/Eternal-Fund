package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LogMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		t := time.Now()

		latency := time.Since(t)
		clientIP := ctx.ClientIP()
		method := ctx.Request.Method
		statusCode := ctx.Writer.Status()
		userAgent := ctx.Request.UserAgent()
		path := ctx.Request.URL.Path

		log.Printf("[LOG] %s - [%v] \"%s %s %d %v \"%s\"\n", clientIP, t, method, path, statusCode, latency, userAgent)
	}
}

// func LogMiddleware2(ctx *gin.Context) {

// 	t := time.Now()

// 	latency := time.Since(t)
// 	clientIP := ctx.ClientIP()
// 	method := ctx.Request.Method
// 	statusCode := ctx.Writer.Status()
// 	userAgent := ctx.Request.UserAgent()
// 	path := ctx.Request.URL.Path

// 	log.Printf("[LOG] %s - [%v] \"%s %s %d %v \"%s\"\n", clientIP, t, method, path, statusCode, latency, userAgent)

// }
