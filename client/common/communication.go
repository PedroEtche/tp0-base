package common

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
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
	BATCH                 byte = 0x00
	WINNERS_REQUEST       byte = 0x01
	ACK                   byte = 0x00
	NACK                  byte = 0xFF
)

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

func sendMessage(conn net.Conn, msg []byte) error {
	written := 0

	for written < len(msg) {
		n, err := conn.Write(msg[written:])
		if err != nil {
			log.Errorf("Error writting to conn: %v", err)
			return err
		}
		written += n
	}
	return nil
}

func sendBatchStartRequest(c net.Conn, id uint8) error {
	msg := byte((BATCH << 7) | id)
	return sendMessage(c, []byte{msg})
}

func sendWinnersRequest(c net.Conn, id uint8) error {
	msg := byte((WINNERS_REQUEST << 7) | id)
	return sendMessage(c, []byte{msg})
}

func recvACK(c net.Conn) error {
	ack, err := bufio.NewReader(c).ReadByte()
	if err != nil {
		return err
	}

	if ack == NACK {
		return fmt.Errorf("Received NACK from server")
	}

	return nil
}

/*
Procotol for serializing winners:

- 1 byte: len(winners) or NACK (uint8)
- N bytes: winners (uint32) -> N = len(winners) * 4

1111 1111 is reserved for NACK, so max winners is 254.
*/

func recvWinnersRespond(c net.Conn) ([]uint32, error) {
	header := make([]byte, 1)

	if _, err := io.ReadFull(c, header); err != nil {
		return nil, fmt.Errorf("error reading winner header: %w", err)
	}

	if header[0] == NACK {
		return nil, fmt.Errorf("received NACK from server")
	}

	count := uint8(header[0])
	winners := make([]uint32, 0, count)
	winner := make([]byte, 4)

	var i uint8 = 0
	for ; i < count; i++ {
		if _, err := io.ReadFull(c, winner); err != nil {
			return nil, fmt.Errorf("error reading winner %d: %w", i, err)
		}
		winners = append(winners, binary.BigEndian.Uint32(winner))
	}

	return winners, nil
}
