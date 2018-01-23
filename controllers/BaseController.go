package controllers

import "github.com/astaxie/beego"

const (
	MSG_OK  = 1
	MSG_ERR = -1
)

type BaseController struct {
	beego.Controller
	controllerName string
	actionName     string
}

//开启跨域
func (c BaseController) Prepare() {
	c.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
}

// 输出json
func (this *BaseController) jsonResult(out interface{}) {
	this.Data["json"] = out
	this.ServeJSON()
	this.StopRun()
}

// {msg_status:ok/false, response: ... }
//
