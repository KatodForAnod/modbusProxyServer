package main

import (
	"log"
	"modbusProxyServer/builder"
	"modbusProxyServer/client"
	"modbusProxyServer/config"
	"modbusProxyServer/controller"
	"modbusProxyServer/memory"
	"time"
)

func main() {
	log.SetFlags(log.Lshortfile)

	iotControll := controller.IoTsController{}
	mem := memory.MemoryFmt{}
	iotControll.Init(mem)

	conf, err := config.LoadConfig()
	if err != nil {
		log.Println(err)
		return
	}

	buildClient := builder.BuildClient{}
	var clients []client.IoTClient
	for _, device := range conf.IoTsDevices {
		ioTClient, err := buildClient.BuildClient(device)
		if err != nil {
			log.Println(err)
			return
		}
		clients = append(clients, ioTClient)
	}

	iotControll.AddIoTsClients(clients)
	//iotControll.ObserveCoils("testDevice", 0x11, 3, time.Second)
	time.Sleep(time.Second * 10)
	return
}
