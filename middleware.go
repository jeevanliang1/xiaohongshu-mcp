package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// corsMiddleware CORS 中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// errorHandlingMiddleware 错误处理中间件
func errorHandlingMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		logrus.Errorf("服务器内部错误: %v, path: %s", recovered, c.Request.URL.Path)

		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR",
			"服务器内部错误", recovered)
	})
}

// panicRecoveryWrapper 包装函数，用于捕获panic并转换为错误
func panicRecoveryWrapper(fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("捕获到panic: %v", r)
			err = fmt.Errorf("操作失败，发生内部错误: %v", r)
		}
	}()

	return fn()
}

// panicRecoveryWrapperWithResult 包装函数，用于捕获panic并返回错误结果
func panicRecoveryWrapperWithResult[T any](fn func() (T, error)) (result T, err error) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("捕获到panic: %v", r)
			err = fmt.Errorf("操作失败，发生内部错误: %v", r)
		}
	}()

	return fn()
}
