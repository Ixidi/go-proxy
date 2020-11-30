package conn

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"net"
)

type ConnectionState int

const (
	HandshakeState ConnectionState = 0
	StatusState    ConnectionState = 1
	LoginState     ConnectionState = 2
	PlayState      ConnectionState = 3
)

type Connection interface {
	ReadPacket() (Packet, error)
	WritePacket(packet Packet) error
	UpdateState(state ConnectionState)
	State() ConnectionState
	Disconnect(reason string)
	Threshold() int
	SetThreshold(threshold int)
}

type connection struct {
	conn      net.Conn
	state     ConnectionState
	threshold int
}

func (c *connection) ReadPacket() (Packet, error) {
	var size VarInt
	var err error
	if err = size.Read(c.conn); err != nil {
		return nil, err
	}

	b := make([]byte, size)
	n, err := c.conn.Read(b)
	if err != nil {
		return nil, err
	}
	if n != int(size) {
		return nil, fmt.Errorf("%d bytes were expected, got %d", size, n)
	}

	buff := bytes.NewBuffer(b)
	if c.threshold != 0 {

		var dataLength VarInt
		if err = dataLength.Read(buff); err != nil {
			return nil, err
		}
		if dataLength != 0 {
			if int(dataLength) < c.threshold {
				return nil, errors.New("data length must be equal to or over the threshold")
			}

			var decompressed bytes.Buffer
			compressed, err := zlib.NewReader(buff)
			if err != nil {
				return nil, err
			}

			if _, err = io.Copy(&decompressed, compressed); err != nil {
				return nil, err
			}
			buff = &decompressed
			if err = compressed.Close(); err != nil {
				return nil, err
			}
		}

	}

	var id VarInt
	if err := id.Read(buff); err != nil {
		return nil, err
	}

	return NewPacket(id, buff.Bytes()), nil
}

func (c *connection) WritePacket(packet Packet) error {
	var buff bytes.Buffer
	var err error
	if err = packet.Id().Write(&buff); err != nil {
		return err
	}
	if _, err = buff.Write(packet.Data()); err != nil {
		return err
	}

	var data bytes.Buffer
	if c.threshold > 0 {
		if buff.Len() < c.threshold {
			if err := VarInt(0).Write(&data); err != nil {
				return err
			}
			data.Write(buff.Bytes())
		} else {
			compressed := zlib.NewWriter(&data)
			if err := VarInt(buff.Len()).Write(&data); err != nil {
				return err
			}

			if _, err := compressed.Write(buff.Bytes()); err != nil {
				return err
			}
			if err := compressed.Flush(); err != nil {
				return err
			}
			if err := compressed.Close(); err != nil {
				return err
			}
		}
	} else {
		data.Write(buff.Bytes())
	}

	if err = VarInt(data.Len()).Write(c.conn); err != nil {
		return err
	}
	if _, err = c.conn.Write(data.Bytes()); err != nil {
		return err
	}
	return nil
}

func (c *connection) UpdateState(state ConnectionState) {
	c.state = state
}

func (c *connection) State() ConnectionState {
	return c.state
}

func (c *connection) Disconnect(reason string) {
	c.conn.Close()
}

func (c *connection) Threshold() int {
	return c.threshold
}

func (c *connection) SetThreshold(threshold int) {
	c.threshold = threshold
}

func NewConnection(conn net.Conn) Connection {
	return &connection{
		conn:      conn,
		state:     HandshakeState,
		threshold: 0,
	}
}
