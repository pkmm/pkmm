package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"net/url"
)

func Init() {
	host := beego.AppConfig.String("db.host")
	port := beego.AppConfig.String("db.port")
	user := beego.AppConfig.String("db.user")
	pwd := beego.AppConfig.String("db.password")
	name := beego.AppConfig.String("db.name")
	timezone := beego.AppConfig.String("db.timezone")

	if port == "" {
		port = "3306"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", user, pwd, host, port, name)

	if timezone != "" {
		dsn = dsn + "&loc=" + url.QueryEscape(timezone)
	}

	orm.RegisterDataBase("default", "mysql", dsn)

	// 模型需要在这里注册
	orm.RegisterModel(
		new(Forum),
		new(User),
		new(Stu), new(Score))

	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	}
}

// 方便给表加前缀
func TableName(name string) string {
	return beego.AppConfig.String("db.prefix") + name
}
