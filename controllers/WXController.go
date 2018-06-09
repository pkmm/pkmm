package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"net/http"
)

type WXController struct {
	BaseController
}

func (this *WXController) Login() {
	code := this.GetString("code")
	client := &http.Client{}
	appId := beego.AppConfig.String("wx.appId")
	secret := beego.AppConfig.String("wx.secret")
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"https://api.weixin.qq.com/sns/jscode2session?appId=%s&secret=%s&js_code=%s&grant_type=authorization_code",
			appId,
			secret,
			code,
		),
		nil,
	)

	out := make(map[string]interface{})

	if err != nil {
		out["msg"] = err
		this.jsonResult(out)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		out["msg"] = err
		this.jsonResult(out)
		return
	}
	defer resp.Body.Close()
	jsonResult, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		out["msg"] = err
		this.jsonResult(out)
		return
	}
	this.Ctx.WriteString(string(jsonResult))
}
