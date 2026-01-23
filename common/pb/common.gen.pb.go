/*
* 本代码由pbtool工具生成，请勿手动修改
 */

package pb

import (
	"github.com/golang/protobuf/proto"
)

func (d LoginType) Uint32() uint32 {
	return uint32(d.Number())
}

func (d NodeType) Uint32() uint32 {
	return uint32(d.Number())
}

func (d Platform) Uint32() uint32 {
	return uint32(d.Number())
}

func (d ErrorCode) Uint32() uint32 {
	return uint32(d.Number())
}

func (d *SessionData) ToDB() ([]byte, error) {
	if d == nil {
		return nil, nil
	}
	return proto.Marshal(d)
}

func (d *SessionData) FromDB(val []byte) error {
	if len(val) <= 0 {
		return nil
	}
	return proto.Unmarshal(val, d)
}

func (d *AccountData) ToDB() ([]byte, error) {
	if d == nil {
		return nil, nil
	}
	return proto.Marshal(d)
}

func (d *AccountData) FromDB(val []byte) error {
	if len(val) <= 0 {
		return nil
	}
	return proto.Unmarshal(val, d)
}
