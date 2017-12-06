package controllers

import (
	"pkmm/models"
	"encoding/json"
)

type ForumController struct {
	BaseController
}

type ReplyJson struct {
	Ctime      string      `json:"ctime"`
	ErrorCode  string      `json:"error_code"`
	ErrorMsg   string      `json:"error_msg"`
	Info       []string    `json:"info"`
	Logid      string      `json:"logid"`
	ServerTime string      `json:"server_time"`
	Time       string      `json:"time"`
	UserInfo   interface{} `json:"user_info"`
}
type ResponseJson struct {
	Id        int
	Kw        string
	LastSign  int
	CreatedAt int64
	SignInfo  ReplyJson
}

// 获取一个用户的所有的贴吧
func (this *ForumController) GetForums() {
	userId := this.GetString("userId", "1")
	forums, total := models.GetForumsByUserId(userId)
	out := make(map[string]interface{})
	var replyJson ReplyJson
	result := make(map[int]ResponseJson)
	for index, forum := range forums {
		json.Unmarshal([]byte(forum.ReplyJson), &replyJson)
		result[index] = ResponseJson{
			forum.Id,
			forum.Kw,
			forum.LastSign,
			forum.CreatedAt,
			replyJson,
		}
	}

	out["total"] = total
	out["forums"] = result
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
