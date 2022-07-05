package builder

import (
	"modbusProxyServer/config"
	"testing"
)

func TestBuildClientTCPClient(t *testing.T) {
	builder := BuildClient{}
	conf := config.IotConfig{
		DeviceName:     "qwerty",
		TypeClient:     config.ClientType{Cl: config.TCPClientType},
		SlaveId:        0,
		ComPort:        "",
		BaudRate:       0,
		DataBits:       0,
		StopBits:       0,
		Parity:         "",
		TimeoutSeconds: 0,
	}

	_, err := builder.BuildClient(conf)
	if err != nil {
		t.Errorf("function BuildClient() is corrupted: unexpected error: %s", err)
		return
	}
}

func TestBuildClientNotExistClientType(t *testing.T) {
	builder := BuildClient{}
	conf := config.IotConfig{
		DeviceName:     "qwerty",
		TypeClient:     config.ClientType{Cl: "notExistType"},
		SlaveId:        0,
		ComPort:        "",
		BaudRate:       0,
		DataBits:       0,
		StopBits:       0,
		Parity:         "",
		TimeoutSeconds: 0,
	}

	_, err := builder.BuildClient(conf)
	if err == nil {
		t.Error("BuildClient() should return error in that case")
	}
}

func TestBuildClientRTUClient(t *testing.T) {
	builder := BuildClient{}
	conf := config.IotConfig{
		DeviceName:     "qwerty",
		TypeClient:     config.ClientType{Cl: config.RTUClientType},
		SlaveId:        0,
		ComPort:        "",
		BaudRate:       0,
		DataBits:       0,
		StopBits:       0,
		Parity:         "",
		TimeoutSeconds: 0,
	}

	_, err := builder.BuildClient(conf)
	if err != nil {
		t.Errorf("function BuildClient() is corrupted: unexpected error: %s", err)
		return
	}
}
