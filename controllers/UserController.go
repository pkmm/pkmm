package controllers

import (
	"pkmm/models"
	"encoding/json"
	"github.com/astaxie/beego/utils"
	myUtils "pkmm/utils"
	"time"
)

type UserController struct {
	BaseController
}

// @router /register [post]
func (this *UserController) UserRegister() {
	var user models.User
	// 以json格式传递数据
	json.Unmarshal(this.Ctx.Input.RequestBody, &user)

	user.Salt = string(utils.RandomCreateBytes(10))
	user.Password = myUtils.Md5([]byte(user.Salt + "pkmm" + user.Password))
	user.CreatedAt = time.Now().Unix()

	uid := models.UserAdd(&user)
	out := make(map[string]interface{}, 0)
	out["user"] = user
	out["uid"] = uid
	this.jsonResult(out)
}
