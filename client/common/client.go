package common

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

const WAIT_TIME = time.Millisecond * 5000

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

// sendBatch Send a batch of bets to the server. The method returns the bets that were not sent in the batch
func (c *Client) sendBatch(pending []Bet) ([]Bet, error) {
	msg, left := CreateBatch(pending)
	return left, sendMessage(c.conn, msg)
}

func sendAllBatches(c *Client) bool {
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

	if err := c.createClientSocket(); err != nil {
		log.Fatalf("action: connect | result: fail | client_id: %v | error: %v", c.config.ID, err)
	}
	defer c.conn.Close()

	if err := sendBatchStartRequest(c.conn, c.config.ID); err != nil {
		log.Errorf("action: send_batch_start_request | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return true
	}

	for !eof || len(pending) > 0 {
		if !stopReading {
			pending = c.populateBatch(&eof, pending, reader, &stopReading)
		}

		if len(pending) == 0 {
			break
		}

		left, err := c.sendBatch(pending)
		if err != nil {
			log.Errorf("action: send_batch | result: fail | client_id: %v | error: %v", c.config.ID, err)
			return true
		}

		if err := recvACK(c.conn); err != nil {
			if errors.Is(err, ErrNACK) {
				log.Error("action: receive_message | result: fail")
				// Corrupted batch. Try next batch
				pending = pending[:0]
				continue
			} else {
				log.Errorf(
					"action: receive_message | result: fail | client_id: %v | error: %v",
					c.config.ID,
					err,
				)
				return true
			}
		}

		log.Info("action: receive_message | result: success")
		pending = left
	}
	return false
}

func showWinners(c *Client) {
	// Polls the server for the winners
	for {
		if err := c.createClientSocket(); err != nil {
			log.Fatalf("action: connect | result: fail | client_id: %v | error: %v", c.config.ID, err)
		}

		if err := sendWinnersRequest(c.conn, c.config.ID); err != nil {
			log.Errorf("action: send_winners_request | result: fail | client_id: %v | error: %v", c.config.ID, err)
			break
		}

		winners, err := recvWinnersRespond(c.conn)
		if err != nil {
			if errors.Is(err, ErrNACK) {
				log.Info("action: consultar_ganadores | description: todavia se esperan apuestas")
				c.conn.Close()
				time.Sleep(WAIT_TIME)
				continue
			} else {
				log.Errorf(
					"action: receive_message | result: fail | client_id: %v | error: %v",
					c.config.ID,
					err,
				)
				break
			}
		}

		log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %v", len(winners))
		break
	}
	c.conn.Close()
}

// StartClientLoop Send bets to the Server an the polls for the winner
func (c *Client) StartClientLoop() {
	shouldReturn := sendAllBatches(c)
	if shouldReturn {
		return
	}
	log.Info("action: apuestas_enviadas | result: success")

	showWinners(c)
}
