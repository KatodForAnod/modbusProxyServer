package client

import (
	"fmt"
	"github.com/goburrow/modbus"
	"github.com/goburrow/serial"
	log "github.com/sirupsen/logrus"
	"time"
)

type RTUClient struct {
	handler *modbus.RTUClientHandler
	BaseClient
}

func (c *RTUClient) Connect() error {
	log.Println("Connect RTUClient with com port:", c.conf.ComPort)
	handler := modbus.NewRTUClientHandler(c.conf.ComPort)
	handler.SlaveId = c.conf.SlaveId

	handler.Config = serial.Config{
		Address:  c.conf.ComPort,
		BaudRate: c.conf.BaudRate,
		DataBits: c.conf.DataBits,
		StopBits: c.conf.StopBits,
		Parity:   c.conf.Parity,
		Timeout:  time.Duration(c.conf.TimeoutSeconds) * time.Second,
		RS485:    serial.RS485Config{},
	}
	if err := handler.Connect(); err != nil {
		return fmt.Errorf("rtuclient connect: %s", err.Error())
	}

	c.handler = handler
	c.client = modbus.NewClient(c.handler)

	return nil
}

func (c *RTUClient) Disconnect() error {
	log.Println("Disconnecting RTUClient from port:", c.handler.Config.Address)
	if err := c.handler.Close(); err != nil {
		return fmt.Errorf("rtuclient dosconnect: %s", err.Error())
	}

	return nil
}
