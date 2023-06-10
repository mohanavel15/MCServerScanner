package mcstatus

import (
	"encoding/binary"
	"fmt"
)

type Packet struct {
	buf []byte
}

func NewPacket() Packet {
	return Packet{
		buf: []byte{},
	}
}

func (p *Packet) VarInt(value int) []byte {
	buf := []byte{}
	for i := 0; i < 5; i++ {
		if value&-0x80 == 0 {
			buf = append(buf, byte(value))
			break
		}
		buf = append(buf, byte((value&0x7F)|0x80))
		value >>= 7
	}
	return buf
}

func (p *Packet) WriteVarInt(value int) {
	p.buf = append(p.buf, p.VarInt(value)...)
}

func (p *Packet) WriteUShort(value uint16) {
	p.buf = binary.BigEndian.AppendUint16(p.buf, value)
}

func (p *Packet) WriteString(value string) {
	p.WriteVarInt(len(value))
	addressBytes := []byte(value)
	p.buf = append(p.buf, addressBytes...)
}

func (p *Packet) ReadVarInt() (int, error) {
	value := 0
	position := 0
	currentByte := 0

	for i := 0; i < 5; i++ {
		currentByte = int(p.buf[0])
		p.buf = p.buf[1:]

		value |= (currentByte & 0x7F) << position

		if (currentByte & 0x80) == 0 {
			break
		}

		position += 7

		if position >= 32 {
			return 0, fmt.Errorf("%s", "VarInt is too big")
		}
	}

	return value, nil
}

func (p *Packet) ReadString() (string, error) {
	length, err := p.ReadVarInt()
	if err != nil {
		return "", err
	}

	byteString := p.buf[:length]
	p.buf = p.buf[length:]
	return string(byteString), nil
}

func (p *Packet) Buffer() []byte {
	buf := []byte{}
	buf = append(buf, p.VarInt(len(p.buf))...)
	buf = append(buf, p.buf...)
	return buf
}

func (p *Packet) Clear() {
	p.buf = []byte{}
}
