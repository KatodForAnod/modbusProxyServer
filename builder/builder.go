package builder

import (
	"errors"
	"github.com/goburrow/serial"
	"log"
	"modbusprottocol/client"
	"modbusprottocol/config"
	"time"
)

type BuildClient struct {
}

func (c *BuildClient) BuildClient(iotConfig config.IotConfig) (client.IoTClient, error) {
	conf := serial.Config{
		Address:  iotConfig.ComPort,
		BaudRate: iotConfig.BaudRate,
		DataBits: iotConfig.DataBits,
		StopBits: iotConfig.StopBits,
		Parity:   iotConfig.Parity,
		Timeout:  time.Duration(iotConfig.TimeoutSeconds) * time.Second,
		RS485:    serial.RS485Config{},
	}

	if TCPClient == iotConfig.TypeClient.String() {
		tcpClient := client.TCPClient{}
		tcpClient.Init(iotConfig.DeviceName)
		err := tcpClient.Connect(conf, iotConfig.SlaveId)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return &tcpClient, nil
	} else if RTUClient == iotConfig.TypeClient.String() {
		rtuClient := client.RTUClient{}
		rtuClient.Init(iotConfig.DeviceName)
		err := rtuClient.Connect(conf, iotConfig.SlaveId)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return &rtuClient, nil
	}

	return nil, errors.New("type " +
		iotConfig.TypeClient.String() + " does not support")
}
