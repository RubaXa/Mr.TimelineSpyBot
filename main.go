package main

import (
	"./api"
	"./env"
	"./server"
	"./space"
	"fmt"
)

func main() {
	go func() {
		id := space.Records.GetLast()
		fmt.Println(id)
	}()

	err := (&server.HttpServer{env.Get("http")["host"]}).Start(
		server.Routes{
			"/ping/":        api.Ping,
			"/project/get/": api.Get,
			"/project/reg/": api.Reg,
		},
	)

	if err != nil {
		panic(err)
	}

	//projects := &space.Projects{box.GetSpace(Env.Space["projects"])}
	//project := projects.Get(24)

	//project.AddChat
	//spaceProjects.SelectOne(24, project)

	//fmt.Printf("project: %#v\n", project)

	//if !project.HasId() {
	//	project.Name = "MailRu::Timeline"
	//	err = spaceProjects.Insert(project)
	//	fmt.Println("Res:", err)
	//}

	//records := box.GetSpace(env.Space["records"])
	//
	//bot := CreateBot(env.Bot["uin"], env.Bot["nick"], env.Bot["aimsid"])
	//events, err := bot.FetchEvents()
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//for _, rec := range toTimelineRecords(events.Events) {
	//	err := records.Insert(&rec)
	//	fmt.Printf("Insert: %#v\n", rec)
	//	fmt.Printf("Error: %#v\n", err)
	//}
}

func toTimelineRecords(events []Event) []space.RecordsEntry {
	records := make([]space.RecordsEntry, 0, len(events))

	for _, rawEvt := range events {
		if rawEvt.Type == "im" {
			evtData := rawEvt.GetIMData()
			records = append(records, space.RecordsEntry{
				MsgId: evtData.MsgID,
				TS:    evtData.Timestamp,
				Source: space.RecordsEntrySource{
					Id:   evtData.Source.AimID,
					Name: evtData.Source.Friendly,
				},
				Author: space.RecordsEntryAuthor{
					Login: evtData.MChatAttrs.Sender,
					Name:  evtData.MChatAttrs.SenderName,
				},
				Body: evtData.Message,
			})
		}
	}

	return records
}
