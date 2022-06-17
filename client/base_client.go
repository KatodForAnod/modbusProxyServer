package client

import (
	"errors"
	"github.com/goburrow/modbus"
	"log"
	"modbusprottocol/config"
	"time"
)

type IoTClient interface {
	GetDeviceName() string
	Connect() error
	Disconnect() error
	ReadCoils(address, quantity uint16) ([]byte, error)
	WriteSingleCoil(address uint16, coil CoilType) ([]byte, error)
	WriteMultipleCoils(address, quantity uint16, coils []CoilType) ([]byte, error)
	WriteSingleRegister(address, value uint16) ([]byte, error)
	WriteMultipleRegisters(address, quantity uint16, values []uint16) ([]byte, error)
	ReadDiscreteInputs(address, quantity uint16) ([]byte, error)
	ReadInputRegisters(address, quantity uint16) ([]byte, error)
	ReadHoldingRegisters(address, quantity uint16) ([]byte, error)
	ReadWriteMultipleRegisters(readAddress uint16, readQuantity uint16,
		writeAddress uint16, writeQuantity uint16, values []uint16) ([]byte, error)
	ReadFIFOQueue(address uint16) ([]byte, error)
	MaskWriteRegister(address uint16, andMask uint16, orMask uint16) ([]byte, error)

	StartObserveInform(save func() error, duration time.Duration) error
	IsObserveInformProcess() bool
	StopObserveInform() error
}

type BaseClient struct {
	deviceName             string
	clientType             config.ClientType
	client                 modbus.Client
	isObserveInformProcess *bool
	conf                   config.IotConfig
}

func (c *BaseClient) Init(conf config.IotConfig) {
	c.isObserveInformProcess = new(bool)
	c.deviceName = conf.DeviceName
	c.conf = conf
	c.clientType = conf.TypeClient
}

func (c *BaseClient) GetDeviceName() string {
	return c.deviceName
}

func (c *BaseClient) Connect() error {
	return errors.New("base client - override this method")
}

func (c *BaseClient) Disconnect() error {
	return errors.New("base client - override this method")
}

// ReadCoils - чтение текущего состояния (ON/OFF) дискретных выходов. 1 бит. Диапазон 00001-10000
func (c *BaseClient) ReadCoils(address, quantity uint16) ([]byte, error) {
	log.Println("ReadCoils address:", address, "quantity:", quantity)

	result, err := c.client.ReadCoils(address, quantity)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	return result, nil
}

// WriteSingleCoil - изменения состояния дискретного выхода в ON или OFF. 1 бит. Диапазон 00001-10000
func (c *BaseClient) WriteSingleCoil(address uint16, coil CoilType) ([]byte, error) {
	log.Println("WriteSingleCoil address:", address, "coil:", coil.String())

	result, err := c.client.WriteSingleCoil(address, coil.ToUint16())
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	return result, nil
}

// WriteMultipleCoils - изменения состояния нескольких дискретных выходов в ON или OFF. 1 бит. Диапазон 00001-10000
func (c *BaseClient) WriteMultipleCoils(address, quantity uint16, coils []CoilType) ([]byte, error) {
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
func (c *BaseClient) WriteSingleRegister(address, value uint16) ([]byte, error) {
	log.Println("WriteSingleRegister address:", address, "value:", value)

	result, err := c.client.WriteSingleRegister(address, value)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	return result, nil
}

// WriteMultipleRegisters - запись нескольких регистров. 16 бит. Диапазон 40001 - 50000
func (c *BaseClient) WriteMultipleRegisters(address, quantity uint16, values []uint16) ([]byte, error) {
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
func (c *BaseClient) ReadDiscreteInputs(address, quantity uint16) ([]byte, error) {
	log.Println("ReadDiscreteInputs address:", address, "quantity:", quantity)

	result, err := c.client.ReadDiscreteInputs(address, quantity)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	return result, nil
}

// ReadInputRegisters - чтение входных регистров. 16 бит. Диапазон 30001-40000
func (c *BaseClient) ReadInputRegisters(address, quantity uint16) ([]byte, error) {
	log.Println("ReadInputRegisters address:", address, "quantity:", quantity)

	result, err := c.client.ReadInputRegisters(address, quantity)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	return result, nil
}

// ReadHoldingRegisters - чтение регистров хранения. 16 бит. Диапазон 40001-50000
func (c *BaseClient) ReadHoldingRegisters(address, quantity uint16) ([]byte, error) {
	log.Println("ReadHoldingRegisters address:", address, "quantity:", quantity)

	result, err := c.client.ReadHoldingRegisters(address, quantity)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	return result, nil
}

// ReadWriteMultipleRegisters - чтение/запись нескольких регистров. 16 бит. Диапазон 40001-50000
func (c *BaseClient) ReadWriteMultipleRegisters(readAddress uint16, readQuantity uint16,
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
func (c *BaseClient) ReadFIFOQueue(address uint16) ([]byte, error) {
	log.Println("ReadFIFOQueue address:", address)

	result, err := c.client.ReadFIFOQueue(address)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	return result, nil
}

// MaskWriteRegister - маскированная запись регистра. 16 бит. Диапазон 40001-50000
func (c *BaseClient) MaskWriteRegister(address uint16, andMask uint16, orMask uint16) ([]byte, error) {
	log.Println("MaskWriteRegister address:", address,
		"andMask:", andMask, "orMask:", orMask)

	result, err := c.client.MaskWriteRegister(address, andMask, orMask)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	return result, nil
}

func (c *BaseClient) hexDecimalInBytes(values []uint16) ([]byte, error) {
	return []byte{}, nil
}

func (c *BaseClient) StartObserveInform(save func() error, duration time.Duration) error {
	log.Println("StartObserveInform device:", c.deviceName)
	if *c.isObserveInformProcess {
		err := errors.New("already observe")
		log.Println(err)
		return err
	}

	tr := true
	c.isObserveInformProcess = &tr

	for *c.isObserveInformProcess {
		if err := save(); err != nil {
			log.Println(err)
		}
		time.Sleep(duration)
	}
	return nil
}

func (c *BaseClient) IsObserveInformProcess() bool {
	return *c.isObserveInformProcess
}

func (c *BaseClient) StopObserveInform() error {
	log.Println("StopObserveInform device:", c.deviceName)
	fl := false
	c.isObserveInformProcess = &fl
	return nil
}
