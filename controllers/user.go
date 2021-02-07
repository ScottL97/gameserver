package controllers

import (
	"encoding/json"
	"fmt"
	"gameserver/models"
	"github.com/astaxie/beego"
)

type UserController struct {
	beego.Controller
}

func (c *UserController) Post() {
	var user models.User
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &user)
	if err != nil {
		fmt.Println("json.Unmarshal is err:", err.Error())
	}
	fmt.Println(user)
	// 用户登录
	id, err := loginUser(user)

	c.Ctx.WriteString(id)
}
