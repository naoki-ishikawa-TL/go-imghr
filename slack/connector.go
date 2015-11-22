package slack

import (
	"../bot"
	"encoding/json"
	"errors"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type SlackConnector struct {
	eventChan chan *bot.Event

	token      string
	breakChan  chan error
	bufChan    chan []byte
	startTime  int
	connection *websocket.Conn
}

type RtmStart struct {
	Url string
}

type PostMessage struct {
	Ok bool
}

type Event struct {
	Type string
	Ts   string
	Raw  []byte
}

type Message struct {
	Type    string
	SubType string `json:"subtype"`
	Ts      string
	Channel string
	User    string
	Text    string
	ts      string
}

type UserTyping struct {
	Type    string
	Channel string
	User    string
}

func NewSlackConnector(token string) *SlackConnector {
	startTime := int(time.Now().Unix())

	return &SlackConnector{
		token:     token,
		startTime: startTime,
		eventChan: make(chan *bot.Event),
	}
}

func (this *SlackConnector) Connect() {
	v := url.Values{}
	v.Set("token", this.token)
	response, _ := http.PostForm("https://slack.com/api/rtm.start", v)
	dec := json.NewDecoder(response.Body)
	var data RtmStart
	dec.Decode(&data)
	log.Print("start connect to ", data.Url)
	ws, err := websocket.Dial(data.Url, "", "http://localhost")
	if err != nil {
		log.Fatal(err)
	}

	this.connection = ws
	this.startReading()
}

func (this *SlackConnector) Listen() error {
	for {
		select {
		case buf := <-this.bufChan:
			var event Event
			json.Unmarshal(buf, &event)

			if event.Ts != "" {
				ts := strings.Split(event.Ts, ".")[0]
				i, _ := strconv.Atoi(ts)
				if i < this.startTime {
					log.Print("skip event")
					continue
				}
			}

			var botEvent *bot.Event
			var eventType string
			var message string
			var channel string
			var user string
			var mentionName string
			var timestamp time.Time
			switch event.Type {
			case "message":
				var messageEvent Message
				json.Unmarshal(buf, &messageEvent)
				if messageEvent.User == "" {
					continue
				}
				
				eventType = bot.MessageEvent
				message = messageEvent.Text
				channel = messageEvent.Channel
				user = messageEvent.User
				mentionName = "<@" + messageEvent.User + ">"
			case "user_typing":
				var userTypingEvent UserTyping
				json.Unmarshal(buf, &userTypingEvent)
				if userTypingEvent.User == "" {
					continue
				}

				eventType = bot.UserTypingEvent
				channel = userTypingEvent.Channel
				user = userTypingEvent.User
			default:
				eventType = bot.UnknownEvent
			}

			botEvent = &bot.Event{
				Type:        eventType,
				Message:     message,
				Channel:     channel,
				User:        user,
				MentionName: mentionName,
				Timestamp:   timestamp,
			}

			this.eventChan <- botEvent
		case <-this.breakChan:
			log.Print("disconnect server")
			return errors.New("disconnect")
		}
	}
}

func (this *SlackConnector) ReceivedEvent() chan *bot.Event {
	return this.eventChan
}

func (this *SlackConnector) Send(event *bot.Event, username string, text string) error {
	v := url.Values{}
	v.Set("token", this.token)
	v.Set("channel", event.Channel)
	v.Set("text", text)
	v.Set("username", username)

	response, _ := http.PostForm("https://slack.com/api/chat.postMessage", v)
	dec := json.NewDecoder(response.Body)
	var data PostMessage
	dec.Decode(&data)

	if data.Ok != true {
		errors.New("post error")
	}

	return nil
}

func (this *SlackConnector) startReading() {
	log.Print("start reading")
	var msg []byte
	this.bufChan = make(chan []byte)
	this.breakChan = make(chan error)

	go func() {
		for {
			var tmp = make([]byte, 512)
			n, err := this.connection.Read(tmp)
			if err == io.EOF {
				log.Print("EOF")
				this.breakChan <- errors.New("eof")
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
				this.bufChan <- msg
				msg = nil
			}
		}
	}()
}