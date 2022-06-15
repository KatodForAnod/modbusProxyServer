package client

import (
	"github.com/goburrow/modbus"
	"github.com/goburrow/serial"
	"log"
)

type RTUClient struct {
	handler *modbus.RTUClientHandler
	client  modbus.Client
}

func (c *RTUClient) Connect(conf serial.Config, slaveId byte) error {
	log.Println("Connect with com port:", conf.Address)
	c.handler = modbus.NewRTUClientHandler(conf.Address)
	c.handler.SlaveId = slaveId

	c.handler.Config = conf
	if err := c.handler.Connect(); err != nil {
		log.Println(err)
		return err
	}

	c.client = modbus.NewClient(c.handler)

	return nil
}

func (c *RTUClient) Disconnect() error {
	log.Println("Disconnecting from port:", c.handler.Config.Address)
	if err := c.handler.Close(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// ReadCoils - чтение текущего состояния (ON/OFF) дискретных выходов. 1 бит. Диапазон 00001-10000
func (c *RTUClient) ReadCoils(address, quantity uint16) ([]byte, error) {
	log.Println("ReadCoils address:", address, "quantity:", quantity)

	result, err := c.client.ReadCoils(address, quantity)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	return result, nil
}

// WriteSingleCoil - изменения состояния дискретного выхода в ON или OFF. 1 бит. Диапазон 00001-10000
func (c *RTUClient) WriteSingleCoil(address uint16, coil CoilType) ([]byte, error) {
	log.Println("WriteSingleCoil address:", address, "coil:", coil.String())

	result, err := c.client.WriteSingleCoil(address, coil.ToUint16())
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	return result, nil
}

// WriteMultipleCoils - изменения состояния нескольких дискретных выходов в ON или OFF. 1 бит. Диапазон 00001-10000
func (c *RTUClient) WriteMultipleCoils(address, quantity uint16, coils []CoilType) ([]byte, error) {
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

// WriteSingleRegister - запись одного регистра. 16 бит. Диапазон 40001 - 50000
func (c *RTUClient) WriteSingleRegister(address, value uint16) ([]byte, error) {
	log.Println("WriteSingleRegister address:", address, "value:", value)

	result, err := c.client.WriteSingleRegister(address, value)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	return result, nil
}

// WriteMultipleRegisters - запись нескольких регистров. 16 бит. Диапазон 40001 - 50000
func (c *RTUClient) WriteMultipleRegisters(address, quantity uint16, values []uint16) ([]byte, error) {
	log.Println("WriteMultipleRegisters address:", address, "quantity:", quantity, "values", values)

	bytes, err := c.hexDecimalInBytes(values)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	result, err := c.client.WriteMultipleRegisters(address, quantity, bytes)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	return result, nil
}

// ReadDiscreteInputs - чтение текущего состояния (ON/OFF) дискретных выходов/входов. 1 бит. Диапазон 00001-20000
func (c *RTUClient) ReadDiscreteInputs(address, quantity uint16) ([]byte, error) {
	log.Println("ReadDiscreteInputs address:", address, "quantity:", quantity)

	result, err := c.client.ReadDiscreteInputs(address, quantity)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	return result, nil
}

// ReadInputRegisters - чтение входных регистров. 16 бит. Диапазон 30001-40000
func (c *RTUClient) ReadInputRegisters(address, quantity uint16) ([]byte, error) {
	log.Println("ReadInputRegisters address:", address, "quantity:", quantity)

	result, err := c.client.ReadInputRegisters(address, quantity)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	return result, nil
}

// ReadHoldingRegisters - чтение регистров хранения. 16 бит. Диапазон 40001-50000
func (c *RTUClient) ReadHoldingRegisters(address, quantity uint16) ([]byte, error) {
	log.Println("ReadHoldingRegisters address:", address, "quantity:", quantity)

	result, err := c.client.ReadHoldingRegisters(address, quantity)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	return result, nil
}

// ReadWriteMultipleRegisters - чтение/запись нескольких регистров. 16 бит. Диапазон 40001-50000
func (c *RTUClient) ReadWriteMultipleRegisters(readAddress uint16, readQuantity uint16,
	writeAddress uint16, writeQuantity uint16, values []uint16) ([]byte, error) {
	log.Println("ReadWriteMultipleRegisters readAddress:", readAddress,
		"readQuantity:", readQuantity, "writeAddress:", writeAddress,
		"values:", values)

	bytes, err := c.hexDecimalInBytes(values)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	result, err := c.client.ReadWriteMultipleRegisters(readAddress, readQuantity,
		writeAddress, writeQuantity, bytes)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	return result, nil
}

// ReadFIFOQueue - чтение содержимого очереди FIFO. 16 бит. Диапазон 40001-50000
func (c *RTUClient) ReadFIFOQueue(address uint16) ([]byte, error) {
	log.Println("ReadFIFOQueue address:", address)

	result, err := c.client.ReadFIFOQueue(address)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	return result, nil
}

// MaskWriteRegister - маскированная запись регистра. 16 бит. Диапазон 40001-50000
func (c *RTUClient) MaskWriteRegister(address uint16, andMask uint16, orMask uint16) ([]byte, error) {
	log.Println("MaskWriteRegister address:", address,
		"andMask:", andMask, "orMask:", orMask)

	result, err := c.client.MaskWriteRegister(address, andMask, orMask)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	return result, nil
}

func (c *RTUClient) hexDecimalInBytes(values []uint16) ([]byte, error) {
	return []byte{}, nil
}
