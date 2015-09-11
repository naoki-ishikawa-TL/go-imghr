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

func main() {
    token := os.Getenv("SLACK_TOKEN")
    if token == "" {
        log.Fatal("not set SLACK_TOKEN")
    }
    log.Print("init")
    eventHandler := slack.NewEventHandler()
    eventHandler.SetExceptionHandler(func (event slack.Event) {
        log.Print("Unknown Event: ", event.Type)
    })
    messageEventHandler := slack.NewMessageEventHandler()
    userTypingEventHandler := slack.NewUserTypingEventHandler()
    eventHandler.AddHandler("message", messageEventHandler.Handle)
    eventHandler.AddHandler("user_typing", userTypingEventHandler.Handle)
again:
    ws := slack.ConnectSocket(token)
    if ws == nil {
        log.Fatal("Could not connect websocket")
    }
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
