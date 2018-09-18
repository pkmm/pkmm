package main

import (
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	_ "pkmm/database/migrations"
	"pkmm/models"
	_ "pkmm/routers"
	"pkmm/utils"
	"fmt"
	"time"
)

func main() {

	// 初始化model
	models.Init()

	// 数据库迁移
	//migration.Upgrade(0)

	// 设置日志
	beego.SetLogger("file", `{"filename":"logs/pkmm.log","level":7,"daily":true,"maxdays":2}`)

	// 设置静态资源文件, eg. /static/images/xx.png
	beego.SetStaticPath("/static", "static")

	// 部署Email提醒
	runMode := beego.AppConfig.String("runmode")
	if (runMode != "dev") {
		utils.SendMail(
			"xiaoccla@qq.com",
			"xiaoccla",
			"部署HOOK通知",
			fmt.Sprintf(
				"pkmm代码已经重新部署, Time: [%s], IP: [%s]",
				beego.Date(time.Now(), "Y-m-d H:i:s"),
				utils.IpAddressOfLocal(),
			),
			[]string{},
		)
	}

	beego.Run()
}
