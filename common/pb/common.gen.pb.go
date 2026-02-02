/*
* 本代码由pbtool工具生成，请勿手动修改
 */

package pb

import (
	"google.golang.org/protobuf/proto"
)

func (d CMD) Integer() uint32 {
	return uint32(d.Number())
}

func (d PlayerDataType) Integer() uint32 {
	return uint32(d.Number())
}

func (d LoginType) Integer() uint32 {
	return uint32(d.Number())
}

func (d NodeType) Integer() uint32 {
	return uint32(d.Number())
}

func (d Platform) Integer() uint32 {
	return uint32(d.Number())
}

func (d ErrorCode) Integer() uint32 {
	return uint32(d.Number())
}

func (d *AuthRsp) SetRspHead(code int32, msg string) {
	d.Head = &RspHead{Code: code, Msg: msg}
}

func (d *AuthRsp) GetRspHead() (int32, string) {
	return d.Head.Code, d.Head.Msg
}

func (d *LoginRsp) SetRspHead(code int32, msg string) {
	d.Head = &RspHead{Code: code, Msg: msg}
}

func (d *LoginRsp) GetRspHead() (int32, string) {
	return d.Head.Code, d.Head.Msg
}

func (d *HeartRsp) SetRspHead(code int32, msg string) {
	d.Head = &RspHead{Code: code, Msg: msg}
}

func (d *HeartRsp) GetRspHead() (int32, string) {
	return d.Head.Code, d.Head.Msg
}

func (d *KickNotify) ToDB() ([]byte, error) {
	if d == nil {
		return nil, nil
	}
	return proto.Marshal(d)
}

func (d *KickNotify) FromDB(val []byte) error {
	if len(val) <= 0 {
		return nil
	}
	return proto.Unmarshal(val, d)
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

func (d *PlayerData) ToDB() ([]byte, error) {
	if d == nil {
		return nil, nil
	}
	return proto.Marshal(d)
}

func (d *PlayerData) FromDB(val []byte) error {
	if len(val) <= 0 {
		return nil
	}
	return proto.Unmarshal(val, d)
}

func (d *PlayerBaseData) ToDB() ([]byte, error) {
	if d == nil {
		return nil, nil
	}
	return proto.Marshal(d)
}

func (d *PlayerBaseData) FromDB(val []byte) error {
	if len(val) <= 0 {
		return nil
	}
	return proto.Unmarshal(val, d)
}

func (d *PlayerBagData) ToDB() ([]byte, error) {
	if d == nil {
		return nil, nil
	}
	return proto.Marshal(d)
}

func (d *PlayerBagData) FromDB(val []byte) error {
	if len(val) <= 0 {
		return nil
	}
	return proto.Unmarshal(val, d)
}

func (d *ItemData) ToDB() ([]byte, error) {
	if d == nil {
		return nil, nil
	}
	return proto.Marshal(d)
}

func (d *ItemData) FromDB(val []byte) error {
	if len(val) <= 0 {
		return nil
	}
	return proto.Unmarshal(val, d)
}

func (d *RspHead) ToDB() ([]byte, error) {
	if d == nil {
		return nil, nil
	}
	return proto.Marshal(d)
}

func (d *RspHead) FromDB(val []byte) error {
	if len(val) <= 0 {
		return nil
	}
	return proto.Unmarshal(val, d)
}
