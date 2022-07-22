package controller

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"modbusProxyServer/client"
	"modbusProxyServer/memory"
	"time"
)

type IoTsController struct {
	ioTDevices map[string]client.IoTClient
	mem        memory.Memory
}

func (c *IoTsController) Init(mem memory.Memory) {
	c.ioTDevices = make(map[string]client.IoTClient)
	c.mem = mem
}

func (c *IoTsController) GetIoTsClients(deviceName []string) (founded []client.IoTClient, err error) {
	for _, s := range deviceName {
		if iotDevice, isExist := c.ioTDevices[s]; isExist {
			founded = append(founded, iotDevice)
		} else {
			return []client.IoTClient{}, fmt.Errorf("iotscontroller GetIoTsClients: device %s not found", s)
		}
	}

	return
}

func (c *IoTsController) AddIoTsClients(devices []client.IoTClient) error {
	for _, device := range devices {
		if _, isExist := c.ioTDevices[device.GetDeviceName()]; isExist {
			return fmt.Errorf("iotscontroller AddIoTsClients:"+
				" device %s already exist", device.GetDeviceName())
		}
	}

	for _, device := range devices {
		err := device.Connect()
		if err != nil {
			log.Errorln(err)
			continue
		}
		c.ioTDevices[device.GetDeviceName()] = device
	}
	return nil
}

func (c *IoTsController) RemoveIoTsClients(devicesName []string) error {
	var founded []client.IoTClient
	for _, deviceName := range devicesName {
		if iot, isExist := c.ioTDevices[deviceName]; !isExist {
			return fmt.Errorf("iotscontroller RemoveIoTsClients: device %s not exist", deviceName)
		} else {
			founded = append(founded, iot)
		}
	}

	for _, tClient := range founded {
		if tClient.IsObserveInformProcess() {
			if err := tClient.StopObserveInform(); err != nil {
				log.Errorln(err)
			}
		}
		if err := tClient.Disconnect(); err != nil {
			log.Errorln(err)
		}
		delete(c.ioTDevices, tClient.GetDeviceName())
	}

	return nil
}

func (c *IoTsController) ObserveFIFOQueue(deviceName string, address uint16) error {
	return errors.New("func not work")
}

func (c *IoTsController) ObserveHoldingRegisters(deviceName string, address, quantity uint16) error {
	return errors.New("func not work")
}

func (c *IoTsController) ObserveInputRegisters(deviceName string, address, quantity uint16) error {
	return errors.New("func not work")
}

func (c *IoTsController) ObserveDiscreteInputs(deviceName string, address, quantity uint16) error {
	return errors.New("func not work")
}

func (c *IoTsController) ObserveCoils(deviceName string, address, quantity uint16, d time.Duration) error {
	iot, isExist := c.ioTDevices[deviceName]
	if !isExist {
		return fmt.Errorf("iotscontroller ObserveCoils: device %s not exist", deviceName)
	}

	if iot.IsObserveInformProcess() {
		log.Println("device:", iot.GetDeviceName(), "already in observe process")
		return nil
	}

	delay := time.Millisecond * 500
	timer := time.AfterFunc(d+delay, func() {
		if iot.IsObserveInformProcess() {
			log.Println("iot device -", iot.GetDeviceName(), "timeout response")
		}
	})

	saveFunc := func() error {
		coils, err := iot.ReadCoils(address, quantity)
		if err != nil {
			return fmt.Errorf("savefunc readCoils err: %s", err.Error())
		}
		//  memory.MsgType{}??
		if err := c.mem.Save(coils, memory.MsgType{}, iot.GetDeviceName()); err != nil {
			return fmt.Errorf("savefunc memorySave err: %s", err.Error())
		}

		timer.Reset(d + delay)
		return nil
	}

	go func() {
		err := iot.StartObserveInform(saveFunc, d)
		if err != nil {
			log.Errorln(err)
		}
	}()
	return nil
}

func (c *IoTsController) StopObserveIoTDevice(deviceName string) error {
	iot, isExist := c.ioTDevices[deviceName]
	if !isExist {
		return fmt.Errorf("iotscontroller StopObserveIoTDevice:"+
			" device %s not exist", deviceName)
	}

	if err := iot.StopObserveInform(); err != nil {
		return fmt.Errorf("iotscontroller StopObserveIoTDevice: %s", err.Error())
	}

	return nil
}
