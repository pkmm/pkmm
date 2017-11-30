package controllers

import (
	"pkmm/models"
	"encoding/json"
	"github.com/astaxie/beego/utils"
	myUtils "pkmm/utils"
	"time"
	"fmt"
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

// @router /:uid([0-9]+) [get]
func (this *UserController) GetUser() {
	uid := this.Ctx.Input.Param(":uid")
	fmt.Println(uid)
	user, err := models.UserGetById(uid)
	out := make(map[string]interface{}, 0)
	if err != nil {
		out["msg"] = err
	}
	out["user"] = user
	this.jsonResult(out)
}

// @router /update_bduss/:uid [post]
func (this *UserController) UpdateUser() {
	uid := this.Ctx.Input.Param(":uid")
	bduss := this.GetString("bduss")
	out := make(map[string]interface{}, 0)
	if bduss == "" {
		out["msg"] = MSG_ERR
		this.jsonResult(out)
	}
	user, err := models.UserGetById(uid)
	if err != nil {
		out["msg"] = MSG_ERR
		this.jsonResult(out)
	}
	user.Bduss = this.GetString("bduss")
	if user.Bduss != "" {
		user.Update("bduss")
	}
	out["msg"] = MSG_OK
	this.jsonResult(out)
}
