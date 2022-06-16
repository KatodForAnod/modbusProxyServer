package client

import (
	"github.com/goburrow/modbus"
	"github.com/goburrow/serial"
	"log"
)

type TCPClient struct {
	handler *modbus.TCPClientHandler
	BaseClient
}

func (c *TCPClient) Connect(conf serial.Config, slaveId byte) error {
	log.Println("Connect TCPClient with com port:", conf.Address)
	handler := modbus.NewTCPClientHandler(conf.Address)
	handler.SlaveId = slaveId

	handler.Address = conf.Address
	if err := handler.Connect(); err != nil {
		log.Println(err)
		return err
	}

	c.handler = handler
	c.client = modbus.NewClient(c.handler)

	return nil
}

func (c *TCPClient) Disconnect() error {
	log.Println("Disconnecting TCPClient from port:", c.handler.Address)
	if err := c.handler.Close(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
