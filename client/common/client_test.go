package common

import (
	"bytes"
	"net"
	"testing"
)

func TestCreateMessage(t *testing.T) {
	tests := []struct {
		name     string
		config   ClientConfig
		expected []byte
	}{
		{
			name: "create message with valid config",
			config: ClientConfig{
				ID:            1,
				ServerAddress: "localhost:8080",
				Name:          "Santiago Lionel",
				LastName:      "Lorca",
				Document:      12345678,
				BirthYear:     1990,
				BirthMonth:    1,
				BirthDay:      1,
				Number:        1,
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
			client := NewClient(tc.config)
			actual := client.CreateMessage()

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
	}{
		{name: "ack byte must succeed", ackByte: []byte{1}, wantErr: false},
		{name: "nack byte must fail", ackByte: []byte{0}, wantErr: true},
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

			c := &Client{conn: clientConn}
			err := c.RecvACK()
			if (err != nil) != tc.wantErr {
				t.Fatalf("unexpected error state. err=%v, wantErr=%v", err, tc.wantErr)
			}
		})
	}
}
