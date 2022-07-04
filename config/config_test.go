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
		t.Error("cant create temp conf file")
		return "", err
	}
	defer file.Close()

	_, err = file.WriteString(confBody)
	if err != nil {
		t.Error("cant write to temp conf file")
		return "", err
	}

	return file.Name(), nil
}

func deleteConfig(filePath string, t *testing.T) error {
	err := os.Remove(filePath)
	if err != nil {
		t.Error("cant delete temp conf file")
		return err
	}
	return nil
}

func TestLoadConfig_Success(t *testing.T) {
	filePath, err := createConfig(t, confBody)
	if err != nil {
		t.Error(err)
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
		t.Error(err)
		return
	}
}

func TestLoadConfig_Fail(t *testing.T) {
	filePath, err := createConfig(t, confBody2)
	if err != nil {
		t.Error(err)
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
		t.FailNow()
		return
	}
}
