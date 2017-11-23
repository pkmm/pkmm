package main

import (
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"pkmm/models"
	"pkmm/controllers"
	_ "pkmm/utils"
)

func main() {

	models.Init()

	beego.Router("/", &controllers.ForumController{}, "*:GetForums")

	beego.Run()
}
