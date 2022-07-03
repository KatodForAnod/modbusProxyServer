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

func createConfig(t *testing.T, confBody string) error {
	file, err := os.CreateTemp("", configPath)
	if err != nil {
		t.Error("cant create temp conf file")
		return err
	}

	configPath = file.Name()

	_, err = file.WriteString(confBody)
	if err != nil {
		t.Error("cant write to temp conf file")
		return err
	}

	return nil
}

func deleteConfig(t *testing.T) error {
	err := os.Remove(configPath)
	if err != nil {
		t.Error("cant delete temp conf file")
		return err
	}
	return nil
}

func TestLoadConfig_Success(t *testing.T) {
	err := createConfig(t, confBody)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = LoadConfig()
	if err != nil {
		t.Error(err)
		return
	}

	err = deleteConfig(t)
	if err != nil {
		t.Log(err)
		return
	}
}

func TestLoadConfig_Fail(t *testing.T) {
	err := createConfig(t, confBody2)
	if err != nil {
		t.Error(err)
		return
	}

	defer func() {
		err = deleteConfig(t)
		if err != nil {
			t.Log(err)
		}
	}()

	_, err = LoadConfig()
	if err == nil {
		t.FailNow()
		return
	}
}
