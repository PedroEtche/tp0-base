package common

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"net"
)

/*
Procotol for serializing a batch of bets:

- 2 byte: number of bets in the batch (uint16)
- N bytes: bets (serialized using the protocol defined in createMessage)
*/

const (
	MAX_PACKET_SIZE       int  = 1024 * 8            // 8KB
	MAX_BATCH_PACKET_SIZE int  = MAX_PACKET_SIZE - 2 // 2 byte for the number of bets
	NACK                  byte = 0
)

var ErrNACK = errors.New("Received NACK from server")

// CreateBatch create a Batch Packet from a list of bets. It returns the byte slice representing the packet.
// If the Bet slice is too big, it returns a slice with the bets that were not included in the batch and should be sent in the next batch.
func CreateBatch(bets []Bet) ([]byte, []Bet) {
	var buf bytes.Buffer
	var leftBets []Bet
	included := 0

	// Placeholder para cantidad de bets; se pisa al final con el valor real.
	binary.Write(&buf, binary.BigEndian, uint16(0))

	// Agrega bets hasta llenar el paquete o terminar la lista.
	for i, bet := range bets {
		msg := createMessage(&bet)
		if buf.Len()+len(msg) > MAX_BATCH_PACKET_SIZE {
			leftBets = bets[i:]
			break
		}
		_, _ = buf.Write(msg)
		included++
	}

	// Escribe la cantidad real de bets incluidas en este paquete.
	out := buf.Bytes()
	binary.BigEndian.PutUint16(out[0:2], uint16(included))

	return out, leftBets
}

/*
Procotol for serializing a bet:

- 1 bytes: ID (uint8)
- 1 byte: length of name (uint8)
- N bytes: name (string) -> Max 255 characters
- 1 byte: length of last name (uint8)
- M bytes: last name (string) -> Max 255 characters
- 4 bytes: document (uint32)
- 2 bytes: year (uint16)
- 1 byte: month (uint8)
- 1 byte: day (uint8)
- 4 bytes: number (uint32)

Packet lenght in bytes = 15 + len(name) + len(last name)
*/

func createMessage(bet *Bet) []byte {
	var buf bytes.Buffer

	binary.Write(&buf, binary.BigEndian, bet.ID)

	binary.Write(&buf, binary.BigEndian, uint8(len(bet.Name)))
	buf.WriteString(bet.Name)

	binary.Write(&buf, binary.BigEndian, uint8(len(bet.LastName)))
	buf.WriteString(bet.LastName)

	binary.Write(&buf, binary.BigEndian, bet.Document)

	binary.Write(&buf, binary.BigEndian, bet.GetYear())
	binary.Write(&buf, binary.BigEndian, bet.GetMonth())
	binary.Write(&buf, binary.BigEndian, bet.GetDay())

	binary.Write(&buf, binary.BigEndian, bet.Number)

	return buf.Bytes()
}

func sendMessage(conn net.Conn, msg []byte) {
	written := 0

	for written < len(msg) {
		n, err := conn.Write(msg[written:])
		if err != nil {
			log.Fatalf("Error writting to conn: %v", err)
		}
		written += n
	}
}

func recvACK(c net.Conn) error {
	ack, err := bufio.NewReader(c).ReadByte()
	if err != nil {
		return err
	}

	if ack == 0 {
		return ErrNACK
	}

	return nil
}
