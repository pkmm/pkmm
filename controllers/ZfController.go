package controllers

import (
	"github.com/astaxie/beego"
	"pkmm/models"
	"pkmm/utils"
	"time"
	"strconv"
	"fmt"
)

type ZfController struct {
	BaseController
}

func (this *ZfController) GetScore() {
	num := this.GetString("num")
	pwd := this.GetString("pwd")
	version := this.GetString("version")
	stu := models.CreatedOrUpdate(num, pwd)
	beego.Debug("num: ", num, "pwd :", pwd)
	// 用户存在(刚才创建或者已经有的
	if stu.Id != 0 {
		scores := models.GetScoresByStuId(stu.Id)
		if len(scores) != 0 {
			if version == "1" {
				this.makeScoreResult(scores)
			} else {
				this.jsonResult(scores)
			}
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
	if version == "1" {
		this.makeScoreResult(scores)
	} else {
		this.jsonResult(scores)
	}
}

func (this *ZfController) makeScoreResult(scores []models.Score) {
	mp := make(map[string]interface{})
	mp["scores"] = scores
	mp["total"] = len(scores)
	var cnt int
	var sum float64
	var year = time.Now().Year()
	XN := fmt.Sprintf("%d-%d", year-1, year) // 学年的绩点
	fmt.Println(XN)
	var sumXN float64
	var cntXN float64
	for _, s := range scores {
		if s.Type == "必修课" {
			val, err := strconv.ParseFloat(s.Jd, 64)
			if err != nil {
				val = 0
			}
			sum += val
			cnt ++

			if s.Xn == XN {
				sumXN += val
				cntXN ++
			}
		}

	}
	mp["total_jd"] = sum
	mp["total_cnt"] = cnt
	var totalAvg float64
	if cnt != 0 {
		totalAvg = sum / float64(cnt)
	}
	mp["total_jd_avg"] = fmt.Sprintf("%.2f", totalAvg)

	mp["current_xn_jd"] = sumXN
	var currentAvg float64
	if cntXN != 0 {
		currentAvg = sumXN / cntXN
	}
	mp["current_xn_jd_avg"] = fmt.Sprintf("%.2f", currentAvg)
	mp["current_xn_cnt"] = cntXN

	mp["cur_xn"] = XN
	this.jsonResult(mp)
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

func (this *ZfController) CheckAccount() {
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

func (this *ZfController) UpdateAccount() {
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
