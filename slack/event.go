package slack

import (
	"log"
)

type Event struct {
	Type string
	Ts   string
	Raw  []byte
}

type EventHandler struct {
	Table            map[string]func(Event)
	ExceptionHandler func(Event)
}

func NewEventHandler() *EventHandler {
	table := make(map[string]func(Event))
	table["hello"] = func(event Event) {
		log.Print("Successfully connected")
	}

	return &EventHandler{Table: table}
}

func (this *EventHandler) SetExceptionHandler(handler func(Event)) {
	this.ExceptionHandler = handler
}

func (this *EventHandler) AddHandler(eventType string, handler func(Event)) {
	this.Table[eventType] = handler
}

func (this *EventHandler) Handle(event Event) {
	go func() {
		handler, ok := this.Table[event.Type]
		if ok {
			handler(event)
		} else {
			if this.ExceptionHandler != nil {
				this.ExceptionHandler(event)
			}
		}
	}()
}
