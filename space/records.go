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

const RecordsEntrySize = 7

type RecordsEntry struct {
	tnt.SpaceEntry

	SeqNum uint
	MsgId  string
	TS     uint
	Source RecordsEntrySource
	Author RecordsEntryAuthor
	Body   string
}

func (rec *RecordsEntry) EncodeMsgpack(e *msgpack.Encoder) error {
	rec.InitEncode(e, RecordsEntrySize)

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

	rec.SeqNum, _ = d.DecodeUint()
	rec.MsgId, _ = d.DecodeString()
	rec.TS, _ = d.DecodeUint()

	d.DecodeArrayLen()
	rec.Source.Id, _ = d.DecodeString()
	rec.Source.Name, _ = d.DecodeString()

	d.DecodeArrayLen()
	rec.Author.Login, _ = d.DecodeString()
	rec.Source.Name, _ = d.DecodeString()

	d.DecodeArrayLen()
	rec.Body, _ = d.DecodeString()

	return nil
}

type RecordsEntrySource struct {
	Id   string
	Name string
}

type RecordsEntryAuthor struct {
	Login string
	Name  string
}
