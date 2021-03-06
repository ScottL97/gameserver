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
	beego.Router("/gamectrl/ready/:username", &controllers.GameController{}, "get:UserReady")
	beego.Router("/gamectrl/cancel/:username", &controllers.GameController{}, "get:UserCancel")
	beego.Router("/gamectrl/finish/:username", &controllers.GameController{}, "get:UserFinish")
	beego.Router("/gamectrl/killvirus", &controllers.GameController{}, "post:KillVirus")
	beego.Router("/gamectrl/move", &controllers.GameController{}, "post:UserMove")
	beego.Router("/gamectrl/buildinstitute", &controllers.GameController{}, "post:BuildInstitute")
	beego.Router("/gamectrl/doresearch", &controllers.GameController{}, "post:DoResearch")
	beego.Router("/", &controllers.UserController{}, "get:Get")
	beego.SetStaticPath("/node_modules", "node_modules")
}
