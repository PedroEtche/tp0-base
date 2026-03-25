package common

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
	"strings"
	"testing"
)

func TestCreateMessage(t *testing.T) {
	tests := []struct {
		name     string
		bet      Bet
		expected []byte
	}{
		{
			name: "create message with valid config",
			bet: Bet{
				ID:        1,
				Name:      "Santiago Lionel",
				LastName:  "Lorca",
				Document:  12345678,
				Birthdate: "1990-01-01",
				Number:    1,
			},
			expected: []byte{
				1,
				15, 'S', 'a', 'n', 't', 'i', 'a', 'g', 'o', ' ', 'L', 'i', 'o', 'n', 'e', 'l',
				5, 'L', 'o', 'r', 'c', 'a',
				0, 188, 97, 78, // DNI: 12345678
				7, 198, // Year: 1990
				1,          // Month
				1,          // Day
				0, 0, 0, 1, // Number
			},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := createMessage(&tc.bet)

			if !bytes.Equal(actual, tc.expected) {
				t.Errorf("Test %v - %s FAIL:\nexpected: %v\nactual:   %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestRecvACK(t *testing.T) {
	tests := []struct {
		name        string
		ackByte     []byte
		closeBefore bool
		wantErr     bool
		wantErrIs   error
	}{
		{name: "ack byte must succeed", ackByte: []byte{1}, wantErr: false},
		{name: "nack byte must fail", ackByte: []byte{NACK}, wantErr: true, wantErrIs: ErrNACK},
		{name: "eof before ack must fail", closeBefore: true, wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			clientConn, serverConn := net.Pipe()
			defer clientConn.Close()

			if tc.closeBefore {
				serverConn.Close()
			} else {
				go func() {
					defer serverConn.Close()
					_, _ = serverConn.Write(tc.ackByte)
				}()
			}

			err := recvACK(clientConn)
			if (err != nil) != tc.wantErr {
				t.Fatalf("unexpected error state. err=%v, wantErr=%v", err, tc.wantErr)
			}

			if tc.wantErrIs != nil && !errors.Is(err, tc.wantErrIs) {
				t.Fatalf("expected error %v, got %v", tc.wantErrIs, err)
			}
		})
	}
}

func TestCreateBatch(t *testing.T) {
	bets := []Bet{
		{
			ID:        1,
			Name:      "Ana",
			LastName:  "Lopez",
			Document:  12345678,
			Birthdate: "1990-01-01",
			Number:    7,
		},
		{
			ID:        2,
			Name:      "Juan",
			LastName:  "Perez",
			Document:  87654321,
			Birthdate: "1985-12-31",
			Number:    10,
		},
	}

	packet, left := CreateBatch(bets)

	if len(left) != 0 {
		t.Fatalf("expected no pending bets, got %d", len(left))
	}

	if len(packet) < 2 {
		t.Fatalf("expected packet with header, got len=%d", len(packet))
	}

	count := binary.BigEndian.Uint16(packet[0:2])
	if count != uint16(len(bets)) {
		t.Fatalf("expected %d bets in header, got %d", len(bets), count)
	}

	msg1 := createMessage(&bets[0])
	msg2 := createMessage(&bets[1])
	expectedPayload := append(msg1, msg2...)
	if !bytes.Equal(packet[2:], expectedPayload) {
		t.Fatalf("batch payload mismatch")
	}
}

func TestCreateBatchBoundaryLeavesOneBetWhenNextExceedsMaxPacket(t *testing.T) {
	baseBet := Bet{
		ID:        1,
		Name:      strings.Repeat("A", 250),
		LastName:  strings.Repeat("B", 250),
		Document:  12345678,
		Birthdate: "1990-01-01",
		Number:    7,
	}

	bets := make([]Bet, 16)
	for i := range bets {
		bets[i] = baseBet
	}

	packet, left := CreateBatch(bets)

	if got := binary.BigEndian.Uint16(packet[:2]); got != 15 {
		t.Fatalf("expected 15 bets in batch header, got %d", got)
	}

	if len(left) != 1 {
		t.Fatalf("expected 1 pending bet, got %d", len(left))
	}

	if len(packet) > MAX_BATCH_PACKET_SIZE {
		t.Fatalf("batch packet exceeds max size: got=%d max=%d", len(packet), MAX_BATCH_PACKET_SIZE)
	}
}

func TestRecvWinnersRespondOK(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()

	go func() {
		defer serverConn.Close()
		// Header: 2 winners. Then each winner as uint32 (4 bytes).
		_, _ = serverConn.Write([]byte{
			2,
			0, 0, 0, 7,
			0, 0, 0, 10,
		})
	}()

	got, err := recvWinnersRespond(clientConn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []uint32{7, 10}
	if len(got) != len(want) {
		t.Fatalf("unexpected winners count. got=%d, want=%d", len(got), len(want))
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("unexpected winner at position %d. got=%d, want=%d", i, got[i], want[i])
		}
	}
}

func TestRecvWinnersRespondNACK(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()

	go func() {
		defer serverConn.Close()
		_, _ = serverConn.Write([]byte{NACK})
	}()

	_, err := recvWinnersRespond(clientConn)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !errors.Is(err, ErrNACK) {
		t.Fatalf("expected error %v, got %v", ErrNACK, err)
	}
}

func TestRecvWinnersRespondEOFBeforeHeader(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()

	// Server closes immediately, so client cannot read header byte.
	serverConn.Close()

	_, err := recvWinnersRespond(clientConn)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestRecvWinnersRespondEOFWhileReadingWinner(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()

	go func() {
		defer serverConn.Close()
		// Says there are 2 winners, but only sends 1 complete and 2 bytes of the next one.
		_, _ = serverConn.Write([]byte{
			2,
			0, 0, 0, 7,
			0, 0,
		})
	}()

	_, err := recvWinnersRespond(clientConn)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
