package common

import (
	"bufio"
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
	// TODO: Chequear que el nombre y apellido sean menor a 255 de largo
	Name       string
	LastName   string
	Document   uint32
	BirthYear  uint16
	BirthMonth uint8
	BirthDay   uint8
	Number     uint32
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
	}
	c.conn = conn
	return nil
}

func createMessage(c *Client) []byte {
	var msg []byte
	bufU32 := make([]byte, 4)
	bufU16 := make([]byte, 2)

	msg = append(msg, byte(c.config.ID))

	binary.BigEndian.PutUint32(bufU32, uint32(len(c.config.Name)))
	msg = append(msg, bufU32...)
	msg = append(msg, []byte(c.config.Name)...)

	binary.BigEndian.PutUint32(bufU32, uint32(len(c.config.LastName)))
	msg = append(msg, bufU32...)
	msg = append(msg, []byte(c.config.LastName)...)

	binary.BigEndian.PutUint32(bufU32, c.config.Document)
	msg = append(msg, bufU32...)

	binary.BigEndian.PutUint16(bufU16, c.config.BirthYear)
	msg = append(msg, bufU16...)
	msg = append(msg, byte(c.config.BirthMonth))
	msg = append(msg, byte(c.config.BirthDay))

	binary.BigEndian.PutUint32(bufU32, c.config.Number)
	msg = append(msg, bufU32...)

	return msg
}

func (c *Client) sendMessage(msg []byte) {
	bytes := []byte(msg)
	written := 0

	for written < len(bytes) {
		n, err := c.conn.Write(bytes[written:])
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
	term := false

	// There is an autoincremental msgID to identify every message sent
	// Messages if the message amount threshold has not been surpassed
	for msgID := 1; msgID <= c.config.LoopAmount && !term; msgID++ {
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

		log.Infof("action: receive_message | result: success | dni: %v | numero: %v",
			c.config.Document,
			c.config.Number,
		)

		// Wait a time between sending one message and the next one
		time.Sleep(c.config.LoopPeriod)

		// At this point every resource has been free. It is safe to exit if the signal has been received
		listenForSigTerm(channel, &term)
	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}

func listenForSigTerm(channel chan os.Signal, term *bool) {
	select {
	case sig := <-channel:
		fmt.Println("Received signal", sig)
		*term = true
	default:
	}
}
