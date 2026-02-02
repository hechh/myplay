package player

import (
	"sync/atomic"
	"time"

	"github.com/hechh/framework"
	"github.com/hechh/framework/handler"
	"github.com/hechh/framework/packet"
	"github.com/hechh/library/uerror"
	"google.golang.org/protobuf/proto"
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
	code, msg := irsp.GetRspHead()
	if code != 0 {
		return nil, uerror.Err(code, msg)
	}
	d.rsp.Store(pack)
	return irsp, nil
}
