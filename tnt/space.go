package tnt

import (
	"fmt"
	"github.com/tarantool/go-tarantool"
	"gopkg.in/vmihailenco/msgpack.v2"
	"math"
)

type Space struct {
	Name string
	box  *Box
}

func (space *Space) Eval(code string, args []interface{}) (*tarantool.Response, error) {
	res, err := space.box.Client.Eval(code, args)

	if err != nil {
		return nil, fmt.Errorf("[tnt] box.space.%s:eval() failed: %s", space.Name, err.Error())
	}

	return res, nil
}

func (space *Space) EvalTyped(code string, args []interface{}, res interface{}) error {
	err := space.box.Client.EvalTyped(code, args, &res)

	if err != nil {
		return fmt.Errorf("[tnt] box.space.%s:evalTyped() failed: %s", space.Name, err.Error())
	}

	return nil
}

func (space *Space) Insert(entry ISpaceEntry) error {
	if !entry.HasId() {
		if id, err := space.box.GetNextId(); err != nil {
			return err
		} else {
			entry.SetId(id)
		}
	}

	_, err := space.box.Client.Insert(space.Name, entry)

	if err != nil {
		return fmt.Errorf("[tnt] box.space.%s:insert() failed: %s", space.Name, err.Error())
	}

	return nil
}

func (space *Space) Update(entry ISpaceEntry) error {
	if !entry.HasId() {
		if id, err := space.box.GetNextId(); err != nil {
			return err
		} else {
			entry.SetId(id)
		}
	}

	_, err := space.box.Client.Update(
		space.Name,
		"primary",
		[]interface{}{entry.GetId()},
		entry,
	)

	if err != nil {
		return fmt.Errorf("[tnt] box.space.%s:update() failed: %s", space.Name, err.Error())
	}

	return nil
}

func (space *Space) SelectOne(id uint64, entry ISpaceEntry) error {
	err := space.box.Client.GetTyped(
		space.Name,
		"primary",
		[]interface{}{id},
		&entry,
	)

	return err
}

func (space *Space) SelectAll(list *[]interface{}) error {
	return space.box.Client.SelectTyped(
		space.Name,
		"primary",
		0,
		math.MaxInt32,
		tarantool.IterAll,
		[]interface{}{},
		list,
	)
}

type ISpaceEntry interface {
	HasId() bool
	GetId() uint64
	SetId(id uint64)
}

type SpaceEntry struct {
	Id uint64 `json:"id"`
}

func (entry *SpaceEntry) HasId() bool {
	return entry.Id > 0
}

func (entry *SpaceEntry) GetId() uint64 {
	return entry.Id
}

func (entry *SpaceEntry) SetId(id uint64) {
	entry.Id = id
}

func (entry *SpaceEntry) InitEncode(e *msgpack.Encoder, s int) {
	if err := e.EncodeArrayLen(s); err != nil {
		panic(err)
	}

	if err := e.EncodeUint64(entry.Id); err != nil {
		panic(err)
	}
}

func (entry *SpaceEntry) InitDecode(d *msgpack.Decoder, s int) {
	var n int
	var err error

	if n, err = d.DecodeArrayLen(); err != nil {
		panic(err)
	} else if n != s {
		panic(fmt.Errorf("SpaceEntry size doesn't match: %d != %d", n, s))
	}

	if entry.Id, err = d.DecodeUint64(); err != nil {
		panic(err)
	}
}
