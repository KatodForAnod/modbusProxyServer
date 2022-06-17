package controller

import (
	"log"
	"modbusprottocol/client"
	"modbusprottocol/logsetting"
	"modbusprottocol/memory"
)

type Controller struct {
	mem            memory.Memory
	ioTsController IoTsController
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

func (c *Controller) AddIoTDevice(device client.IoTClient) error {
	log.Println("controller AddIoTDevice")
	err := c.ioTsController.AddIoTsClients([]client.IoTClient{device})
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
