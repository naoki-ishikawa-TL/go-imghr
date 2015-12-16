package main

import (
	"./amesh"
	"./env"
	"./giphy"
	"./google"
	"./jma"
	"./wikipedia"
	"github.com/f110/go-ihr"
	"github.com/f110/go-ihr-console"
	"github.com/f110/go-ihr-slack"
	"log"
	"math/rand"
	"os"
	"time"
)

const (
	BotName     = "imghr"
	IhrID       = "U037GMSJ9"
	ProfileIcon = "http://go-imghr.ds-12.com/data/ihr_icon.png"
)

var (
	Version string
)

var _ = slack.NewSlackConnector
var _ = console.NewConsoleConnector

func imageSearch(word string) string {
	url := google.ImageSearch(word)
	if url == "" {
		return "( ˘ω˘ )ｽﾔｧ"
	} else {
		return url + "#.png"
	}
}

func main() {
	token := os.Getenv("SLACK_TOKEN")
	if env.DEBUG == false {
		if token == "" {
			log.Fatal("not set SLACK_TOKEN")
		}
	}
	log.Print("init")
	ignoreUsers := []string{}

	ameshImageGenerator := amesh.NewAmeshImageGenerator()
	jmaImageGenerator := jma.NewJmaImageGenerator()
	lastNeteroTime := time.Now().Add(time.Duration(-13) * time.Hour)

	var connector ihr.Connector
	if env.DEBUG == true {
		connector = console.NewConsoleConnector()
	} else {
		connector = slack.NewSlackConnector(token, ProfileIcon)
	}

	robot := ihr.NewBot(connector, BotName, ignoreUsers)
	robot.Command(
		"help",
		"これ",
		func(msg ihr.Event) {
			msg.Say(robot.ShowHelp())
		},
	)
	robot.Command(
		"version",
		"show version",
		func(msg ihr.Event) {
			msg.Say("version: " + Version)
		},
	)
	robot.Command(
		"ping",
		"ピンポン",
		func(msg ihr.Event) {
			msg.Say("pong")
		},
	)
	robot.Command(
		"amesh",
		"アメッシュの画像を生成する",
		func(msg ihr.Event) {
			targetDate := time.Now().Add(time.Duration(-1) * time.Minute).Truncate(5 * time.Minute).Format("200601021504")
			imgPath := ameshImageGenerator.Generate(targetDate)

			msg.Say("http://go-imghr.ds-12.com/" + imgPath)
		},
	)
	robot.Command(
		"jma",
		"高解像度ナウキャストの画像を生成する",
		func(msg ihr.Event) {
			targetDate := time.Now().UTC().Add(time.Duration(-5) * time.Minute).Truncate(5 * time.Minute).Format("200601021504")
			imgPath := jmaImageGenerator.Generate(targetDate)

			msg.Say("http://go-imghr.ds-12.com/" + imgPath)
		},
	)
	robot.Command(
		"gif",
		"ランダムにアニメーションgifを表示する",
		func(msg ihr.Event) {
			msg.Say(giphy.Random())
		},
	)
	robot.CommandWithArgv(
		"img",
		"画像検索",
		func(msg ihr.Event) {
			msg.Say(imageSearch(msg.Argv))
		},
	)
	robot.CommandWithArgv(
		"ihr",
		"ihrの画像検索",
		func(msg ihr.Event) {
			words := []string{"aww cat", "柴犬", "豆柴", "aww dog", "aww shiba"}
			msg.Say(imageSearch(words[rand.Intn(len(words))]))
		},
	)
	robot.CommandWithArgv(
		"wikipedia",
		"wikipedia検索",
		func(msg ihr.Event) {
			summary, err := wikipedia.GetSummary(msg.Argv)
			if err != nil {
				msg.Say("ページがないよ")
			} else {
				msg.Say(summary + "\n" + wikipedia.GenerateJaWikipediaURL(msg.Argv))
			}
		},
	)
	robot.CommandWithArgv(
		"gif",
		"アニメーションGif検索",
		func(msg ihr.Event) {
			msg.Say(giphy.Search(msg.Argv))
		},
	)

	robot.Hear("(?i)(hanakin|花金|金曜|ファナキン|tgif)", func(msg ihr.Event) {
		if time.Now().Weekday() != time.Friday {
			return
		}

		msg.Say("花金だーワッショーイ！テンションAGEAGEマック")
	})
	robot.Hear("\\An\\z", func(msg ihr.Event) {
		msg.Say(imageSearch("長澤まさみ"))
	})
	robot.Hear("\\Ak\\z", func(msg ihr.Event) {
		msg.Say(imageSearch("木村文乃"))
	})
	robot.Hear("\\Aa\\z", func(msg ihr.Event) {
		msg.Say(imageSearch("有村架純"))
	})
	robot.Hear("\\Af\\z", func(msg ihr.Event) {
		msg.Say(imageSearch("Ferrari 458 Italia"))
	})
	robot.Hear("\\Ap\\z", func(msg ihr.Event) {
		msg.Say(imageSearch("Porsche 991 GT3 RS"))
	})
	robot.Appearance(IhrID, func(msg ihr.Event) {
		if time.Now().Sub(lastNeteroTime) < 12*time.Hour {
			return
		}
		if 14 > time.Now().Hour() || time.Now().Hour() > 17 {
			return
		}
		lastNeteroTime = time.Now()
		msg.Say("I H R は 寝 て ろ ！ ！")
	})

	robot.Start()
}
