package controllers

import (
	"pkmm/utils/zf"
	"fmt"
)

type ZfController struct {
	BaseController
}

func (this *ZfController) Post() {
	num := this.GetString("num")
	pwd := this.GetString("pwd")
	out := make(map[string]string)
	out["num"] = num
	out["pwd"] = pwd


	ret:= zf.Login(num, pwd)
	fmt.Println(out)
	this.jsonResult(ret)
}
