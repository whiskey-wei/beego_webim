package routers

import (
	"webim/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.UserController{}, `get:PageIndex`)
	beego.Router("/signin", &controllers.UserController{}, `get:PageSignin`)
	beego.Router("/signup", &controllers.UserController{}, `get:PageSignup`)

	beego.Router("/signin", &controllers.UserController{}, `post:Signin`)
	beego.Router("/signup", &controllers.UserController{}, `post:Signup`)

	beego.Router("/chat", &controllers.UserController{}, `get:PageChat`)

	beego.Router("/ws", &controllers.WebsocketController{}, `get:HandleWs`)
}
