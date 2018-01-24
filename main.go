package main

import (
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"pkmm/controllers"
	_ "pkmm/database/migrations"
	"pkmm/models"
	_ "pkmm/routers"
)

func main() {

	// 初始化model
	models.Init()

	// 数据库迁移
	//migration.Upgrade(0)
	beego.Router("/", &controllers.ForumController{}, "*:GetForums")
	beego.Router("/zf", &controllers.ZfController{}, "*:GetScore")
	beego.Router("/zf/check_account", &controllers.ZfController{}, "*:CheckAccount")
	beego.Router("/zf/update_account", &controllers.ZfController{}, "post:UpdateAccount")

	beego.SetLogger("file", `{"filename":"logs/pkmm.log","level":7,"daily":true,"maxdays":2}`)
	beego.Run()
}
