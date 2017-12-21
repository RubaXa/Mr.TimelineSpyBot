package space

import (
	"../tnt"
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/vmihailenco/msgpack.v2"
	"time"
)

type sTokens struct {
	baseSpace
}

func (t *sTokens) Has(value string) bool {
	token, _ := t.GetByValue(value)
	return token != nil && token.HasId()
}

func (t *sTokens) GetByValue(value string) (*TokensEntry, error) {
	var tokens []TokensEntry

	err := t.space.Select(&tokens, tnt.Request{
		Index:  "token",
		Values: []interface{}{value},
		Limit:  1,
	})

	if err != nil {
		return nil, err
	} else if len(tokens) > 0 {
		return &tokens[0], nil
	}

	return nil, fmt.Errorf("Token not found")
}

func (t *sTokens) Create(pId uint64) (*TokensEntry, error) {
	uid, err := uuid.NewUUID()

	if err != nil {
		return nil, err
	}

	token := &TokensEntry{
		ProjectId: pId,
		Value:     uid.String(),
		TS:        time.Now().Unix(),
	}
	err = t.Save(token)

	if err != nil {
		return nil, err
	}

	return token, nil
}

const TokenEntrySize = 4

type TokensEntry struct {
	tnt.SpaceEntry

	ProjectId uint64 `json:"project_id"`
	Value     string `json:"value"`
	TS        int64  `json:"ts"`
}

func (t *TokensEntry) EncodeMsgpack(e *msgpack.Encoder) error {
	t.InitEncode(e, TokenEntrySize)
	e.EncodeUint64(t.ProjectId)
	e.EncodeString(t.Value)
	e.EncodeInt64(t.TS)
	return nil
}

func (t *TokensEntry) DecodeMsgpack(d *msgpack.Decoder) error {
	t.InitDecode(d, TokenEntrySize)
	t.ProjectId, _ = d.DecodeUint64()
	t.Value, _ = d.DecodeString()
	t.TS, _ = d.DecodeInt64()
	return nil
}
