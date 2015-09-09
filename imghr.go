package main

import (
    "log"
    "encoding/json"
    "os"
    "time"
    "regexp"
    "strings"
    "strconv"
    "./amesh"
    "./jma"
    "./slack"
    "./google"
)

const BOT_NAME = "imghr"

type Event struct {
    Type string
    Ts string
    Raw []byte
}

type Message struct {
    Event
    Channel string
    User string
    Text string
    ts string
}

type EventHandler struct {
    Table map[string] func(Event)
    ExceptionHandler func(Event)
}

func NewEventHandler() *EventHandler {
    table := make(map[string] func(Event))
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

func IsBotComandAlias(message string) bool {
    matched, _ := regexp.MatchString("^(n)$", message)
    return matched
}

func IsBotCommand(message string) bool {
    matched, _ := regexp.MatchString("^imghr\\s+", message)
    return matched
}

func ParseCommand(message string) (string, string) {
    re := regexp.MustCompile("^imghr\\s+(\\w+)(?:\\s+(.+))*")
    matched := re.FindStringSubmatch(message)
    if len(matched) == 2 {
        return matched[1], ""
    }
    return matched[1], matched[2]
}

type MessageEventHandler struct {
    AmeshImageGenerator *amesh.AmeshImageGenerator
    JmaImageGenerator *jma.JmaImageGenerator
}

func NewMessageEventHandler() *MessageEventHandler {
    ameshImageGenerator := amesh.NewAmeshImageGenerator()
    jmaImageGenerator := jma.NewJmaImageGenerator()
    return &MessageEventHandler{AmeshImageGenerator: ameshImageGenerator, JmaImageGenerator: jmaImageGenerator}
}

func (this *MessageEventHandler) ExecuteCommand(message Message, command string, argv string) {
    token := os.Getenv("SLACK_TOKEN")
    switch command {
    case "ping":
        slack.PostMessage(token, message.Channel, BOT_NAME, "pong")
    case "img":
        url := google.ImageSearch(argv)
        if url == "" {
            slack.PostMessage(token, message.Channel, BOT_NAME, "( ˘ω˘ )ｽﾔｧ")
        } else {
            slack.PostMessage(token, message.Channel, BOT_NAME, url+"#.png")
        }
    case "amesh":
        targetDate := time.Now().Add(time.Duration(-1)*time.Minute).Truncate(5 * time.Minute).Format("200601021504")
        imgPath := this.AmeshImageGenerator.Generate(targetDate)
        if imgPath == "" {
            time.Sleep(1 * time.Second)
            slack.PostMessage(token, message.Channel, BOT_NAME, "👆")
        } else {
            slack.PostMessage(token, message.Channel, BOT_NAME, "http://go-imghr.ds-12.com/"+imgPath)
        }
    case "jma":
        targetDate := time.Now().UTC().Add(time.Duration(-5)*time.Minute).Truncate(5 * time.Minute).Format("200601021504")
        imgPath := this.JmaImageGenerator.Generate(targetDate)
        if imgPath == "" {
            time.Sleep(1 * time.Second)
            slack.PostMessage(token, message.Channel, BOT_NAME, "👆")
        } else {
            slack.PostMessage(token, message.Channel, BOT_NAME, "http://go-imghr.ds-12.com/"+imgPath)
        }
    }
}

func (this *MessageEventHandler) Handle(event Event) {
    var message Message
    json.Unmarshal(event.Raw, &message)
    if IsBotComandAlias(message.Text) == true {
        switch message.Text {
        case "n":
            this.ExecuteCommand(message, "img", "長澤まさみ")
        }
        return
    }

    if IsBotCommand(message.Text) == false {
        return
    }

    command, argv := ParseCommand(message.Text)
    this.ExecuteCommand(message, command, argv)
}

func main() {
    token := os.Getenv("SLACK_TOKEN")
    if token == "" {
        log.Fatal("not set SLACK_TOKEN")
    }
    log.Print("init")
    eventHandler := NewEventHandler()
    eventHandler.SetExceptionHandler(func (event Event) {
        log.Print("Unknown Event: ", event.Type)
    })
    messageEventHandler := NewMessageEventHandler()
    eventHandler.AddHandler("message", messageEventHandler.Handle)
again:
    ws := slack.ConnectSocket(token)
    eventChan, resultChan := slack.StartReading(ws)
    startTime := int(time.Now().Unix())
    for {
        select {
        case buf := <-eventChan:
            var event Event
            json.Unmarshal(buf, &event)
            if event.Ts == "" {
                continue
            }
            ts := strings.Split(event.Ts, ".")[0]
            i, _ := strconv.Atoi(ts)
            if i < startTime {
                log.Print("skip event")
                continue
            }
            event.Raw = buf
            eventHandler.Handle(event)
        case <-resultChan:
            log.Print("disconnect server")
            goto again
        }
    }
}
