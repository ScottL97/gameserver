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

func (c *UserController) LoginUser() {
	var user models.User
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &user)
	if err != nil {
		fmt.Println("[LoginUser]json.Unmarshal is err:", err.Error())
	}
	fmt.Println(user)
	// 用户登录
	id, err := loginUser(user)

	c.Ctx.WriteString(id)
}

func (c *UserController) CheckUser() {
	var userCheck models.UserCheck
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &userCheck)
	if err != nil {
		fmt.Println("[CheckUser]json.Unmarshal is err:", err.Error())
	}
	fmt.Println(userCheck)
	getUserManager().loginUsersMutex.Lock()
	defer getUserManager().loginUsersMutex.Unlock()
	if userId, ok := getUserManager().loginUsers[userCheck.Username]; ok {
		if userId == userCheck.Id {
			c.Ctx.WriteString("ok")
			return
		}
	}
	c.Ctx.WriteString("fail")
}

func (c *UserController) Get() {
	c.TplName = "login.tpl"
}
