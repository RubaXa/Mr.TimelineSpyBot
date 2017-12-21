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
	"sync"
	"time"
)

func main() {
	var srv *server.HttpServer

	botEnv := env.Get("bot")
	bot := botAPI.CreateBot(botEnv["uin"], botEnv["nick"], botEnv["token"])

	lock := &sync.WaitGroup{}
	lock.Add(1)

	srv = server.CreateHttpServer(env.Get("http")["host"], server.Routes{
		"/":              api.WwwDashboard,
		"/ping/":         api.Ping,
		"/project/get/":  api.ProjectGet,
		"/project/reg/":  api.ProjectReg,
		"/token/create/": api.TokenCreate,
		"/record/list/":  api.RecordList,
		"/buddy/list/":   api.BuddyList,
	})

	go func() {
		srv.Start()
		lock.Done()
	}()

	go func() {
		var records []space.RecordsEntry

		record, err := space.Records.GetLast()
		lastSeqNum := uint(0)

		if err != nil {
			log.Fatal(err)
		} else if record != nil {
			lastSeqNum = record.SeqNum
		}

		for {
			fmt.Println("Last seqNum:", lastSeqNum)

			srv.PushEvent(0, "next", map[string]uint{
				"seq_num": lastSeqNum,
			})

			result, err := bot.FetchEvents(lastSeqNum)

			if err == nil {
				lastSeqNum, records = events.Processing(lastSeqNum, result.Events)

				for _, r := range records {
					srv.PushEvent(record.ProjectId, "record", r)
				}

				//break
			} else {
				log.Fatal(err)
			}

			time.Sleep(10 * time.Second)
		}
	}()

	lock.Wait()
}
