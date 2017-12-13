package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"../flow"
)

type Env struct {
	Bot   map[string]string
	Store map[string]string
}

func loadEnv() (env *Env, err error) {
	var body []byte

	err = flow.Go(
		func() (err error) {
			body, err = ioutil.ReadFile("./.env.json")
			return
		},

		func() (err error) {
			return json.Unmarshal(body, &env)
		},
	)

	return
}

func main() {
	env, err := loadEnv()

	if err != nil {
		panic(err)
	}

	store, err := InitStore(env.Store)

	if err != nil {
		panic(err)
	}

	bot := CreateBot(env.Bot["uin"], env.Bot["nick"], env.Bot["aimsid"])
	events, err := bot.FetchEvents()

	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	for _, rec := range toTimelineRecords(events.Events) {
		fmt.Printf("Insert: %#v\n", rec)

		res, err := store.Records.Insert(rec)
		fmt.Printf("Result: %#v, %#v\n", err, res)
	}
}

func toTimelineRecords(events []Event) []TimelineRecord {
	records := make([]TimelineRecord, 0, len(events))

	for _, rawEvt := range events {
		if rawEvt.Type == "im" {
			evtData := rawEvt.GetIMData()
			records = append(records, TimelineRecord{
				Id: evtData.MsgID,
				TS: evtData.Timestamp,
				Source: TimelineRecordSource{
					Id:   evtData.Source.AimID,
					Name: evtData.Source.Friendly,
				},
				Author: TimelineRecordAuthor{
					Login: evtData.MChatAttrs.Sender,
					Name:  evtData.MChatAttrs.SenderName,
				},
				Body: evtData.Message,
			})
		}
	}

	return records
}
