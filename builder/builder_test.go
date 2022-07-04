package builder

import (
	"modbusProxyServer/config"
	"testing"
)

func TestBuildClient_BuildClient(t *testing.T) {
	t.Log("Checking tcp client building")
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
		t.Error(err)
		return
	}
}

func TestBuildClient_BuildClient2(t *testing.T) {
	t.Log("Checking not exist client type")
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
		t.FailNow()
	}
}

func TestBuildClient_BuildClient3(t *testing.T) {
	t.Log("Checking rtu client building")
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
		t.Error(err)
		return
	}
}
