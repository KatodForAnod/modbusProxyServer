package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	ProxyServerAddr string      `json:"proxy_server_addr"`
	IoTsDevices     []IotConfig `json:"iots_devices"`
}

type IotConfig struct {
	DeviceName     string     `json:"device_name"`
	TypeClient     ClientType `json:"type_client"`
	SlaveId        byte       `json:"slave_id"`
	ComPort        string     `json:"com_port"`
	BaudRate       int        `json:"baud_rate"`
	DataBits       int        `json:"data_bits"`
	StopBits       int        `json:"stop_bits"`
	Parity         string     `json:"parity"`
	TimeoutSeconds int        `json:"timeout_seconds"`
}

var configPath = "conf.config"

func LoadConfig() (loadedConf Config, err error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Println(err)
		return Config{}, err
	}

	err = json.Unmarshal(data, &loadedConf)
	return
}
