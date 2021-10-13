package middleware

import (
	"fmt"
	"gin_demo/response"
	"github.com/gin-gonic/gin"
)

func RecoveryMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				response.Fail(ctx, fmt.Sprint(err), nil)
			}
		}()
		ctx.Next()
	}
}