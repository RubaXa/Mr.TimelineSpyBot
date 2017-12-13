package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/vmihailenco/msgpack.v2"
)

type Any interface{}
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
	raw := data.(map[string]interface{})
	attrs.ChatName = raw["chat_name"].(string)
	attrs.Sender = raw["sender"].(string)
	attrs.SenderName = raw["senderName"].(string)
}

type EventRawMsg struct {
	Base64Msg string
}

func (rawMsg *EventRawMsg) Parse(data interface{}) {
	rawMsg.Base64Msg = data.(map[string]interface{})["base64Msg"].(string)
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
		res.Imf = raw["imf"].(string)
		res.Message = raw["message"].(string)
		res.MsgID = raw["msgId"].(string)
		res.Notification = raw["notification"].(string)
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

type TimelineRecordSource struct {
	Id   string
	Name string
}

func (s TimelineRecordSource) EncodeMsgpack(e *msgpack.Encoder) error {
	if err := e.EncodeArrayLen(2); err != nil {
		return err
	}

	if err := e.EncodeString(s.Id); err != nil {
		return err
	}

	if err := e.EncodeString(s.Name); err != nil {
		return err
	}

	return nil
}

type TimelineRecordAuthor struct {
	Login string
	Name  string
}

func (a TimelineRecordAuthor) EncodeMsgpack(e *msgpack.Encoder) error {
	if err := e.EncodeArrayLen(2); err != nil {
		return err
	}

	if err := e.EncodeString(a.Login); err != nil {
		return err
	}

	if err := e.EncodeString(a.Name); err != nil {
		return err
	}

	return nil
}

type TimelineRecord struct {
	Id     string
	TS     uint
	Source TimelineRecordSource
	Author TimelineRecordAuthor
	Body   string
}

func (rec TimelineRecord) EncodeMsgpack(e *msgpack.Encoder) error {
	if err := e.EncodeArrayLen(5); err != nil {
		return err
	}

	if err := e.EncodeString(rec.Id); err != nil {
		return err
	}

	if err := e.EncodeUint(rec.TS); err != nil {
		return err
	}

	e.Encode(rec.Source)
	e.Encode(rec.Author)

	if err := e.EncodeString(rec.Body); err != nil {
		return err
	}

	return nil
}

type Project struct {
	Id     uint
	Name   string
	LastId uint
}
