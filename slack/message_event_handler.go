package slack

import (
    "../amesh"
    "../jma"
    "../google"
    "time"
    "os"
    "encoding/json"
    "regexp"
)

const (
    BOT_NAME = "imghr"
)

type Message struct {
    Event
    Channel string
    User string
    Text string
    ts string
}

type MessageEventHandler struct {
    AmeshImageGenerator *amesh.AmeshImageGenerator
    JmaImageGenerator *jma.JmaImageGenerator
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

func NewMessageEventHandler() *MessageEventHandler {
    ameshImageGenerator := amesh.NewAmeshImageGenerator()
    jmaImageGenerator := jma.NewJmaImageGenerator()

    return &MessageEventHandler{AmeshImageGenerator: ameshImageGenerator, JmaImageGenerator: jmaImageGenerator}
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

func (this *MessageEventHandler) ExecuteCommand(message Message, command string, argv string) {
    token := os.Getenv("SLACK_TOKEN")
    switch command {
    case "ping":
        PostMessage(token, message.Channel, BOT_NAME, "pong")
    case "img":
        url := google.ImageSearch(argv)
        if url == "" {
            PostMessage(token, message.Channel, BOT_NAME, "( ˘ω˘ )ｽﾔｧ")
        } else {
            PostMessage(token, message.Channel, BOT_NAME, url+"#.png")
        }
    case "amesh":
        targetDate := time.Now().Add(time.Duration(-1)*time.Minute).Truncate(5 * time.Minute).Format("200601021504")
        imgPath := this.AmeshImageGenerator.Generate(targetDate)
        if imgPath == "" {
            time.Sleep(1 * time.Second)
            PostMessage(token, message.Channel, BOT_NAME, "👆")
        } else {
            PostMessage(token, message.Channel, BOT_NAME, "http://go-imghr.ds-12.com/"+imgPath)
        }
    case "jma":
        targetDate := time.Now().UTC().Add(time.Duration(-5)*time.Minute).Truncate(5 * time.Minute).Format("200601021504")
        imgPath := this.JmaImageGenerator.Generate(targetDate)
        if imgPath == "" {
            time.Sleep(1 * time.Second)
            PostMessage(token, message.Channel, BOT_NAME, "👆")
        } else {
            PostMessage(token, message.Channel, BOT_NAME, "http://go-imghr.ds-12.com/"+imgPath)
        }
    }
}
