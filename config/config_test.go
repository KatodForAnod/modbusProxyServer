package config

import (
	"os"
	"testing"
)

const (
	confBodyTCP = `{
    "proxy_server_addr":"127.0.0.1:5300",
    "iots_devices":[
        {
			"type_client":"tcp"
		}
    ]}`

	confBodyRTU = `{
    "proxy_server_addr":"127.0.0.1:5300",
    "iots_devices":[
        {
			"type_client":"rtu"
		}
    ]}`

	confBodyWrongType = `{
    "proxy_server_addr":"127.0.0.1:5300",
    "iots_devices":[
        {
			"type_client":"wrongtype"
		}
    ]}`
)

func createConfig(t *testing.T, confBody string) (filePath string, err error) {
	file, err := os.CreateTemp("", configPath)
	if err != nil {
		t.Errorf("cant create temp conf file: unexpected error: %s", err)
		return "", err
	}
	defer file.Close()

	_, err = file.WriteString(confBody)
	if err != nil {
		t.Errorf("cant write to temp conf file: unexpected error: %s", err)
		return "", err
	}

	return file.Name(), nil
}

func deleteConfig(filePath string, t *testing.T) error {
	err := os.Remove(filePath)
	if err != nil {
		t.Errorf("cant delete temp conf file: unexpected error: %s", err)
		return err
	}
	return nil
}

func TestLoadConfigSuccessTCPClient(t *testing.T) {
	t.Log("testing loading config")
	filePath, err := createConfig(t, confBodyTCP)
	if err != nil {
		t.Errorf("function createConfig() is corrupted: unexpected error: %s", err)
		return
	}

	temp := configPath
	configPath = filePath

	defer func() {
		configPath = temp
		err = deleteConfig(filePath, t)
		if err != nil {
			t.Log(err)
			return
		}
	}()

	_, err = LoadConfig()
	if err != nil {
		t.Errorf("function LoadConfig() is corrupted: unexpected error: %s", err)
		return
	}
}

func TestLoadConfigSuccessRTUClient(t *testing.T) {
	t.Log("testing loading config")
	filePath, err := createConfig(t, confBodyRTU)
	if err != nil {
		t.Errorf("function createConfig() is corrupted: unexpected error: %s", err)
		return
	}

	temp := configPath
	configPath = filePath

	defer func() {
		configPath = temp
		err = deleteConfig(filePath, t)
		if err != nil {
			t.Log(err)
			return
		}
	}()

	_, err = LoadConfig()
	if err != nil {
		t.Errorf("function LoadConfig() is corrupted: unexpected error: %s", err)
		return
	}
}

func TestLoadConfigCheckField(t *testing.T) {
	t.Log("testing loading config")
	filePath, err := createConfig(t, confBodyRTU)
	if err != nil {
		t.Errorf("function createConfig() is corrupted: unexpected error: %s", err)
		return
	}

	temp := configPath
	configPath = filePath

	defer func() {
		configPath = temp
		err = deleteConfig(filePath, t)
		if err != nil {
			t.Log(err)
			return
		}
	}()

	conf, err := LoadConfig()
	if err != nil {
		t.Errorf("function LoadConfig() is corrupted: unexpected error: %s", err)
		return
	}

	expectedLen := 1
	if len(conf.IoTsDevices) < expectedLen {
		t.Errorf("expected len > %d, instead got: %d", expectedLen, len(conf.IoTsDevices))
	}

	if conf.IoTsDevices[0].TypeClient.String() != RTUClientType {
		t.Errorf("expected client type = %s, instead got: %s",
			RTUClientType, conf.IoTsDevices[0].TypeClient.String())
	}
}

func TestLoadConfigFail(t *testing.T) {
	t.Log("testing unsuccessful loading config")
	filePath, err := createConfig(t, confBodyWrongType)
	if err != nil {
		t.Errorf("function createConfig() is corrupted: unexpected error: %s", err)
		return
	}

	temp := configPath
	configPath = filePath

	defer func() {
		configPath = temp
		err = deleteConfig(filePath, t)
		if err != nil {
			t.Log(err)
			return
		}
	}()

	_, err = LoadConfig()
	if err == nil {
		t.Error("LoadConfig() should return error in that case")
		return
	}
}
