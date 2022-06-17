package controller

import (
	"errors"
	"log"
	"modbusprottocol/client"
	"modbusprottocol/memory"
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
	log.Println("GetIoTsClients")
	for _, s := range deviceName {
		if iotDevice, isExist := c.ioTDevices[s]; isExist {
			founded = append(founded, iotDevice)
		} else {
			err = errors.New("not found")
			log.Println(err)
			return []client.IoTClient{}, err
		}
	}

	return
}

func (c *IoTsController) AddIoTsClients(devices []client.IoTClient) error {
	log.Println("AddIoTsClients")
	for _, device := range devices {
		if _, isExist := c.ioTDevices[device.GetDeviceName()]; isExist {
			err := errors.New("device " + device.GetDeviceName() + " already exist")
			log.Println(err)
			return err
		}
	}

	for _, device := range devices {
		err := device.Connect()
		if err != nil {
			log.Println(err)
			continue
		}
		c.ioTDevices[device.GetDeviceName()] = device
	}
	return nil
}

func (c IoTsController) RemoveIoTsClients(devicesName []string) error {
	log.Println("RemoveIoTsClients")
	var founded []client.IoTClient
	for _, deviceName := range devicesName {
		if iot, isExist := c.ioTDevices[deviceName]; !isExist {
			err := errors.New("device " + deviceName + " not exist")
			log.Println(err)
			return err
		} else {
			founded = append(founded, iot)
		}
	}

	for _, tClient := range founded {
		if tClient.IsObserveInformProcess() {
			if err := tClient.StopObserveInform(); err != nil {
				log.Println(err)
			}
		}
		if err := tClient.Disconnect(); err != nil {
			log.Println(err)
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
	log.Println("ObserveCoils device:", deviceName)
	iot, isExist := c.ioTDevices[deviceName]
	if !isExist {
		err := errors.New("device not exist")
		log.Println(err)

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
			log.Println(err)
			return err
		}
		//  memory.MsgType{}??
		if err := c.mem.Save(coils, memory.MsgType{}, iot.GetDeviceName()); err != nil {
			log.Println(err)
			return err
		}

		timer.Reset(d + delay)
		return nil
	}

	go func() {
		err := iot.StartObserveInform(saveFunc, d)
		if err != nil {
			log.Println(err)
		}
	}()
	return nil
}
