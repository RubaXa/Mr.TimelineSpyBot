package space

import (
	"../tnt"
	"gopkg.in/vmihailenco/msgpack.v2"
)

type sRecords struct {
	baseSpace
}

func (r *sRecords) GetLast() *RecordsEntry {
	r.space.
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

	e.EncodeUint64(rec.Id)
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
