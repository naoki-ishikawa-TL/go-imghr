package slack

import (
	"../amesh"
	"../google"
	"../jma"
	"encoding/json"
	"regexp"
	"time"
)

type Message struct {
	Event
	Channel string
	User    string
	Text    string
	ts      string
}

type MessageEventHandler struct {
	AmeshImageGenerator *amesh.AmeshImageGenerator
	JmaImageGenerator   *jma.JmaImageGenerator
}

func isBotCommandAlias(message string) bool {
	matched, _ := regexp.MatchString("^(n|f)$", message)

	return matched
}

func isBotCommand(message string) bool {
	matched, _ := regexp.MatchString("^imghr\\s+", message)

	return matched
}

func parseCommand(message string) (string, string) {
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
	if isBotCommandAlias(message.Text) == true {
		switch message.Text {
		case "n":
			this.ExecuteCommand(message, "img", "長澤まさみ")
		case "f":
			this.ExecuteCommand(message, "img", "Ferrari")
		}
		return
	}

	if isBotCommand(message.Text) == false {
		return
	}

	command, argv := parseCommand(message.Text)
	this.ExecuteCommand(message, command, argv)
}

func (this *MessageEventHandler) ExecuteCommand(message Message, command string, argv string) {
	switch command {
	case "ping":
		PostMessage(message.Channel, BOT_NAME, "pong")
	case "img":
		url := google.ImageSearch(argv)
		if url == "" {
			PostMessage(message.Channel, BOT_NAME, "( ˘ω˘ )ｽﾔｧ")
		} else {
			PostMessage(message.Channel, BOT_NAME, url+"#.png")
		}
	case "amesh":
		targetDate := time.Now().Add(time.Duration(-1) * time.Minute).Truncate(5 * time.Minute).Format("200601021504")
		imgPath := this.AmeshImageGenerator.Generate(targetDate)
		PostMessage(message.Channel, BOT_NAME, "http://go-imghr.ds-12.com/"+imgPath)
	case "jma":
		targetDate := time.Now().UTC().Add(time.Duration(-5) * time.Minute).Truncate(5 * time.Minute).Format("200601021504")
		imgPath := this.JmaImageGenerator.Generate(targetDate)
		PostMessage(message.Channel, BOT_NAME, "http://go-imghr.ds-12.com/"+imgPath)
	}
}
