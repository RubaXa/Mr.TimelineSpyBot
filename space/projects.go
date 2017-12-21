package space

import (
	"../tnt"
	"fmt"
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

func (p *sProjects) GetByTokenValue(t string) (*ProjectsEntry, error) {
	token, err := Tokens.GetByValue(t)

	if err != nil {
		return nil, err
	}

	project := p.Get(token.ProjectId)

	if project == nil {
		return nil, fmt.Errorf("Project not found")
	}

	return project, nil
}

func (p *sProjects) GetAll() ([]ProjectsEntry, error) {
	var list []ProjectsEntry
	err := p.space.SelectAll(&list)

	if err != nil {
		return nil, err
	}

	return list, nil
}

const ProjectsEntrySize = 4

type ProjectsEntry struct {
	tnt.SpaceEntry

	Name  string   `json:"name"`
	Key   string   `json:"key"`
	Chats []string `json:"chats"`
}

func (entry *ProjectsEntry) HasChat(id string) bool {
	for _, x := range entry.Chats {
		if x == id {
			return true
		}
	}

	return false
}

func (entry *ProjectsEntry) AddChat(id string) {
	entry.Chats = append(entry.Chats, id)
}

func (entry *ProjectsEntry) RemoveChat(id string) {
	for i, x := range entry.Chats {
		if x == id {
			entry.Chats = append(entry.Chats[:i], entry.Chats[i+1:]...)
			break
		}
	}
}

func (entry *ProjectsEntry) EncodeMsgpack(e *msgpack.Encoder) (err error) {
	entry.InitEncode(e, ProjectsEntrySize)

	e.EncodeString(entry.Name)
	e.EncodeString(entry.Key)

	e.EncodeArrayLen(len(entry.Chats))

	for _, id := range entry.Chats {
		e.EncodeString(id)
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
			entry.Chats[i], _ = d.DecodeString()
		}
	}

	return nil
}
