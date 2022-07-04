package controller

import (
	"log"
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
	log.Println("controller get information of iot device", deviceName)

	load, err := c.mem.Load(deviceName)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	return load, nil
}

func (c *Controller) GetLastNRowsLogs(nRows int) ([]string, error) {
	log.Println("controller get lastNRowsLogs")
	file, err := logsetting.OpenLastLogFile()
	if err != nil {
		log.Println(err)
		return []string{}, err
	}

	logs, err := logsetting.GetNLastLines(file, nRows)
	if err != nil {
		log.Println(err)
		return []string{}, err
	}

	return logs, nil
}

func (c *Controller) AddIoTDevice(device config.IotConfig) error {
	log.Println("controller AddIoTDevice")

	buildClient := builder.BuildClient{}
	iotClient, err := buildClient.BuildClient(device)
	if err != nil {
		log.Println(err)
		return err
	}

	err = c.ioTsController.AddIoTsClients([]client.IoTClient{iotClient})
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (c *Controller) RmIoTDevice(deviceName string) error {
	log.Println("controller RmIoTDevice")

	err := c.ioTsController.RemoveIoTsClients([]string{deviceName})
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (c *Controller) StopObserveDevice(deviceName string) error {
	log.Println("controller stop observe device")

	if err := c.ioTsController.StopObserveIoTDevice(deviceName); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (c *Controller) ObserveIoTCoils(deviceName, address, quantity, timeSecondsDuration string) error {
	log.Println("controller ObserveIoTCoils deviceName:", deviceName)
	quantityUint, err := strconv.ParseUint(quantity, 10, 16)
	if err != nil {
		log.Println(err)
		return err
	}
	addressUint, err := strconv.ParseUint(address, 10, 16)
	if err != nil {
		log.Println(err)
		return err
	}
	timeInt, err := strconv.ParseInt(timeSecondsDuration, 10, 64)
	if err != nil {
		log.Println(err)
		return err
	}

	if err := c.ioTsController.ObserveCoils(deviceName,
		uint16(addressUint), uint16(quantityUint), time.Duration(timeInt)*time.Second); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
