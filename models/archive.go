package models

import "container/list"

type EvenType int

//消息类型
const (
	EVENT_JOIN = iota
	EVENT_LEVEA
	EVENT_MESSAGE
)

type Event struct {
	Type      EvenType
	User      string
	Timestamp int //time
	Content   string
}

const archiveSize = 50 //聊天记录

var archive = list.New() //事件队列

func NewArchive(event Event) {
	for archive.Len() >= archiveSize {
		archive.Remove(archive.Front())
	}
	archive.PushBack(event)
}

//根据传过来的时间戳返回时间戳之后的所有时间消息
func GetEvents(lastReceived int) []Event {
	events := make([]Event, 0, archive.Len())
	for event := archive.Front(); event != nil; event = event.Next() {
		e := event.Value.(Event)
		if e.Timestamp > int(lastReceived) {
			events = append(events, e)
		}
	}
	return events
}
