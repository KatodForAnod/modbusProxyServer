package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"modbusprottocol/builder"
)

type Config struct {
	ProxyServerAddr string      `json:"proxy_server_addr"`
	IoTsDevices     []IotConfig `json:"iots_devices"`
}

type IotConfig struct {
	TypeClient        builder.ClientType `json:"type_client"`
	ComPort           string             `json:"com_port"`
	BaudRate          int                `json:"baud_rate"`
	DataBits          int                `json:"data_bits"`
	StopBits          int                `json:"stop_bits"`
	Parity            string             `json:"parity"`
	TimeoutSeconds    int                `json:"timeout_seconds"`
	TypeIotConnection string             `json:"type_iot_connection"`
}

const configPath = "conf.config"

func LoadConfig() (loadedConf Config, err error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Println(err)
		return Config{}, err
	}

	err = json.Unmarshal(data, &loadedConf)
	return
}
