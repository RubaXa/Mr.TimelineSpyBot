package tnt

import (
	"github.com/tarantool/go-tarantool"
	"time"
)

type Box struct {
	Client  *tarantool.Connection
	options tarantool.Opts
	spaces  map[string]*Space
}

func (box *Box) GetNextId() (uint64, error) {
	res, err := box.Client.Eval(
		"return box.sequence.GID:next()",
		[]interface{}{},
	)

	if err != nil {
		return 0, err
	}

	return res.Data[0].(uint64), nil
}

func (box *Box) GetSpace(name string) *Space {
	if _, ok := box.spaces[name]; !ok {
		box.spaces[name] = &Space{name, box}
	}

	return box.spaces[name]
}

func InitBox(server, user, pass string, timeout time.Duration) (*Box, error) {
	if timeout == 0 {
		timeout = 500
	}

	opts := tarantool.Opts{
		Timeout:       timeout * time.Millisecond,
		Reconnect:     1 * time.Second,
		MaxReconnects: 3,
		User:          user,
		Pass:          pass,
	}

	client, err := tarantool.Connect(server, opts)

	if err != nil {
		return nil, err
	}

	return &Box{
		Client:  client,
		options: opts,
		spaces:  make(map[string]*Space),
	}, nil
}
