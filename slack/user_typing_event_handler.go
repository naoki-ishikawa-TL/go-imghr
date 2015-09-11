package slack

import (
    "time"
    "encoding/json"
)

const IHR_ID = "U037GMSJ9"

type UserTyping struct {
    Event
    Channel string
    User string
}

type UserTypingEventHandler struct {
    PostFlag bool
    Time time.Time
    FireChan chan Event
}

func NewUserTypingEventHandler() *UserTypingEventHandler {
    defaultTime := time.Now().Add(time.Duration(-13)*time.Hour)
    fireChan := make(chan Event)

    this := &UserTypingEventHandler{PostFlag: false, Time: defaultTime, FireChan: fireChan}

    go func() {
        for {
            select {
            case event := <-this.FireChan:
                if this.IsEnable() != true {
                    continue
                }
                var userTyping UserTyping
                json.Unmarshal(event.Raw, &userTyping)

				if userTyping.User != IHR_ID {
					continue
				}

                this.PostFlag = true
                this.Time = time.Now()
                PostMessage(userTyping.Channel, BOT_NAME, "I H R は 寝 て ろ ！ ！")
            }
        }
    }()

    return this
}

func (this *UserTypingEventHandler) IsEnable() bool {
    if time.Now().Sub(this.Time) > 12 * time.Hour {
        this.PostFlag = false
    }
    if this.PostFlag == true {
        return false
    }

    if 14 <= time.Now().Hour() && time.Now().Hour() < 17 {
        return true
    } else {
        return false
    }
}

func (this *UserTypingEventHandler) Handle(event Event) {
    this.FireChan <- event
}
