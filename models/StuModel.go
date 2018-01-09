package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type Stu struct {
	Id        int64
	Name      string
	Num       string
	Pwd       string
	CreatedAt time.Time
	WxUid     string
}

func CreatedOrUpdate(num, pwd string) *Stu {
	if num == "" || pwd == "" {
		return &Stu{}
	}
	o := orm.NewOrm()
	stu := Stu{}
	err := o.QueryTable("stu").Filter("num", num).One(&stu)
	if err == nil {
		stu.Num = num
		stu.Pwd = pwd
		o.Update(&stu, "num", "pwd")
	} else {
		stu.CreatedAt = time.Now()
		o.Insert(&stu)
	}
	return &stu
}

func GetPwdByNum(num string) string {
	o := orm.NewOrm()
	var stu Stu
	o.QueryTable("stu").Filter("num", num).One(&stu)
	return stu.Pwd
}

func GetStuByWxUid(WxUid string) *Stu {
	o := orm.NewOrm()
	var stu Stu
	o.QueryTable("stu").Filter("wx_uid", WxUid).One(&stu)
	return &stu
}
