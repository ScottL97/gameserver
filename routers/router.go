package routers

import (
	"gameserver/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/ws", &controllers.WebSocketController{})
	beego.Router("/user", &controllers.UserController{})
	beego.SetStaticPath("/node_modules", "node_modules")
}
