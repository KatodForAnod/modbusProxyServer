package builder

import (
	"errors"
	"modbusprottocol/client"
	"modbusprottocol/config"
)

type BuildClient struct {
}

func (c *BuildClient) BuildClient(iotConfig config.IotConfig) (client.IoTClient, error) {
	if client.TCPClientType == iotConfig.TypeClient.String() {
		tcpClient := client.TCPClient{}
		tcpClient.Init(iotConfig)
		return &tcpClient, nil
	} else if client.RTUClientType == iotConfig.TypeClient.String() {
		rtuClient := client.RTUClient{}
		rtuClient.Init(iotConfig)
		return &rtuClient, nil
	}

	return nil, errors.New("type " +
		iotConfig.TypeClient.String() + " does not support")
}
