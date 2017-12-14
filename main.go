package main

import (
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"pkmm/models"
	_ "pkmm/utils"
	_ "pkmm/routers"
	"pkmm/controllers"
	"github.com/astaxie/beego/migration"
	_ "pkmm/database/migrations"
)

func main() {

	// 初始化model
	models.Init()

	// 数据库迁移
	migration.Upgrade(0)

	beego.Router("/", &controllers.ForumController{}, "*:GetForums")

	beego.Run()
}
