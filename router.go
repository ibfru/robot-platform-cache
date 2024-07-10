package main

import (
	"fmt"
	sdk "git-platform-sdk"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type header struct {
	Token string `form:"token"`
}

const (
	AuthErrorMessage = "authorized failed"
)

func authRequired(c *gin.Context) {
	var h header

	// error happened, gin framework handler
	_ = c.BindHeader(&h)

	if h.Token == "111" {
		return
	}

	c.AbortWithStatusJSON(http.StatusUnauthorized, AuthErrorMessage)
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	println(gin.Mode())
	if gin.Mode() == gin.ReleaseMode {
		gin.DisableConsoleColor()
	}

	router := gin.New()
	// LoggerWithFormatter 中间件会写入日志到 gin.DefaultWriter
	// 默认 gin.DefaultWriter = os.Stdout
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 你的自定义格式
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.DateTime),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	router.Use(gin.Recovery())

	// 认证路由组
	controller := Controller{
		group:  router.Group("/", authRequired),
		client: sdk.GetClientInstance(""),
	}
	controller.registerRoutePath()

	return router
}
