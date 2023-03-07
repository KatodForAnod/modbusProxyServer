package builder

import (
	"errors"
	"github.com/KatodForAnod/modbusProxyServer/client"
	"github.com/KatodForAnod/modbusProxyServer/config"
)

type BuildClient struct {
}

func (c *BuildClient) BuildClient(iotConfig config.IotConfig) (client.IoTClient, error) {
	if config.TCPClientType == iotConfig.TypeClient.String() {
		tcpClient := client.TCPClient{}
		tcpClient.Init(iotConfig)
		return &tcpClient, nil
	} else if config.RTUClientType == iotConfig.TypeClient.String() {
		rtuClient := client.RTUClient{}
		rtuClient.Init(iotConfig)
		return &rtuClient, nil
	}

	return nil, errors.New("type " +
		iotConfig.TypeClient.String() + " does not support")
}
