package main

import (
	"./api"
	botAPI "./bot"
	"./env"
	"./events"
	"./server"
	"./space"
	"fmt"
	"log"
	"time"
)

func main() {
	botEnv := env.Get("bot")
	bot := botAPI.CreateBot(botEnv["uin"], botEnv["nick"], botEnv["token"])

	go func() {
		record, err := space.Records.GetLast()
		lastSeqNum := uint(0)

		if err != nil {
			log.Fatal(err)
		} else if record != nil {
			lastSeqNum = record.SeqNum
		}

		for {
			fmt.Println("Last seqNum:", lastSeqNum)
			result, err := bot.FetchEvents(lastSeqNum)

			if err == nil {
				lastSeqNum = events.Processing(lastSeqNum, result.Events)
				//break
			} else {
				log.Fatal(err)
			}

			time.Sleep(10 * time.Second)
		}
	}()

	err := (&server.HttpServer{env.Get("http")["host"]}).Start(
		server.Routes{
			"/":              api.WwwDashboard,
			"/ping/":         api.Ping,
			"/project/get/":  api.ProjectGet,
			"/project/reg/":  api.ProjectReg,
			"/token/create/": api.TokenCreate,
			"/record/list/":  api.RecordList,
			"/buddy/list/":   api.BuddyList,
		},
	)

	if err != nil {
		panic(err)
	}
}
