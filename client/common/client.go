package common

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            uint8
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
	Name          string
	LastName      string
	Document      uint32
	BirthYear     uint16
	BirthMonth    uint8
	BirthDay      uint8
	Number        uint32
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return err
	}
	c.conn = conn
	return nil
}

func createMessage(c *Client) []byte {
	var buf bytes.Buffer

	buf.WriteByte(c.config.ID)

	binary.Write(&buf, binary.BigEndian, uint32(len(c.config.Name)))
	buf.WriteString(c.config.Name)

	binary.Write(&buf, binary.BigEndian, uint32(len(c.config.LastName)))
	buf.WriteString(c.config.LastName)

	binary.Write(&buf, binary.BigEndian, c.config.Document)

	binary.Write(&buf, binary.BigEndian, c.config.BirthYear)
	buf.WriteByte(c.config.BirthMonth)
	buf.WriteByte(c.config.BirthDay)

	binary.Write(&buf, binary.BigEndian, c.config.Number)

	return buf.Bytes()
}

func (c *Client) sendMessage(msg []byte) {
	written := 0

	for written < len(msg) {
		n, err := c.conn.Write(msg[written:])
		if err != nil {
			log.Fatalf("Error writting to conn: %v", err)
		}
		written += n
	}
}

func recvACK(c *Client) error {
	ack, err := bufio.NewReader(c.conn).ReadByte()
	if err != nil {
		return err
	}

	if ack == 0 {
		return fmt.Errorf("Received NACK from server")
	}

	return nil
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, syscall.SIGTERM)

	// There is an autoincremental msgID to identify every message sent
	// Messages if the message amount threshold has not been surpassed
	for msgID := 1; msgID <= c.config.LoopAmount; msgID++ {
		// At this point every resource has been free. It is safe to exit if the signal has been received
		if listenForSigTerm(channel) {
			break
		}

		// Create the connection the server in every loop iteration. Send an
		if err := c.createClientSocket(); err != nil {
			continue
		}

		msg := createMessage(c)
		c.sendMessage(msg)

		err := recvACK(c)
		c.conn.Close()

		if err != nil {
			if err.Error() == "Received NACK from server" {
				log.Error("action: receive_message | result: nack")
			} else {
				log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
					c.config.ID,
					err,
				)
				return
			}
		}

		log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v",
			c.config.Document,
			c.config.Number,
		)

		// Wait a time between sending one message and the next one
		time.Sleep(c.config.LoopPeriod)

	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}

func listenForSigTerm(channel chan os.Signal) bool {
	select {
	case sig := <-channel:
		fmt.Println("Received signal", sig)
		return true
	default:
		return false
	}
}
