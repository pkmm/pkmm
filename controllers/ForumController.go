package controllers

import (
	"encoding/json"
	"pkmm/models"
)

type ForumController struct {
	BaseController
}

// 对应的类型必须是匹配的否则是解析不到的，不能把string转化为int等
type ReplyJson struct {
	Ctime      int64       `json:"ctime"`
	ErrorCode  string      `json:"error_code"`
	ErrorMsg   string      `json:"error_msg"`
	Info       []string    `json:"info"`
	Logid      int64       `json:"logid"`
	ServerTime string      `json:"server_time"`
	Time       int64       `json:"time"`
	UserInfo   interface{} `json:"user_info"`
}

// 后面加上tag标签，可以使得转化之后的json是小写的
type ResponseJson struct {
	Id        int       `json:"id"`
	Kw        string    `json:"kw"`
	LastSign  int       `json:"last_sign"`
	CreatedAt int64     `json:"created_at"`
	SignInfo  ReplyJson `json:"sign_info"`
}

func prepareResponseData(forums []*models.Forum) *map[int]ResponseJson {
	result := make(map[int]ResponseJson)
	for index, forum := range forums {
		replyJson := ReplyJson{}
		json.Unmarshal([]byte(forum.ReplyJson), &replyJson)
		result[index] = ResponseJson{
			forum.Id,
			forum.Kw,
			forum.LastSign,
			forum.CreatedAt,
			replyJson,
		}
	}
	return &result
}

// 获取一个用户的所有的贴吧
func (this *ForumController) GetForums() {
	userId := this.GetString("userId", "1")
	forums, total := models.GetForumsByUserId(userId)
	out := make(map[string]interface{})
	out["total"] = total
	out["forums"] = prepareResponseData(forums)
	this.jsonResult(out)
}

func (this *ForumController) AddForum() {

}

// @router /:uid([0-9]+) [get]
func (this *ForumController) Get() {
	uid := this.Ctx.Input.Param(":uid")
	forums, _ := models.GetForumsByUserId(uid)
	this.jsonResult(prepareResponseData(forums))
}
