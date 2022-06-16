package client

import (
	"github.com/goburrow/modbus"
	"github.com/goburrow/serial"
	"log"
)

type RTUClient struct {
	handler *modbus.RTUClientHandler
	BaseClient
}

func (c *RTUClient) Connect(conf serial.Config, slaveId byte) error {
	log.Println("Connect RTUClient with com port:", conf.Address)
	handler := modbus.NewRTUClientHandler(conf.Address)
	handler.SlaveId = slaveId

	handler.Config = conf
	if err := handler.Connect(); err != nil {
		log.Println(err)
		return err
	}

	c.handler = handler
	c.client = modbus.NewClient(c.handler)

	return nil
}

func (c *RTUClient) Disconnect() error {
	log.Println("Disconnecting RTUClient from port:", c.handler.Config.Address)
	if err := c.handler.Close(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
