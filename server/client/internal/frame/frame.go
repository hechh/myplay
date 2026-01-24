package frame

import (
	"encoding/binary"
	"myplay/common/pb"

	"github.com/hechh/framework"
	"github.com/hechh/framework/handler"
	"github.com/hechh/framework/packet"
	"github.com/hechh/library/uerror"
)

type Frame struct{}

// 解码数据包
func (d *Frame) Decode(buf []byte) (*packet.Packet, error) {
	msg := &packet.Packet{
		Head: &packet.Head{
			SendType:    packet.SendType_POINT,
			SrcNodeType: framework.GetSelfType(),
			SrcNodeId:   framework.GetSelfId(),
			Cmd:         binary.BigEndian.Uint32(buf[4:]),
			Id:          binary.BigEndian.Uint64(buf[8:]),
			ActorId:     binary.BigEndian.Uint64(buf[16:]),
			Seq:         binary.BigEndian.Uint32(buf[24:]),
			Version:     binary.BigEndian.Uint32(buf[28:]),
			Extra:       binary.BigEndian.Uint32(buf[32:]),
		},
		Body: buf[36:],
	}

	// 获取rpc
	hh := handler.GetCmdRpc(msg.Head.Cmd)
	if hh == nil {
		return msg, uerror.Err(pb.ErrorCode_CmdNotSupported, "Cmd(%d)没有注册", msg.Head.Cmd)
	}
	msg.Head.DstNodeType = hh.GetNodeType()
	msg.Head.ActorFunc = hh.GetCrc32()
	return msg, nil
}

func (d *Frame) Encode(pack *packet.Packet) (buf []byte) {
	size := uint32(32 + len(pack.Body))
	buf = make([]byte, size+4)
	// 组包
	pos := 0
	binary.BigEndian.PutUint32(buf, size)
	pos += 4
	binary.BigEndian.PutUint32(buf[pos:], pack.Head.Cmd) // cmd
	pos += 4
	binary.BigEndian.PutUint64(buf[pos:], pack.Head.Id) // uid
	pos += 8
	binary.BigEndian.PutUint64(buf[pos:], 0) // router id
	pos += 8
	binary.BigEndian.PutUint32(buf[pos:], pack.Head.Seq) // seq
	pos += 4
	binary.BigEndian.PutUint32(buf[pos:], 0) // version
	pos += 4
	binary.BigEndian.PutUint32(buf[pos:], 0) // extra
	pos += 4
	copy(buf[pos:], pack.Body)
	return
}
