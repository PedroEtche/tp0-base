package common

import (
	"bytes"
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
				0, 0, 0, 15, 'S', 'a', 'n', 't', 'i', 'a', 'g', 'o', ' ', 'L', 'i', 'o', 'n', 'e', 'l',
				0, 0, 0, 5, 'L', 'o', 'r', 'c', 'a',
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
			actual := createMessage(client)

			if !bytes.Equal(actual, tc.expected) {
				t.Errorf("Test %v - %s FAIL:\nexpected: %v\nactual:   %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
