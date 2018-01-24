package controllers

import (
	"github.com/astaxie/beego"
	"pkmm/models"
	"pkmm/utils"
	"time"
)

type ZfController struct {
	BaseController
}

func (this *ZfController) GetScore() {
	num := this.GetString("num")
	pwd := this.GetString("pwd")
	stu := models.CreatedOrUpdate(num, pwd)
	beego.Debug("num: ", num, "pwd :", pwd)
	// 用户存在(刚才创建或者已经有的
	if stu.Id != 0 {
		scores := models.GetScoresByStuId(stu.Id)
		if len(scores) != 0 {
			this.jsonResult(scores)
			return
		}
	}
	// 新来的， 需要添加
	crawl := utils.NewCrawl(num, pwd)
	scores, _ := crawl.Login()
	if len(scores) != 0 && stu.Id != 0 {
		preInsertScores := make([]models.Score, 0)
		for _, score := range scores {
			score.StuId = stu.Id
			score.CreatedAt = time.Now()
			preInsertScores = append(preInsertScores, score)
		}
		_, err := models.InsertScores(preInsertScores)
		beego.Debug(err)
	}
	this.jsonResult(scores)
}

func (this *ZfController) Post() {
	num := this.GetString("num")
	pwd := this.GetString("pwd")
	if num == "" || pwd == "" {
		this.jsonResult(map[string]string{"status": "error"})
	}
	this.Ctx.SetCookie("num", num)
	this.Ctx.SetCookie("pwd", pwd)

	this.jsonResult(map[string]string{"status": "success"})
}

func (this *ZfController) CheckAccount(){
	num := this.GetString("num")
	pwd := this.GetString("pwd")
	ok, reason := utils.ValidAccount(num, pwd)
	mp := make(map[string]interface{}, 0)
	if ok {
		mp["msg_status"] = MSG_OK
	} else {
		mp["msg_status"] = MSG_ERR
	}
	mp["response"] = reason
	this.jsonResult(mp)
}

func (this *ZfController) UpdateAccount(){
	num := this.GetString("num")
	pwd := this.GetString("pwd")
	stu := models.CreatedOrUpdate(num, pwd)
	out := make(map[string]interface{}, 0)
	if stu.Id != 0 {
		out["msg_status"] = MSG_OK
		out["response"] = "更新成功"
	} else {
		out["msg_status"] = MSG_ERR
		out["msg_status"] = "更新失败"
	}
	this.jsonResult(out)
}
