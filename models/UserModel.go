package models

import "github.com/astaxie/beego/orm"

type User struct {
	Id        int    `orm:"column(id);auto"`
	UserName  string `orm:"column(user_name);size(40)" description:"用户名"`
	Password  string `orm:"column(password);size(60)" description:"密码"`
	Email     string `orm:"column(email);size(30)"`
	Salt      string `orm:"column(salt);size(10)" description:"密码盐"`
	LastLogin int    `orm:"column(last_login)" description:"最后登录的时间"`
	Status    int8   `orm:"column(status)" description:"状态, 0 正常 -1 禁用"`
	CreatedAt int64  `orm:"column(created_at)" description:"创建时间"`
	Bduss     string
}

func (t *User) TableName() string {
	return TableName("user")
}

func (this *User) GetUserById(id int) (*User, error) {
	user := &User{
		Id: id,
	}
	err := orm.NewOrm().Read(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetAllUsers() ([]*User, int64, error){
	users := make([]*User, 0)
	total, err := orm.NewOrm().QueryTable(TableName("user")).Exclude("bduss__isnull", true).All(&users)
	return users, total, err
}
