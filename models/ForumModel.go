package models

import (
	"github.com/astaxie/beego/orm"
	"fmt"
	"time"
)

type Forum struct {
	Id         int
	UserId     int
	Kw         string
	Fid        int
	LastSign   int
	SignStatus int
	CreatedAt  int64
	ReplyJson  string
}

// 每一个表需要实现 这个方法 才能使模型被正确的注册
// 否则默认的话 模型的名称必须与表名称一样
func (t *Forum) TableName() string {
	return TableName("forums")
}

func GetForumsByUserId(userId string) ([]*Forum, int64) {
	forums := make([]*Forum, 0)
	total, err := orm.NewOrm().QueryTable(TableName("forums")).
		Filter("user_id", userId).
		All(&forums)
	if err != nil {
		return nil, 0
	}
	return forums, total
}

func NeedSignForumsByUserId(userId int) ([]*Forum, int64) {
	forums := make([]*Forum, 0)
	total, err := orm.NewOrm().QueryTable(TableName("forums")).
		Exclude("last_sign", time.Now().Day()).
		Filter("user_id", userId).
		Exclude("fid", -1).
		All(&forums)
	if err != nil {
		return nil, 0
	}
	return forums, total
}

func AddForum(forum *Forum) (int64, error) {
	if forum.UserId == 0 {
		return 0, fmt.Errorf("用户的userid不能为空")
	}
	if forum.Fid == 0 {
		return 0, fmt.Errorf("贴吧的fid不能为空")
	}
	if forum.Kw == "" {
		return 0, fmt.Errorf("贴吧的kw不能为空")
	}
	if forum.LastSign == 0 {
		forum.LastSign = -1
	}
	if forum.CreatedAt == 0 {
		forum.CreatedAt = time.Now().Unix()
	}
	return orm.NewOrm().InsertOrUpdate(forum)
}

func AllForums() ([]*Forum, int64) {
	forums := make([]*Forum, 0)
	total, err := orm.NewOrm().QueryTable(TableName("forums")).All(&forums)
	if err != nil {
		return nil, 0
	}
	return forums, total
}

func (f *Forum) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(f, fields...); err != nil {
		return err
	}
	return nil
}
