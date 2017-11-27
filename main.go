package main

import (
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"pkmm/models"
	_ "pkmm/utils"
	"pkmm/controllers"
)

func main() {

	models.Init()

	beego.Router("/", &controllers.ForumController{}, "*:GetForums")

	beego.Run()
}
