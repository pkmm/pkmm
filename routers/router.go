// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"pkmm/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",

		beego.NSNamespace("/forum",
			beego.NSInclude(
				&controllers.ForumController{},
			),
		),

		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
	)
	beego.AddNamespace(ns)

	// 指定路由
	beego.Router("/", &controllers.ForumController{}, "*:GetForums")
	beego.Router("/zf", &controllers.ZfController{}, "*:GetScore")
	beego.Router("/zf/check_account", &controllers.ZfController{}, "*:CheckAccount")
	beego.Router("/zf/update_account", &controllers.ZfController{}, "post:UpdateAccount")
	beego.Router("/zf/get_failed_lessons", &controllers.ZfController{}, "*:GetFailedLessons")

	// 自动路由
	beego.AutoRouter(&controllers.WXController{})


}
