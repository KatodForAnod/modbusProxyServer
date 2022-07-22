package controller

import (
	log "github.com/sirupsen/logrus"
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
		log.Errorln(err)
		return []byte{}, err
	}

	return load, nil
}

func (c *Controller) GetLastNRowsLogs(nRows int) ([]string, error) {
	file, err := logsetting.OpenLastLogFile()
	if err != nil {
		log.Errorln(err)
		return []string{}, err
	}

	logs, err := logsetting.GetNLastLines(file, nRows)
	if err != nil {
		log.Errorln(err)
		return []string{}, err
	}

	return logs, nil
}

func (c *Controller) AddIoTDevice(device config.IotConfig) error {
	buildClient := builder.BuildClient{}
	iotClient, err := buildClient.BuildClient(device)
	if err != nil {
		log.Errorln(err)
		return err
	}

	err = c.ioTsController.AddIoTsClients([]client.IoTClient{iotClient})
	if err != nil {
		log.Errorln(err)
		return err
	}

	return nil
}

func (c *Controller) RmIoTDevice(deviceName string) error {
	err := c.ioTsController.RemoveIoTsClients([]string{deviceName})
	if err != nil {
		log.Errorln(err)
		return err
	}

	return nil
}

func (c *Controller) StopObserveDevice(deviceName string) error {
	if err := c.ioTsController.StopObserveIoTDevice(deviceName); err != nil {
		log.Errorln(err)
		return err
	}

	return nil
}

func (c *Controller) ObserveIoTCoils(deviceName, address, quantity, timeSecondsDuration string) error {
	quantityUint, err := strconv.ParseUint(quantity, 10, 16)
	if err != nil {
		log.Errorln(err)
		return err
	}
	addressUint, err := strconv.ParseUint(address, 10, 16)
	if err != nil {
		log.Errorln(err)
		return err
	}
	timeInt, err := strconv.ParseInt(timeSecondsDuration, 10, 64)
	if err != nil {
		log.Errorln(err)
		return err
	}

	if err := c.ioTsController.ObserveCoils(deviceName,
		uint16(addressUint), uint16(quantityUint), time.Duration(timeInt)*time.Second); err != nil {
		log.Errorln(err)
		return err
	}

	return nil
}
