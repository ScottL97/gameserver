package routers

import (
	"gameserver/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/game", &controllers.MainController{})
	beego.Router("/ws", &controllers.WebSocketController{})
	beego.Router("/user", &controllers.UserController{}, "post:LoginUser")
	beego.Router("/checkuser", &controllers.UserController{}, "post:CheckUser")
	beego.Router("/", &controllers.UserController{}, "get:Get")
	beego.SetStaticPath("/node_modules", "node_modules")
}
