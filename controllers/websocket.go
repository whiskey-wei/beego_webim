package controllers

import (
	"net/http"

	"webim/models"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
)

type WebsocketController struct {
	beego.Controller
}

func (c *WebsocketController) HandleWs() {
	ws, err := websocket.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(c.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		beego.Error("can not setup websocket connection:", err)
	}
	username := c.GetSession("username")
	if username == nil {
		ws.Close()
	} else if value, ok := username.(string); ok == false {
		ws.Close()
	} else {
		Join(value, ws)
		defer Leave(value)

		for {
			_, msg, err := ws.ReadMessage()
			if err == nil {
				publish <- newEvent(models.EVENT_MESSAGE, value, string(msg))
			}
		}
	}
}
