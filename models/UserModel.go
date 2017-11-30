package models

import "github.com/astaxie/beego/orm"

type User struct {
	Id        int
	UserName  string
	Password  string
	Email     string
	Salt      string
	LastLogin int
	Status    int8
	CreatedAt int64
	Bduss     string
}

func (t *User) TableName() string {
	return TableName("user")
}


func (u *User) Update(fields ...string) error {
	if _,err:= orm.NewOrm().Update(u, fields...); err != nil {
		return err
	}
	return nil
}


func UserGetById(id int) (*User, error) {
	u := new(User)
	err := orm.NewOrm().QueryTable(TableName("user")).Filter("id", id).One(u)
	if err != nil {
		return u, err
	}
	return u, nil
}

func UserGetByName(name string) (*User, error) {
	u := new(User)
	err := orm.NewOrm().QueryTable(TableName("user")).Filter("username", name).One(u)
	if err != nil {
		return u, err
	}
	return u, nil
}

func GetAllUsers() ([]*User, int64, error) {
	users := make([]*User, 0)
	total, err := orm.NewOrm().QueryTable(TableName("user")).Exclude("bduss__isnull", true).All(&users)
	return users, total, err
}

func UserAdd(u *User) int64 {
	uid, err := orm.NewOrm().Insert(u)
	if err != nil {
		return -1
	}
	return uid
}

func UserUpdate(user *User, fields ...string) error {
	_, err := orm.NewOrm().Update(user, fields...)
	return err
}
