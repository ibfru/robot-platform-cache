package main

import (
	"encoding/json"
	"fmt"
	sdk "git-platform-sdk"
	"io"
	"net/http"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

//var Controller *gin.RouterGroup

type Controller struct {
	group  *gin.RouterGroup
	client *sdk.ClientTarget
}

func (ctl *Controller) registerRoutePath() {

	g := ctl.group
	g.Use()

	// Ping test
	g.POST("/atomgit-hook", func(c *gin.Context) {
		fmt.Printf("%+v\n", c.Request.Header)
		var b map[string]any
		b1, _ := io.ReadAll(c.Request.Body)
		_ = json.Unmarshal(b1, &b)
		fmt.Printf("%+v\n", b)
		c.String(http.StatusOK, "pong")
	})

	// 嵌套路由组
	testing := g.Group("door-control")
	testing.GET("/1", func(c *gin.Context) {
		err := loadCacheFormGitPlatform("")
		if err != nil {
			c.String(http.StatusOK, "no data")
		} else {
			c.Data(http.StatusOK, binding.MIMEMultipartPOSTForm, DoorControlCache.Get(nil, []byte("ibforuorg/community-test/raw/sig/Test/sig-info.yaml")))
		}

	})

	testing.GET("/0", func(c *gin.Context) {
		flushCache()
		c.Data(http.StatusOK, binding.MIMEMultipartPOSTForm, DoorControlCache.Get(nil, []byte("conflictCheck.py")))
	})

	ctl.openGitPlatformApi(g.Group("/git-platform"))
}

var orgValid validator.Func = func(fl validator.FieldLevel) bool {
	val, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	l := len(val)
	if l == 0 || l > 16 {
		return false
	}

	return true
}

type prBody sdk.PRParameter

func (ctl *Controller) openGitPlatformApi(g *gin.RouterGroup) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("orgValid", orgValid)
	}

	repo := g.Group("/repos")
	repo.GET("/1", func(context *gin.Context) {
		context.String(http.StatusOK, "repos")
	})
	pr := g.Group("/pr")
	pr.GET("/1", func(context *gin.Context) {
		var b prBody
		_ = context.ShouldBindBodyWithJSON(&b)
		context.String(http.StatusOK, "pull request")
	})
	pr.POST("add-comment", func(context *gin.Context) {
		var b sdk.PRParameter
		//q, _ := io.ReadAll(context.Request.Body)
		//if err := json.Unmarshal(q, &b); err == nil {
		//	context.String(http.StatusOK, "111, "+b.Comment+", "+b.Org)
		//} else {
		//	context.String(http.StatusBadRequest, "111, "+b.Comment+", "+b.Org+", err: "+err.Error())
		//}
		if err := context.ShouldBindBodyWithJSON(&b); err == nil {

			err = ctl.client.AddPRComment(&b)
			if err != nil {
				context.String(http.StatusBadRequest, "err, "+err.Error())
			} else {
				context.String(http.StatusOK, "add pull request successful")
			}

		} else {
			context.String(http.StatusBadRequest, "json args parse failed: "+err.Error())
		}

	})
	issue := g.Group("issue")
	issue.GET("/1", func(context *gin.Context) {
		context.String(http.StatusOK, "issue")
	})
	label := g.Group("label")
	label.GET("/1", func(context *gin.Context) {
		context.String(http.StatusOK, "label")
	})
}
