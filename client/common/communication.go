package common

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
)

var ErrNACK = errors.New("Received NACK from server")

const NACK byte = 0

func (c *Client) CreateMessage() []byte {
	var buf bytes.Buffer

	buf.WriteByte(c.config.ID)

	binary.Write(&buf, binary.BigEndian, uint8(len(c.config.Name)))
	buf.WriteString(c.config.Name)

	binary.Write(&buf, binary.BigEndian, uint8(len(c.config.LastName)))
	buf.WriteString(c.config.LastName)

	binary.Write(&buf, binary.BigEndian, c.config.Document)

	binary.Write(&buf, binary.BigEndian, c.config.BirthYear)
	buf.WriteByte(c.config.BirthMonth)
	buf.WriteByte(c.config.BirthDay)

	binary.Write(&buf, binary.BigEndian, c.config.Number)

	return buf.Bytes()
}

// SendMessage sends the given message to the server, ensuring that all bytes are sent.
// If connection fails, it logs the error and terminates the program.
func (c *Client) SendMessage(msg []byte) {
	written := 0

	for written < len(msg) {
		n, err := c.conn.Write(msg[written:])
		if err != nil {
			log.Fatalf("Error writting to conn: %v", err)
		}
		written += n
	}
}

// RecvACK reads a single byte from the server to determine if the message was acknowledged (ACK) or not (NACK).
func (c *Client) RecvACK() error {
	ack, err := bufio.NewReader(c.conn).ReadByte()
	if err != nil {
		return err
	}

	if ack == NACK {
		return ErrNACK
	}

	return nil
}
