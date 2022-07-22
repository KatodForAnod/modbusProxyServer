package client

import (
	"fmt"
	"github.com/goburrow/modbus"
	log "github.com/sirupsen/logrus"
)

type TCPClient struct {
	handler *modbus.TCPClientHandler
	BaseClient
}

func (c *TCPClient) Connect() error {
	log.Println("Connect TCPClient with com port:", c.conf.ComPort)
	handler := modbus.NewTCPClientHandler(c.conf.ComPort)
	handler.SlaveId = c.conf.SlaveId

	handler.Address = c.conf.ComPort
	if err := handler.Connect(); err != nil {
		return fmt.Errorf("tcpclient connect: %s", err.Error())
	}

	c.handler = handler
	c.client = modbus.NewClient(c.handler)

	return nil
}

func (c *TCPClient) Disconnect() error {
	log.Println("Disconnecting TCPClient from port:", c.handler.Address)
	if err := c.handler.Close(); err != nil {
		return fmt.Errorf("tcpclient disconnect: %s", err.Error())
	}

	return nil
}
