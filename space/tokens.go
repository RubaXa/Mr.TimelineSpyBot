package space

import (
	"../tnt"
	"gopkg.in/vmihailenco/msgpack.v2"
)

type sTokens struct {
	baseSpace
}

func (t *sTokens) Create(pId uint64) (*TokensEntry, error) {
	token := &TokensEntry{ProjectId: pId}
	err := t.Save(token)

	if err != nil {
		return nil, err
	}

	return token, nil
}

const TokenEntrySize = 3

type TokensEntry struct {
	tnt.SpaceEntry
	ProjectId uint64 `json:"project_id"`
	Value     string `json:"value"`
}

func (t *TokensEntry) EncodeMsgpack(e *msgpack.Encoder) error {
	t.InitEncode(e, TokenEntrySize)
	e.EncodeUint64(t.ProjectId)
	e.EncodeString(t.Value)
	return nil
}

func (t *TokensEntry) DecodeMsgpack(d *msgpack.Decoder) error {
	t.InitDecode(d, TokenEntrySize)
	t.ProjectId, _ = d.DecodeUint64()
	t.Value, _ = d.DecodeString()
	return nil
}
