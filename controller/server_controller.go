package controller

import (
	"fmt"
	"modbusProxyServer/builder"
	"modbusProxyServer/client"
	"modbusProxyServer/config"
	"modbusProxyServer/logsetting"
	"modbusProxyServer/memory"
	"strconv"
	"time"
)

type ServerController interface {
	GetLastNRowsLogs(nRows int) ([]string, error)
	GetInformation(deviceName string) ([]byte, error)
	AddIoTDevice(device config.IotConfig) error
	RmIoTDevice(deviceName string) error
	StopObserveDevice(deviceName string) error
	ObserveIoTCoils(deviceName, address, quantity, timeSecondsDuration string) error
}

type Controller struct {
	mem            memory.Memory
	ioTsController IoTsController
}

func (c *Controller) Init(mem memory.Memory, controller IoTsController) {
	c.mem = mem
	c.ioTsController = controller
}

func (c *Controller) GetInformation(deviceName string) ([]byte, error) {
	load, err := c.mem.Load(deviceName)
	if err != nil {
		return []byte{}, err
	}

	return load, nil
}

func (c *Controller) GetLastNRowsLogs(nRows int) ([]string, error) {
	file, err := logsetting.OpenLastLogFile()
	if err != nil {
		return []string{}, err
	}

	logs, err := logsetting.GetNLastLines(file, nRows)
	if err != nil {
		return []string{}, err
	}

	return logs, nil
}

func (c *Controller) AddIoTDevice(device config.IotConfig) error {
	buildClient := builder.BuildClient{}
	iotClient, err := buildClient.BuildClient(device)
	if err != nil {
		return err
	}

	err = c.ioTsController.AddIoTsClients([]client.IoTClient{iotClient})
	return err
}

func (c *Controller) RmIoTDevice(deviceName string) error {
	return c.ioTsController.RemoveIoTsClients([]string{deviceName})
}

func (c *Controller) StopObserveDevice(deviceName string) error {
	return c.ioTsController.StopObserveIoTDevice(deviceName)
}

func (c *Controller) ObserveIoTCoils(deviceName, address, quantity, timeSecondsDuration string) error {
	quantityUint, err := strconv.ParseUint(quantity, 10, 16)
	if err != nil {
		return fmt.Errorf("parse quantity: %s", err.Error())
	}
	addressUint, err := strconv.ParseUint(address, 10, 16)
	if err != nil {
		return fmt.Errorf("parse address: %s", err.Error())
	}
	timeInt, err := strconv.ParseInt(timeSecondsDuration, 10, 64)
	if err != nil {
		return fmt.Errorf("parse timeSecondsDuration: %s", err.Error())
	}

	err = c.ioTsController.ObserveCoils(deviceName,
		uint16(addressUint), uint16(quantityUint), time.Duration(timeInt)*time.Second)
	return err
}
