package bot

import (
	"time"
)

type Connector interface {
	Connect()
	Listen() error
	ReceivedEvent() chan *Event
	Send(*Event, string, string) error
}

const (
	MessageEvent    = "message"
	UserTypingEvent = "user_typing"
	UnknownEvent    = "unknown"
)

type Event struct {
	Type        string
	Message     string
	Argv        string
	Channel     string
	User        string
	Timestamp   time.Time
	MentionName string
	Bot         *Bot
}