package main

import (
	"./amesh"
	"./bot"
	"./google"
	"./jma"
	"./slack"
	"./wikipedia"
	"log"
	"math/rand"
	"os"
	"time"
)

const (
	BotName = "imghr"
	IhrID = "U037GMSJ9"
)

var (
	Version string
)

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
	if token == "" {
		log.Fatal("not set SLACK_TOKEN")
	}
	log.Print("init")
	ignoreUsers := []string{"U037G7YJF"}

	ameshImageGenerator := amesh.NewAmeshImageGenerator()
	jmaImageGenerator := jma.NewJmaImageGenerator()
	lastNeteroTime := time.Now().Add(time.Duration(-13) * time.Hour)

	slackConnector := slack.NewSlackConnector(token)
	robot := bot.NewBot(slackConnector, BotName, ignoreUsers)
	robot.Command(
		"help",
		"これ",
		func(msg bot.Event) {
			msg.Say(robot.ShowHelp())
		},
	)
	robot.Command(
		"version",
		"show version",
		func(msg bot.Event) {
			msg.Say("version: "+Version)
		},
	)
	robot.Command(
		"ping",
		"ピンポン",
		func(msg bot.Event) {
			msg.Say("pong")
		},
	)
	robot.Command(
		"amesh",
		"アメッシュの画像を生成する",
		func(msg bot.Event) {
			targetDate := time.Now().Add(time.Duration(-1) * time.Minute).Truncate(5 * time.Minute).Format("200601021504")
			imgPath := ameshImageGenerator.Generate(targetDate)

			msg.Say("http://go-imghr.ds-12.com/" + imgPath)
		},
	)
	robot.Command(
		"jma",
		"高解像度ナウキャストの画像を生成する",
		func(msg bot.Event) {
			targetDate := time.Now().UTC().Add(time.Duration(-5) * time.Minute).Truncate(5 * time.Minute).Format("200601021504")
			imgPath := jmaImageGenerator.Generate(targetDate)

				msg.Say("http://go-imghr.ds-12.com/" + imgPath)
		},
	)
	robot.CommandWithArgv(
		"img",
		"画像検索",
		func(msg bot.Event) {
			msg.Say(imageSearch(msg.Argv))
		},
	)
	robot.CommandWithArgv(
		"ihr",
		"ihrの画像検索",
		func(msg bot.Event) {
			words := []string{"aww cat", "柴犬", "豆柴", "aww dog", "aww shiba"}
			msg.Say(imageSearch(words[rand.Intn(len(words))]))
		},
	)
	robot.CommandWithArgv(
		"wikipedia",
		"wikipedia検索",
		func(msg bot.Event) {
			summary, err := wikipedia.GetSummary(msg.Argv)
			if err != nil {
				msg.Say("ページがないよ")
			} else {
				msg.Say(summary + "\n" + wikipedia.GenerateJaWikipediaURL(msg.Argv))
			}
		},
	)

	robot.Hear("(?i)(hanakin|花金|金曜|ファナキン|tgif)", func(msg bot.Event) {
		if time.Now().Weekday() != time.Friday {
			return
		}

		msg.Say("花金だーワッショーイ！テンションAGEAGEマック")
	})
	robot.Hear("\\An\\z", func(msg bot.Event) {
		msg.Say(imageSearch("長澤まさみ"))
	})
	robot.Hear("\\Ak\\z", func(msg bot.Event) {
		msg.Say(imageSearch("木村文乃"))
	})
	robot.Hear("\\Aa\\z", func(msg bot.Event) {
		msg.Say(imageSearch("有村架純"))
	})
	robot.Hear("\\Af\\z", func(msg bot.Event) {
		msg.Say(imageSearch("Ferrari 458 Italia"))
	})
	robot.Hear("\\Ap\\z", func(msg bot.Event) {
		msg.Say(imageSearch("Porsche 991 GT3 RS"))
	})
	robot.Appearance(IhrID, func(msg bot.Event) {
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
