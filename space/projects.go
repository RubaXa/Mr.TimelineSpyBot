package space

import (
	"../tnt"
	"github.com/google/uuid"
	"gopkg.in/vmihailenco/msgpack.v2"
)

type sProjects struct {
	baseSpace
}

func (p *sProjects) New() *ProjectsEntry {
	return &ProjectsEntry{
		Key:   uuid.New().String(),
		Chats: make([]string, 0),
	}
}

func (p *sProjects) Get(id uint64) *ProjectsEntry {
	entry := p.New()
	p.space.SelectOne(id, entry)

	if entry.HasId() {
		return entry
	} else {
		return nil
	}
}

const ProjectsEntrySize = 4

type ProjectsEntry struct {
	tnt.SpaceEntry

	Name  string   `json:"name"`
	Key   string   `json:"key"`
	Chats []string `json:"chats"`
}

func (entry *ProjectsEntry) EncodeMsgpack(e *msgpack.Encoder) (err error) {
	entry.InitEncode(e, ProjectsEntrySize)

	e.EncodeString(entry.Name)
	e.EncodeString(entry.Key)

	e.EncodeArrayLen(len(entry.Chats))
	for _, uin := range entry.Chats {
		e.EncodeString(uin)
	}

	return nil
}

func (entry *ProjectsEntry) DecodeMsgpack(d *msgpack.Decoder) (err error) {
	entry.InitDecode(d, ProjectsEntrySize)

	entry.Name, _ = d.DecodeString()
	entry.Key, _ = d.DecodeString()

	if n, err := d.DecodeArrayLen(); err != nil {
		panic(err)
	} else {
		entry.Chats = make([]string, n)
		for i := 0; i < n; i++ {
			d.Decode(&entry.Chats[i])
		}
	}

	return nil
}
