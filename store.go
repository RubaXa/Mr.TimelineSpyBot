package main

import (
	"github.com/tarantool/go-tarantool"
	"time"
)

type Store struct {
	Projects *Space
	Records  *Space
}

type Space struct {
	Name   string
	client *tarantool.Connection
}

func (space *Space) Insert(data TimelineRecord) (*tarantool.Response, error) {
	return space.client.Eval("box.space."+space.Name+":auto_increment({...})", data)
}

func InitSpace(client *tarantool.Connection, name string) *Space {
	return &Space{name, client}
}

func InitStore(cfg map[string]string) (*Store, error) {
	opts := tarantool.Opts{
		Timeout:       500 * time.Millisecond,
		Reconnect:     1 * time.Second,
		MaxReconnects: 3,
		User:          cfg["user"],
		Pass:          cfg["pass"],
	}

	client, err := tarantool.Connect(cfg["server"], opts)

	if err != nil {
		return nil, err
	}

	return &Store{
		Projects: InitSpace(client, cfg["space.projects"]),
		Records:  InitSpace(client, cfg["space.records"]),
	}, nil
}
