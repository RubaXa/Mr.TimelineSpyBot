package space

import (
	"../tnt"
	"gopkg.in/vmihailenco/msgpack.v2"
)

type sRecords struct {
	baseSpace
}

func (r *sRecords) New() *RecordsEntry {
	return &RecordsEntry{}
}

func (r *sRecords) Select(pId uint64, offset, limit uint32) ([]RecordsEntry, error) {
	var list []RecordsEntry

	err := r.space.Select(&list, tnt.Request{
		Index:  "project",
		Offset: offset,
		Limit:  limit,
		Values: []interface{}{pId},
	})

	if err != nil {
		return nil, err
	}

	return list, nil
}

func (r *sRecords) GetLast() (*RecordsEntry, error) {
	list := []RecordsEntry{}
	err := r.space.EvalTyped(
		"return box.space."+r.space.Name+".index.primary:max();",
		[]interface{}{},
		&list,
	)

	if err != nil {
		return nil, err
	} else if len(list) > 0 {
		return &list[0], nil
	}

	return nil, nil
}

const RecordsEntrySize = 8

type RecordsEntry struct {
	tnt.SpaceEntry

	ProjectId uint64             `json:"project_id"`
	SeqNum    uint               `json:"seq_num"`
	MsgId     string             `json:"msg_id"`
	TS        uint               `json:"ts"`
	Source    RecordsEntrySource `json:"source"`
	Author    RecordsEntryAuthor `json:"author"`
	Body      string             `json:"body"`
}

func (rec *RecordsEntry) EncodeMsgpack(e *msgpack.Encoder) error {
	rec.InitEncode(e, RecordsEntrySize)

	e.EncodeUint64(rec.ProjectId)
	e.EncodeUint(rec.SeqNum)
	e.EncodeString(rec.MsgId)
	e.EncodeUint(rec.TS)

	e.EncodeArrayLen(2)
	e.EncodeString(rec.Source.Id)
	e.EncodeString(rec.Source.Name)

	e.EncodeArrayLen(2)
	e.EncodeString(rec.Author.Login)
	e.EncodeString(rec.Author.Name)

	e.EncodeString(rec.Body)

	return nil
}

func (rec *RecordsEntry) DecodeMsgpack(d *msgpack.Decoder) error {
	rec.InitDecode(d, RecordsEntrySize)

	rec.ProjectId, _ = d.DecodeUint64()
	rec.SeqNum, _ = d.DecodeUint()
	rec.MsgId, _ = d.DecodeString()
	rec.TS, _ = d.DecodeUint()

	d.DecodeArrayLen()
	rec.Source.Id, _ = d.DecodeString()
	rec.Source.Name, _ = d.DecodeString()

	d.DecodeArrayLen()
	rec.Author.Login, _ = d.DecodeString()
	rec.Author.Name, _ = d.DecodeString()

	rec.Body, _ = d.DecodeString()

	return nil
}

type RecordsEntrySource struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type RecordsEntryAuthor struct {
	Login string `json:"login"`
	Name  string `json:"name"`
}
