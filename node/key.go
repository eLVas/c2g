package node

import (
	"errors"
	"fmt"
	"github.com/c2stack/c2g/meta"
)

func ReadKeys(sel Selection) (values []*Value, err error) {
	if len(sel.Path.key) > 0 {
		return sel.Path.key, nil
	}
	list := sel.Path.meta.(*meta.List)
	values = make([]*Value, len(list.Key))
	var key *Value
	for i, keyIdent := range list.Key {
		keyMeta := meta.FindByIdent2(sel.Path.meta, keyIdent).(meta.HasDataType)
		r := FieldRequest{
			Request:Request {
				Selection: sel,
			},
			Meta: keyMeta,
		}
		var hnd ValueHandle
		if err = sel.Node.Field(r, &hnd); err != nil {
			return nil, err
		}
		if hnd.Val == nil {
			return nil, errors.New(fmt.Sprint("Key value is nil for ", keyIdent))
		}
		key.Type = keyMeta.GetDataType()
		values[i] = hnd.Val
	}
	return
}

