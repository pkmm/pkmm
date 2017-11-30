package controllers

import (
	"pkmm/models"
)

type ForumController struct {
	BaseController
}

// 获取一个用户的所有的贴吧
func (this *ForumController) GetForums() {
	userId := this.GetString("userId", "1")
	forums, total := models.GetForumsByUserId(userId)
	out := make(map[string]interface{})
	out["total"] = total
	out["forums"] = forums
	this.jsonResult(out)
}

func (this *ForumController) AddForum() {

}

// @router /:uid([0-9]+) [get]
func (this *ForumController) Get() {
	uid := this.Ctx.Input.Param(":uid")
	forums, _ := models.GetForumsByUserId(uid)
	this.jsonResult(forums)
}
