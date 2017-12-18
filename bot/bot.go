package bot

import (
	"fmt"
	"strings"
	"time"
)

type Bot struct {
	UIN   string
	NICK  string
	token string
}

func (bot *Bot) Call(method string, query HttpQuery) (*RawAPIResponse, error) {
	url := method
	result := &RawAPIResponse{}

	query["aimsid"] = bot.token
	query["timeout"] = "1000"

	if strings.Index(method, "https://") != 0 {
		url = "https://botapi.icq.net/" + url
	}

	err := HttpGet(url, query).AsJSON(result)

	return result, err
}

func (bot *Bot) GetBuddyList() (*BuddyList, error) {
	resp, err := bot.Call("getBuddyList", HttpQuery{})

	if err != nil {
		return nil, err
	}

	return resp.AsBuddyList()
}

func (bot *Bot) FetchEvents(seqNum uint) (*FetchEventsData, error) {
	query := HttpQuery{"seqNum": fmt.Sprint(seqNum)}

	fetch := func(url string, delay int) (data *FetchEventsData, err error) {
		if delay > 0 {
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}

		resp, err := bot.Call(url, query)

		if err == nil {
			data, err = resp.AsFetchEvents()
		}

		return
	}

	data, err := fetch("fetchEvents", 0)

	if err == nil {
		data, err = fetch(data.FetchBaseURL, data.TimeToNextFetch)
	}

	return data, err
}

func CreateBot(uin, nick, token string) Bot {
	return Bot{
		UIN:   uin,
		NICK:  nick,
		token: token,
	}
}
