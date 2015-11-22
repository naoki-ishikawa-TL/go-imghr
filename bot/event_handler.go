package bot

import (
	"regexp"
)

type EventHandler struct {
	ignoreUsers map[string]bool
	commands map[string][]Command
}

type Command struct {
	eventType string
	description string
	pattern  *regexp.Regexp
	user string
	argv     bool
	callback func(Event)
}

func NewEventHandler(ignoreUsers []string) *EventHandler {
	ignore := make(map[string]bool)
	for _, v := range(ignoreUsers) {
		ignore[v] = true
	}
	
	return &EventHandler{
		ignoreUsers: ignore,
		commands: make(map[string][]Command, 0),
	}
}

func (this *EventHandler) AddHandler(eventType string, command *Command) {
	if this.commands[eventType] == nil {
		this.commands[eventType] = make([]Command, 0)
	}
	
	this.commands[eventType] = append(this.commands[eventType], *command)
}

func (this *EventHandler) AddCommand(pattern *regexp.Regexp, description string, callback func(Event), argv bool) {
	command := &Command{pattern: pattern, description: description, callback: callback, argv: argv}
	this.AddHandler(MessageEvent, command)
}

func (this *EventHandler) Appearance(user string, callback func(Event)) {
	command := &Command{user: user, callback: callback}
	this.AddHandler(UserTypingEvent, command)
}

func (this *EventHandler) Handle(event Event) {
	if _, ok := this.ignoreUsers[event.User]; ok == true {
		return
	}
	for _, command := range this.commands[event.Type] {
		switch event.Type {
		case MessageEvent:
			if command.pattern.MatchString(event.Message) == true {
				if command.argv == true {
					matched := command.pattern.FindStringSubmatch(event.Message)
					event.Argv = matched[1]
					command.callback(event)
				} else {
					command.callback(event)
				}
				return
			}
		case UserTypingEvent:
			if event.User == command.user {
				command.callback(event)
			}
		}
	}
}

func (this *EventHandler) RenderHelp() string {
	help := ""
	for _, v := range this.commands {
		for _, c := range v {
			if c.description != "" {
				help += c.description + "\n"
			}
		}
	}
	
	return help
}