package main

import (
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"pkmm/models"
	_ "pkmm/utils"
	_ "pkmm/routers"
	"pkmm/controllers"
	"github.com/astaxie/beego/migration"
)

func main() {

	models.Init()

	migration.Upgrade(0)
	beego.Router("/", &controllers.ForumController{}, "*:GetForums")

	beego.Run()
}
