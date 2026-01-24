package player

import (
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hechh/framework"
	"github.com/hechh/framework/handler"
	"github.com/hechh/framework/packet"
	"github.com/hechh/library/uerror"
)

type Request struct {
	cmd   uint32
	seq   uint32
	start int64
	end   int64
	rsp   atomic.Pointer[packet.Packet]
}

func NewRequest(cmd uint32, seq uint32) *Request {
	return &Request{
		cmd: cmd,
		seq: seq,
		rsp: atomic.Pointer[packet.Packet]{},
	}
}

func (d *Request) Start() *Request {
	atomic.StoreInt64(&d.end, time.Now().Unix())
	return d
}

func (d *Request) End() *Request {
	atomic.StoreInt64(&d.end, time.Now().Unix())
	return d
}

func (d *Request) GetRsp(pack *packet.Packet) (proto.Message, error) {
	hh := handler.GetCmdRpc(d.cmd)
	rsp := hh.New(1)
	if err := hh.Unmarshal(pack.Body, rsp); err != nil {
		return nil, err
	}
	irsp, _ := rsp.(framework.IResponse)
	msg := irsp.GetRspHead()
	if msg != nil {
		return nil, uerror.Err(msg.Code, msg.Msg)
	}
	d.rsp.Store(pack)
	return irsp, nil
}
