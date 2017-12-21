package bot

import (
	"encoding/json"
	"fmt"
)

type MapOfAny map[string]interface{}

type ServiceURL struct {
	Name string
	Urls struct {
		Profile     string
		Icon        string
		Settings    string
		Association string
		Icon10      string
		Compose     string
		Signup      string
		Inbox       string
	}
}

type Event struct {
	Type      string
	EventData MapOfAny
	SeqNum    uint
}

type ServiceEventData struct {
	ServiceURLs []ServiceURL
}

type EventMChatAttrs struct {
	ChatName   string
	Sender     string
	SenderName string
}

func (attrs *EventMChatAttrs) Parse(data interface{}) {
	raw, ok := data.(map[string]interface{})
	if ok {
		attrs.ChatName = raw["chat_name"].(string)
		attrs.Sender = raw["sender"].(string)
		attrs.SenderName = raw["senderName"].(string)
	}
}

type EventRawMsg struct {
	Base64Msg string
}

func (rawMsg *EventRawMsg) Parse(data interface{}) {
	if data != nil {
		rawMsg.Base64Msg = data.(map[string]interface{})["base64Msg"].(string)
	}
}

type EventSource struct {
	AimID     string
	DisplayID string
	Friendly  string
	State     string
	UserType  string
}

func (source *EventSource) Parse(data interface{}) {
	raw := data.(map[string]interface{})
	source.AimID = raw["aimId"].(string)
	source.DisplayID = raw["displayId"].(string)
	source.Friendly = raw["friendly"].(string)
	source.State = raw["state"].(string)
	source.UserType = raw["userType"].(string)
}

type IMEventData struct {
	MChatAttrs   EventMChatAttrs
	Autoresponse int
	Imf          string
	Message      string
	MsgID        string
	Notification string
	RawMsg       EventRawMsg
	Source       EventSource
	Timestamp    uint
}

func (evt *Event) GetIMData() (res IMEventData) {
	if evt.Type == "im" {
		raw := evt.EventData

		res.RawMsg.Parse(raw["rawMsg"])
		res.MChatAttrs.Parse(raw["MChat_Attrs"])
		res.Source.Parse(raw["source"])

		res.Autoresponse = int(raw["autoresponse"].(float64))
		res.Imf, _ = raw["imf"].(string)
		res.Message, _ = raw["message"].(string)
		res.MsgID, _ = raw["msgId"].(string)
		res.Notification, _ = raw["notification"].(string)
		res.Timestamp = uint(raw["timestamp"].(float64))

	}

	return res
}

type FetchEventsData struct {
	PollTime        int
	Ts              uint
	FetchBaseURL    string
	TimeToNextFetch int
	Events          []Event
}

type FetchEventsResponse struct {
	StatusCode int
	StatusText string
}

type APIResponse struct {
	StatusCode int      `json:"statusCode"`
	StatusText string   `json:"statusText"`
	Data       MapOfAny `json:"Data"`
}

type RawAPIResponse struct {
	Response   APIResponse `json:"response"`
	ParseError error
}

func ParseResponse(contents []byte) *RawAPIResponse {
	raw := &RawAPIResponse{}
	raw.ParseError = json.Unmarshal(contents, raw)

	return raw
}

func (r *RawAPIResponse) AsFetchEvents() (*FetchEventsData, error) {
	if r.ParseError != nil {
		return nil, r.ParseError
	} else if r.Response.StatusCode != 200 {
		return nil, fmt.Errorf("API Failed: %d %s", r.Response.StatusCode, r.Response.StatusText)
	}

	data := r.Response.Data
	rawEvents := data["events"].([]interface{})
	events := make([]Event, 0, len(rawEvents))

	for _, raw := range rawEvents {
		rawEvt := raw.(map[string]interface{})
		events = append(events, Event{
			Type:      rawEvt["type"].(string),
			SeqNum:    uint(rawEvt["seqNum"].(float64)),
			EventData: MapOfAny(rawEvt["eventData"].(map[string]interface{})),
		})
	}

	return &FetchEventsData{
		PollTime:        int(data["pollTime"].(float64)),
		Ts:              uint(data["ts"].(float64)),
		FetchBaseURL:    data["fetchBaseURL"].(string),
		TimeToNextFetch: int(data["timeToNextFetch"].(float64)),
		Events:          events,
	}, nil
}

type Buddy struct {
	AimId     string `json:"aim_id"`
	DisplayId string `json:"display_id"`
	Friendly  string `json:"friendly"`
	UserType  string `json:"user_type"`
}

func (g *Buddy) Parse(raw map[string]interface{}) {
	g.AimId = raw["aimId"].(string)
	g.DisplayId = raw["displayId"].(string)
	g.Friendly = raw["friendly"].(string)
	g.UserType = raw["userType"].(string)
}

type BuddyGroup struct {
	Id      uint64  `json:"id"`
	Name    string  `json:"name"`
	Buddies []Buddy `json:"buddies"`
}

func (g *BuddyGroup) Parse(raw map[string]interface{}) {
	g.Id = uint64(raw["id"].(float64))
	g.Name = raw["name"].(string)

	buddies := raw["buddies"].([]interface{})
	g.Buddies = make([]Buddy, len(buddies))

	for i, b := range buddies {
		g.Buddies[i].Parse(b.(map[string]interface{}))
	}
}

type BuddyList struct {
	Groups []BuddyGroup `json:"groups"`
}

func (bl *BuddyList) Norm() []*Buddy {
	list := make([]*Buddy, 0, 10)

	for _, group := range bl.Groups {
		for _, buddy := range group.Buddies {
			copy := buddy
			list = append(list, &copy)
		}
	}

	return list
}

func (r *RawAPIResponse) AsBuddyList() (*BuddyList, error) {
	list := &BuddyList{}
	groups := r.Response.Data["groups"].([]interface{})
	list.Groups = make([]BuddyGroup, len(groups))

	for i, g := range groups {
		list.Groups[i].Parse(g.(map[string]interface{}))
	}

	return list, nil
}
