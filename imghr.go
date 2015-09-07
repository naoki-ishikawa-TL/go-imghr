package main

import (
    "golang.org/x/net/websocket"
    "log"
    "net/http"
    "net/url"
    "encoding/json"
    "os"
    "io"
    "math/rand"
    "time"
    "regexp"
    "strings"
    "strconv"
    "./amesh"
    "./jma"
    "./image"
)

type RtmStart struct {
    Url string
}

type AbstractRestful struct {
    Ok bool
}

const BOT_NAME = "goihr"

type ImageSearchApi struct {
    ResponseData struct {
        Results []struct {
            UnescapedUrl string
        }
    }
}

func GoogleImageSearch(query string) string {
    rand.Seed(time.Now().UnixNano())
    v := url.Values{}
    v.Set("v", "1.0")
    v.Set("rsz", "8")
    v.Set("q", query)
    v.Set("safe", "active")
    response, _ := http.Get("http://ajax.googleapis.com/ajax/services/search/images?"+v.Encode())
    dec := json.NewDecoder(response.Body)
    var data ImageSearchApi
    dec.Decode(&data)
    if len(data.ResponseData.Results) == 0 {
        return ""
    }
    i := rand.Intn(len(data.ResponseData.Results))

    return data.ResponseData.Results[i].UnescapedUrl
}

func PostMessage(token string, channel string, username string, text string) <-chan bool {
    resultChan := make(chan bool)
    v := url.Values{}
    v.Set("token", token)
    v.Set("channel", channel)
    v.Set("text", text)
    v.Set("username", username)
    go func() {
        response, _ := http.PostForm("https://slack.com/api/chat.postMessage", v)
        dec := json.NewDecoder(response.Body)
        var data AbstractRestful
        dec.Decode(&data)
        resultChan <- data.Ok
    }()

    return resultChan
}

func connectSocket(token string) *websocket.Conn {
    v := url.Values{}
    v.Set("token", token)
    response, _ := http.PostForm("https://slack.com/api/rtm.start", v)
    dec := json.NewDecoder(response.Body)
    var data RtmStart
    dec.Decode(&data)
    log.Print("start connect to ", data.Url)
    ws, err := websocket.Dial(data.Url, "", "http://localhost")
    if err != nil {
        log.Fatal(err)
    }

    return ws
}

func startReading(conn *websocket.Conn) (<-chan []byte, <-chan bool) {
    log.Print("reading...")
    var msg []byte
    sendChan := make(chan []byte)
    breakChan := make(chan bool)

    go func() {
        for {
            var tmp = make([]byte, 512)
            n, err := conn.Read(tmp)
            if err == io.EOF {
                breakChan <- true
                break
            }
            if err != nil {
                log.Fatal(err)
            }
            if msg != nil {
                msg = append(msg, tmp[:n]...)
            } else {
                msg = tmp[:n]
            }
            if n != 512 {
                sendChan <- msg
                msg = nil
            }
        }
    }()

    return sendChan, breakChan
}

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
    ImageGenerator *image.ImageGenerator
}

func NewMessageEventHandler() *MessageEventHandler {
    imageGenerator := image.NewImageGenerator()
    return &MessageEventHandler{ImageGenerator: imageGenerator}
}

func (this *MessageEventHandler) ExecuteCommand(message Message, command string, argv string) {
    token := os.Getenv("SLACK_TOKEN")
    switch command {
    case "ping":
        PostMessage(token, message.Channel, BOT_NAME, "pong")
    case "img":
        url := GoogleImageSearch(argv)
        if url == "" {
            PostMessage(token, message.Channel, BOT_NAME, "( ˘ω˘ )ｽﾔｧ")
        } else {
            PostMessage(token, message.Channel, BOT_NAME, url+"#.png")
        }
    case "amesh":
        targetDate := time.Now().Add(time.Duration(-1)*time.Minute).Truncate(5 * time.Minute).Format("200601021504")
        imgPath := this.ImageGenerator.Generate(command, targetDate)
        if imgPath == "" {
            time.Sleep(1 * time.Second)
            PostMessage(token, message.Channel, BOT_NAME, "👆")
        } else {
            PostMessage(token, message.Channel, BOT_NAME, "http://go-imghr.ds-12.com/"+imgPath)
        }
    case "jma":
        targetDate := time.Now().UTC().Add(time.Duration(-5)*time.Minute).Truncate(5 * time.Minute).Format("200601021504")
        imgPath := this.ImageGenerator.Generate(command, targetDate)
        if imgPath == "" {
            time.Sleep(1 * time.Second)
            PostMessage(token, message.Channel, BOT_NAME, "👆")
        } else {
            PostMessage(token, message.Channel, BOT_NAME, "http://go-imghr.ds-12.com/"+imgPath)
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
    messageEventHandler.ImageGenerator.AddGenerator("amesh", amesh.GenerateAmeshImage)
    messageEventHandler.ImageGenerator.AddGenerator("jma", jma.GenerateJmaImage)
    eventHandler.AddHandler("message", messageEventHandler.Handle)
again:
    ws := connectSocket(token)
    eventChan, resultChan := startReading(ws)
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
