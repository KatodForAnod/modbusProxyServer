package config

import (
	"os"
	"testing"
)

const confBody = `{
    "proxy_server_addr":"127.0.0.1:5300",
    "iots_devices":[
        {
			"type_client":"tcp"
		}
    ]
}`

const confBody2 = `{
    "proxy_server_addr":"127.0.0.1:5300",
    "iots_devices":[
        {
			"type_client":"wrongtype"
		}
    ]
}`

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

func TestLoadConfig_Success(t *testing.T) {
	t.Log("testing loading config")
	filePath, err := createConfig(t, confBody)
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

func TestLoadConfig_Fail(t *testing.T) {
	t.Log("testing unsuccessful loading config")
	filePath, err := createConfig(t, confBody2)
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
