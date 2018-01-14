package controllers

import (
	"pkmm/utils/zf"
	"pkmm/models"
	"github.com/astaxie/beego"
)

type ZfController struct {
	BaseController
}

func (this *ZfController) Get() {
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
	ret, _ := zf.Login(num, pwd)
	if len(ret) != 0 && stu.Id != 0 {
		// mock data
		//ret = [][]string{{"2013-2014", "1", "数学", "22","23", "34", "34", "234"}, {"2013-2014", "1", "数学222", "22","23", "34", "34", "234"}}
		_, err := models.InsertScores(ret, stu.Id)
		beego.Debug(err)
	}
	this.jsonResult(ret)
}