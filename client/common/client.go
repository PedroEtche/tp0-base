package common

import (
	"encoding/csv"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID             uint8
	ServerAddress  string
	BatchMaxAmount int
}

// Client Entity that encapsulates how
type Client struct {
	config  ClientConfig
	conn    net.Conn
	channel chan os.Signal
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig, channel chan os.Signal) *Client {
	client := &Client{
		config:  config,
		channel: channel,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	if err := c.createClientSocket(); err != nil {
		log.Fatalf("action: connect | result: fail | client_id: %v | error: %v", c.config.ID, err)
	}
	defer c.conn.Close()

	file, err := os.Open(fmt.Sprintf("/agency-%v.csv", c.config.ID))
	if err != nil {
		log.Fatalf(
			"action: create_csv_reader | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}
	defer file.Close()
	reader := csv.NewReader(file)

	pending := make([]Bet, 0, c.config.BatchMaxAmount)
	eof := false
	stopReading := false

	for !eof || len(pending) > 0 {
		if !stopReading {
			pending = c.populateBatch(&eof, pending, reader, &stopReading)
		}

		if len(pending) == 0 {
			break
		}

		left := c.sendBatch(pending)

		if err := recvACK(c.conn); err != nil {
			if err.Error() == "Received NACK from server" {
				log.Error("action: receive_message | result: nack")
			} else {
				log.Errorf(
					"action: receive_message | result: fail | client_id: %v | error: %v",
					c.config.ID,
					err,
				)
			}
			return
		}

		log.Info("action: receive_message | result: ack")
		pending = left

		// ACK means more batches need to be sent
		sendACK(c.conn)
	}
	// NACK means the connection has finish
	sendNACK(c.conn)
}

// sendBatch Send a batch of bets to the server. The method returns the bets that were not sent in the batch
func (c *Client) sendBatch(pending []Bet) []Bet {
	msg, left := CreateBatch(pending)
	sendMessage(c.conn, msg)
	return left
}

func (c *Client) populateBatch(eof *bool, batch []Bet, reader *csv.Reader, stopReading *bool) []Bet {
	for !*eof && len(batch) < c.config.BatchMaxAmount {
		if ListenForSigTerm(c.channel) {
			*stopReading = true
			break
		}

		line, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				*eof = true
				break
			}
			log.Errorf(
				"action: read_csv_line | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			continue
		}

		bet, err := CreatBetFromCSVLine(c.config.ID, line)
		if err != nil {
			log.Errorf(
				"action: create_bet_from_csv_line | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			continue
		}
		batch = append(batch, *bet)
	}

	return batch
}
