package slack

import (
	"../amesh"
	"../google"
	"../jma"
	"../tenki"
	"encoding/json"
	"regexp"
	"time"
)

var BAN_USERS map[string]bool = map[string]bool{
	"U037G7YJF": true, // yusuke.shirakawa
	"U02G5LRKZ": false, // for debug
}

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
    Version string
    LastBuild string
}

func isBotCommandAlias(message string) bool {
	matched, _ := regexp.MatchString("^(n|f)$", message)

	return matched
}

func isBotCommand(message string) bool {
	matched, _ := regexp.MatchString("^imghr\\s+", message)

	return matched
}

func isTenkiBotCommand(message string) bool {
	matched, _ := regexp.MatchString("^tenkihr\\s+", message)

	return matched
}

func parseCommand(message string) (string, string) {
	re := regexp.MustCompile("^(?:imghr|tenkihr)\\s+(\\w+)(?:\\s+(.+))*")
	matched := re.FindStringSubmatch(message)
	if len(matched) == 2 {
		return matched[1], ""
	}
	return matched[1], matched[2]
}

func NewMessageEventHandler(version string, lastBuild string) *MessageEventHandler {
	ameshImageGenerator := amesh.NewAmeshImageGenerator()
	jmaImageGenerator := jma.NewJmaImageGenerator()

	return &MessageEventHandler{AmeshImageGenerator: ameshImageGenerator, JmaImageGenerator: jmaImageGenerator, Version: version, LastBuild: lastBuild}
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

	if isBotCommand(message.Text) == true {
		command, argv := parseCommand(message.Text)
		this.ExecuteCommand(message, command, argv)

		return
	}

	if isTenkiBotCommand(message.Text) == true {
		command, argv := parseCommand(message.Text)
		this.ExecuteTenkiCommand(message, command, argv)

		return
	}
}

func (this *MessageEventHandler) ExecuteCommand(message Message, command string, argv string) {
	if this.checkBanUser(message.User) == true {
		this.postMessageToBanUser(message.Channel)
		return
	}

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
    case "version":
        PostMessage(message.Channel, BOT_NAME, "last build:"+this.LastBuild+"  version:"+this.Version)
	}
}

func (this *MessageEventHandler) ExecuteTenkiCommand(message Message, command string, argv string) {
	if this.checkBanUser(message.User) == true {
		this.postMessageToBanUser(message.Channel)
		return
	}

	switch command {
	case "temperature":
		PostMessage(message.Channel, BOT_NAME, "今の気温は "+tenki.GetTemperature()+" 度だよ")
	case "humidity":
		PostMessage(message.Channel, BOT_NAME, "今の湿度は "+tenki.GetHumidity()+" %だよ")
	}
}

func (this *MessageEventHandler) checkBanUser(userId string) bool {
	if value, ok := BAN_USERS[userId]; ok && value {
		return true
	}

	return false
}

func (this *MessageEventHandler) postMessageToBanUser(channel string) {
	PostMessage(channel, BOT_NAME, "( ˘ω˘ )ｽﾔｧ")
}
