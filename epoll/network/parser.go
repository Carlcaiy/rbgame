package network

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/golang/protobuf/proto"
)

const (
	LenLen  = 2
	HeadLen = 12
)

// 封包 沾包问题
// len + cmd + data
// 2   + 2   + len-2-2

var byteOrder binary.ByteOrder = binary.BigEndian

type Header struct {
	cmd    uint16 // cmd
	dest   uint8  // 目标服务器类型
	hold1  uint8  // 保留1
	hold2  uint16 // 保留2
	rand1  uint8  // 随机
	rand2  uint8  // 随机
	userID uint32 // 用户UID
}

func NewHeader(dest ServerType, uid uint32) *Header {
	return &Header{
		dest:   uint8(dest),
		userID: uid,
	}
}

func (h *Header) Cmd() uint16 {
	return h.cmd
}

func (h *Header) Dest() ServerType {
	return ServerType(h.dest)
}

func (h *Header) UserID() uint32 {
	return h.userID
}

// 提取消息头
func (p *Header) extract(buf []byte) error {
	p.cmd = byteOrder.Uint16(buf[:2])
	p.dest = buf[2]
	p.hold1 = buf[3]
	p.hold2 = byteOrder.Uint16(buf[4:6])
	p.rand1 = buf[6]
	p.rand2 = buf[7]
	p.userID = byteOrder.Uint32(buf[8:])
	return nil
}

// 注入消息头
func (p *Header) inject(buf []byte) {
	byteOrder.PutUint16(buf[:2], p.cmd)
	buf[2] = p.dest
	buf[3] = p.hold1
	byteOrder.PutUint16(buf[4:6], p.hold2)
	buf[6] = p.rand1
	buf[7] = p.rand2
	byteOrder.PutUint32(buf[8:], p.userID)
}

func PutLength(buf []byte, len uint16) {
	byteOrder.PutUint16(buf[0:2], len)
}

func PutCmd(buf []byte, cmd uint16) {
	byteOrder.PutUint16(buf[2:4], cmd)
}

func PutDest(buf []byte, st ServerType) {
	buf[4] = byte(st)
}

func PutUserID(buf []byte, uid uint32) {
	byteOrder.PutUint32(buf[10:], uid)
}

type Message struct {
	Header
	body []byte
	all  []byte
}

func (m *Message) Bytes() []byte {
	return m.all
}

func Parse(r io.Reader) (msg *Message, err error) {
	buflen := make([]byte, 2)
	// 长度
	_, err = io.ReadFull(r, buflen[:2])
	if err != nil {
		return nil, err
	}

	// 检测长度
	len := byteOrder.Uint16(buflen[:2])
	if len < 12 {
		return nil, fmt.Errorf("parse error len:%d", len)
	}

	all := make([]byte, len+2)
	_, err = io.ReadFull(r, all[2:])
	if err != nil {
		return nil, err
	}
	copy(all, buflen)

	msg = new(Message)
	if err := msg.extract(all[2:]); err != nil {
		return nil, err
	}

	msg.all = all
	msg.body = all[LenLen+HeadLen:]
	return
}

func NewMessage(uid uint32, dest ServerType) *Message {
	p := new(Message)
	p.userID = uid
	p.dest = uint8(dest)
	return p
}

func (m *Message) Pack(cmd uint16, pro proto.Message) ([]byte, error) {
	body, err := proto.Marshal(pro)
	if err != nil {
		return nil, err
	}
	length := uint16(HeadLen + len(body))
	buf := make([]byte, length+LenLen)
	byteOrder.PutUint16(buf, length)

	m.cmd = cmd
	m.Header.inject(buf[LenLen:])

	copy(buf[LenLen+HeadLen:], body)
	return buf, nil
}

func Pack(uid uint32, dest ServerType, cmd uint16, pro proto.Message) ([]byte, error) {
	body, err := proto.Marshal(pro)
	if err != nil {
		return nil, err
	}
	length := uint16(HeadLen + len(body))
	buf := make([]byte, length+LenLen)

	byteOrder.PutUint16(buf, length)

	h := NewHeader(dest, uid)
	h.cmd = cmd
	h.inject(buf[LenLen:])

	copy(buf[LenLen+HeadLen:], body)
	return buf, nil
}

func (m *Message) UnPack(pro proto.Message) error {
	return proto.Unmarshal(m.body, pro)
}
