package controllers

import (
	"pkmm/models"
)

type ForumController struct {
	BaseController
}

// 获取一个用户的所有的贴吧

func (this *ForumController) GetForums() {

	userId, _ := this.GetInt("userId", 1)
	forums, total := models.GetForumsByUserId(userId)
	out := make(map[string]interface{})
	out["total"] = total
	out["forums"] = forums
	this.jsonResult(out)
}

func (this *ForumController) AddForum() {

}