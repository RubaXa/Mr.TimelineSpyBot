package space

import (
	"../tnt"
)

type baseSpace struct {
	space *tnt.Space
}

func (bs *baseSpace) Init(space *tnt.Space) {
	bs.space = space
}

func (bs *baseSpace) Save(entry tnt.ISpaceEntry) error {
	if entry.HasId() {
		return bs.space.Update(entry)
	} else {
		return bs.space.Insert(entry)
	}
}

func (bs *baseSpace) GetLastId() uint64 {
	res, _ := bs.space.Eval(
		"return box.space."+bs.space.Name+".index.primary:max()",
		[]interface{}{},
	)

	if res.Data != nil {
		return res.Data[0].([]interface{})[0].(uint64)
	}

	return 0
}
