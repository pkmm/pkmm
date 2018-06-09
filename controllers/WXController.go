package controllers

import (
	"net/http"
	"net/url"
	"strings"
	"io/ioutil"
)

type WXController struct {
	BaseController
}

func (this *WXController) Login() {
	code := this.Ctx.Input.Param("code")
	client := &http.Client{}
	params := make(url.Values)
	params.Set("appid", "wx52fe76e5731453e0")
	params.Set("secret", "9fb9be904c1edb5a17b620a0f97f7dba")
	params.Set("js_code", code)
	params.Set("grant_type", "authorization_code")
	req, err := http.NewRequest(
		"GET",
		"https://api.weixin.qq.com/sns/jscode2session",
		strings.NewReader(params.Encode()),
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
