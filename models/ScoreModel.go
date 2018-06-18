package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type Score struct {
	Id        int64
	Kcmc      string
	Type      string
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

func GetScoresByStuId(stuId int64) []Score {
	scores := make([]Score, 0)
	orm.NewOrm().QueryTable("score").Filter("stu_id", stuId).All(&scores)
	return scores
}

func InsertOrUpdateScore(score *Score) {
	o := orm.NewOrm()
	cnt, _ := o.QueryTable("score").
		Filter("xn", score.Xn).
		Filter("xq", score.Xq).
		Filter("kcmc", score.Kcmc).
		Filter("type", score.Type).
		Filter("stu_id", score.StuId).
		Count()
	if cnt != 0 {
		return
	}
	beego.Debug(score)
	o.Insert(score)
}

func GetFailedLessons() (error, []orm.Params) {
	var maps []orm.Params
	o := orm.NewOrm()
	_, err := o.Raw("SELECT kcmc, `type`, COUNT(*) AS cnt " +
		"FROM score JOIN (SELECT id FROM score WHERE jd <=0) AS s2 ON s2.id = score.`id` " +
		"GROUP BY score.`kcmc` ORDER BY cnt DESC limit 10").Values(&maps)

	return err, maps
}
