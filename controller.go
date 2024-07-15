package main

import (
	sdk "git-platform-sdk"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

//var Controller *gin.RouterGroup

func (bot *robot) registerRoutePath() {
	bot.handlerSigEvent()
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

type SigReqArgs struct {
	Org  string `form:"org"`
	Repo string `form:"repo"`
}

func (bot *robot) handlerSigEvent() {
	bot.ctl.GET("sig/name", func(context *gin.Context) {
		var arg SigReqArgs
		// 如果是 `GET` 请求，只使用 `Form` 绑定引擎（`query`）

		if err := context.ShouldBind(&arg); err != nil {
			context.String(http.StatusBadRequest, "err, "+err.Error())
		} else {
			name := getSigName(arg.Org, arg.Repo)
			context.String(http.StatusOK, string(name))
		}
	})

	bot.ctl.GET("sig/maintainers", func(context *gin.Context) {
		var arg SigReqArgs
		// 如果是 `GET` 请求，只使用 `Form` 绑定引擎（`query`）

		if err := context.ShouldBind(&arg); err != nil {
			context.String(http.StatusBadRequest, "err, "+err.Error())
		} else {
			d := getMaintainers(arg.Org, arg.Repo)
			j, _ := ConvertFromBytes(d)
			context.JSON(http.StatusOK, j)
		}
	})

	bot.ctl.GET("sig/committers", func(context *gin.Context) {
		var arg SigReqArgs
		// 如果是 `GET` 请求，只使用 `Form` 绑定引擎（`query`）

		if err := context.ShouldBind(&arg); err != nil {
			context.String(http.StatusBadRequest, "err, "+err.Error())
		} else {
			d := getCommitters(arg.Org, arg.Repo)
			j, _ := ConvertFromBytes(d)
			context.JSON(http.StatusOK, j)
		}
	})

	bot.ctl.GET("sig/file", func(context *gin.Context) {

		context.JSON(http.StatusOK, bot.getFile())
	})
}
