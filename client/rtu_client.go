package client

import (
	"github.com/goburrow/modbus"
	"github.com/goburrow/serial"
	"log"
)

type RTUClient struct {
	handler *modbus.RTUClientHandler
	client modbus.Client
}

func (c *RTUClient) Connect(conf serial.Config, comPort string) error {
	log.Println("Connect with com port:", comPort)
	c.handler = modbus.NewRTUClientHandler("COM3")

	c.handler.Config = conf
	if err := c.handler.Connect(); err != nil {
		log.Println(err)
		return err
	}

	c.client = modbus.NewClient(c.handler)

	return nil
}

func (c *RTUClient) ReadCoils(address, quantity uint16) ([]byte, error) {
	log.Println("ReadCoils address:", address, "quantity:", quantity)

	result, err := c.client.ReadCoils(address, quantity)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	return result, nil
}

func (c *RTUClient) WriteSingleCoil (address, value uint16) ([]byte, error)  {
	log.Println("WriteSingleCoil address:", address, "value:", value)

	result, err := c.client.WriteSingleCoil(address, value)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	return result, nil
}

func (c *RTUClient) WriteMultipleCoils (address, quantity uint16, coils []CoilType) ([]byte, error)  {
	log.Println("WriteMultipleCoils address:", address, "quantity:", quantity, "coils:", coils)

	var outCoils []byte
	for _, coil := range coils {
		outCoils = append(outCoils, coil.ToByteSlice()...)
	}

	result, err := c.client.WriteMultipleCoils(address, quantity, outCoils)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	return result, nil
}



