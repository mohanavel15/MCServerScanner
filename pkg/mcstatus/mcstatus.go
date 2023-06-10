package mcstatus

import (
	"fmt"
	"net"
	"time"
)

func Lookup(addr string, port uint16, timeout time.Duration) (string, error) {
	ip := fmt.Sprintf("%s:%d", addr, port)
	socket, err := net.Dial("tcp", ip)
	if err != nil {
		return "", err
	}

	defer socket.Close()

	proto_version := 47
	next_state := 1

	packet := NewPacket()
	packet.WriteVarInt(0)
	packet.WriteVarInt(proto_version)
	packet.WriteString(addr)
	packet.WriteUShort(port)
	packet.WriteVarInt(next_state)

	_, err = socket.Write(packet.Buffer())
	if err != nil {
		return "", err
	}

	packet.Clear()
	packet.WriteVarInt(0)

	_, err = socket.Write(packet.Buffer())
	if err != nil {
		return "", err
	}

	timeout_at := time.Now().Add(timeout)
	socket.SetDeadline(timeout_at)

	buffer := make([]byte, 1024)
	n, err := socket.Read(buffer)
	if err != nil {
		return "", err
	}

	fmt.Println(buffer[:n])

	packet.Clear()
	packet.buf = buffer

	_, err = packet.ReadVarInt()
	if err != nil {
		return "", fmt.Errorf("%s", "Unable to parse packet")
	}

	id, err := packet.ReadVarInt()
	if err != nil {
		return "", fmt.Errorf("%s", "Unable to parse packet")
	}

	if id != 0 {
		return "", fmt.Errorf("%s", "Unexpected packet")
	}

	status, err := packet.ReadString()
	if err != nil {
		return "", fmt.Errorf("%s", "Unable to parse packet")
	}

	return status, nil
}
