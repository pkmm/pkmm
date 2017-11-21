package models

import (
	"github.com/astaxie/beego"
	"fmt"
	"net/url"
	"github.com/astaxie/beego/orm"
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

	orm.RegisterModel(new(Forum), new(User))
	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	}
}


// 方便给表加前缀
func TableName(name string) string {
	return beego.AppConfig.String("db.prefix") + name
}
