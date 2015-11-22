package bot

import (
	"log"
	"regexp"
)

var (
	botInstance *Bot
)

type Bot struct {
	connector        Connector
	name             string
	connectErrorChan chan bool
	eventHandler     *EventHandler
}

func BotInstance() *Bot {
	return botInstance
}

func NewBot(connector Connector, name string, ignoreUsers []string) *Bot {
	if botInstance != nil {
		return botInstance
	}

	bot := &Bot{
		connector:        connector,
		name:             name,
		connectErrorChan: make(chan bool),
		eventHandler:     NewEventHandler(ignoreUsers),
	}
	botInstance = bot

	bot.Connect()
	return bot
}

func (self *Bot) Connect() {
	self.connector.Connect()

	go func() {
		res := self.connector.Listen()
		if res != nil {
			self.connectErrorChan <- true
		}
	}()
}

func (self *Bot) Start() {
	for {
		select {
		case event := <-self.connector.ReceivedEvent():
			event.Bot = self
			go self.eventHandler.Handle(*event)
		case <-self.connectErrorChan:
			log.Print("reconnect")
			self.Connect()
		}
	}
}

func (self *Bot) Send(event *Event, text string) {
	self.connector.Send(event, self.name, text)
}

func (self *Bot) Hear(pattern string, callback func(Event)) {
	self.eventHandler.AddCommand(regexp.MustCompile(pattern), "", callback, false)
}

func (self *Bot) Command(pattern string, description string, callback func(Event)) {
	self.eventHandler.AddCommand(regexp.MustCompile("\\A"+self.name+"\\s+"+pattern), pattern+" - "+description, callback, false)
}

func (self *Bot) CommandWithArgv(pattern string, description string, callback func(Event)) {
	self.eventHandler.AddCommand(regexp.MustCompile("\\A"+self.name+"\\s+"+pattern+"(?:\\s+(.+))*"), pattern+" - "+description, callback, true)
}

func (self *Bot) Appearance(user string, callback func(Event)) {
	self.eventHandler.Appearance(user, callback)
}

func (self *Bot) ShowHelp() string {
	return self.eventHandler.RenderHelp()
}

func (self *Event) Say(text string) {
	self.Bot.Send(self, text)
}
