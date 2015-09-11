package main

import (
    "log"
    "encoding/json"
    "os"
    "time"
    "strings"
    "strconv"
    "./slack"
)

const (
    BOT_NAME = "imghr"
    IHR_ID = "U037GMSJ9"
)

type UserTyping struct {
    slack.Event
    Channel string
    User string
}

type EventHandler struct {
    Table map[string] func(slack.Event)
    ExceptionHandler func(slack.Event)
}

func NewEventHandler() *EventHandler {
    table := make(map[string] func(slack.Event))
    table["hello"] = func(event slack.Event) {
        log.Print("Successfully connected")
    }

    return &EventHandler{Table: table}
}

func (this *EventHandler) SetExceptionHandler(handler func(slack.Event)) {
    this.ExceptionHandler = handler
}

func (this *EventHandler) AddHandler(eventType string, handler func(slack.Event)) {
    this.Table[eventType] = handler
}

func (this *EventHandler) Handle(event slack.Event) {
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

type UserTypingEventHandler struct {
    PostFlag bool
    Time time.Time
    FireChan chan slack.Event
}

func NewUserTypingEventHandler() *UserTypingEventHandler {
    defaultTime := time.Now().Add(time.Duration(-13)*time.Hour)
    fireChan := make(chan slack.Event)

    this := &UserTypingEventHandler{PostFlag: false, Time: defaultTime, FireChan: fireChan}

    go func() {
        for {
            select {
            case event := <-this.FireChan:
                if this.IsEnable() == true {
                    return
                }
                var userTyping UserTyping
                json.Unmarshal(event.Raw, &userTyping)

                if userTyping.User != IHR_ID {
                    return
                }

                token := os.Getenv("SLACK_TOKEN")
                this.PostFlag = true
                this.Time = time.Now()
                slack.PostMessage(token, userTyping.Channel, BOT_NAME, "I H R は 寝 て ろ ！ ！")
            }
        }
    }()

    return this
}

func (this *UserTypingEventHandler) IsEnable() bool {
    // 12時間以上経っていたらフラグをリセット
    if time.Now().Unix() - this.Time.Unix() > 43200 {
        this.PostFlag = false
    }
    if this.PostFlag == true {
        return false
    }

    if time.Now().Hour() <= 14 && time.Now().Hour() > 17 {
        return true
    } else {
        return false
    }
}

func (this *UserTypingEventHandler) Handle(event slack.Event) {
    this.FireChan <- event
}

func main() {
    token := os.Getenv("SLACK_TOKEN")
    if token == "" {
        log.Fatal("not set SLACK_TOKEN")
    }
    log.Print("init")
    eventHandler := NewEventHandler()
    eventHandler.SetExceptionHandler(func (event slack.Event) {
        log.Print("Unknown Event: ", event.Type)
    })
    messageEventHandler := slack.NewMessageEventHandler()
    userTypingEventHandler := NewUserTypingEventHandler()
    eventHandler.AddHandler("message", messageEventHandler.Handle)
    eventHandler.AddHandler("user_typing", userTypingEventHandler.Handle)
again:
    ws := slack.ConnectSocket(token)
    eventChan, resultChan := slack.StartReading(ws)
    startTime := int(time.Now().Unix())
    for {
        select {
        case buf := <-eventChan:
            var event slack.Event
            json.Unmarshal(buf, &event)
            if event.Ts != "" {
                ts := strings.Split(event.Ts, ".")[0]
                i, _ := strconv.Atoi(ts)
                if i < startTime {
                    log.Print("skip event")
                    continue
                }
            }
            event.Raw = buf
            eventHandler.Handle(event)
        case <-resultChan:
            log.Print("disconnect server")
            goto again
        }
    }
}
