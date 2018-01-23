package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type Score struct {
	Id        int64
	Kcmc      string
	Cj        string
	Cxcj      string
	Bkcj      string
	Xf        string
	Jd        string
	CreatedAt time.Time
	StuId     int64
	Xn        string
	Xq        string
}

func InsertScores(scores []Score) (int64, error) {
	successNum, err := orm.NewOrm().InsertMulti(len(scores), scores)
	return successNum, err
}

func GetScoresByStuId(stuId int64) []*Score {
	scores := make([]*Score, 0)
	orm.NewOrm().QueryTable("score").Filter("stu_id", stuId).All(&scores)
	return scores
}

func InsertOrUpdateScore(score *Score) {
	o := orm.NewOrm()
	cnt, _ := o.QueryTable("score").Filter("xn", score.Xn).
		Filter("xq", score.Xq).Filter("kcmc", score.Kcmc).Filter("stu_id", score.StuId).Count()
	if cnt != 0 {
		return
	}
	beego.Debug(score)
	o.Insert(score)
}
