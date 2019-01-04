package controllers

import (
	"container/list"
	"encoding/json"
	"time"
	"webim/models"

	"github.com/astaxie/beego"

	"github.com/gorilla/websocket"
)

type Subscriber struct {
	Name string
	Conn *websocket.Conn
}

type Subscription struct {
	Archive []models.Event
	New     <-chan models.Event
}

var (
	subscribe = make(chan Subscriber, 30) //channel for new users

	unsubscribe = make(chan string, 30) //channel for exit users

	publish = make(chan models.Event, 40) //send event here to publish them

	subscribers = list.New()
)

func newEvent(ep models.EvenType, user, msg string) models.Event {
	return models.Event{Type: ep, User: user, Timestamp: int(time.Now().Unix()), Content: msg}
}

func Join(user string, ws *websocket.Conn) {
	subscribe <- Subscriber{Name: user, Conn: ws}
}

func Leave(user string) {
	unsubscribe <- user
}

func isUserExit(subscribers *list.List, user string) bool {
	for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
		if sub.Value.(Subscriber).Name == user {
			return true
		}
	}
	return false
}

//change message to json and send json to browse
func broadcastWebsocket(event models.Event) {
	data, err := json.Marshal(event)
	if err != nil {
		beego.Error("Fail to marshal event", err)
		return
	}

	for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
		ws := sub.Value.(Subscriber).Conn
		if ws != nil {
			if ws.WriteMessage(websocket.TextMessage, data) != nil {
				unsubscribe <- sub.Value.(Subscriber).Name
			}
		}
	}
}

//handle chan messages
func chatroom() {
	for {
		select {
		case sub := <-subscribe:
			if !isUserExit(subscribers, sub.Name) {
				subscribers.PushBack(sub) //if not exit add user end of list
				publish <- newEvent(models.EVENT_JOIN, sub.Name, "")
				beego.Info("New User:", sub.Name, ";WenSocket:", sub.Conn != nil)
			} else {
				beego.Info("old user:", sub.Name, "WebSocket:", sub.Conn != nil)
			}
		case event := <-publish:
			broadcastWebsocket(event)
			models.NewArchive(event)
			if event.Type == models.EVENT_MESSAGE {
				beego.Info("Message from", event.User, ";Content", event.Content)
			}
		case unsub := <-unsubscribe:
			for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
				if sub.Value.(Subscriber).Name == unsub {
					subscribers.Remove(sub)
					ws := sub.Value.(Subscriber).Conn
					if ws != nil {
						ws.Close()
						beego.Error("WebSocket closed:", unsub)
					}
					publish <- newEvent(models.EVENT_LEVEA, unsub, "")
					break
				}
			}
		}
	}
}

func init() {
	go chatroom()
}
