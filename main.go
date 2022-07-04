package main

import (
	"log"
	"modbusProxyServer/config"
	"modbusProxyServer/controller"
	"modbusProxyServer/memory"
	"modbusProxyServer/server"
)

func main() {
	log.SetFlags(log.Lshortfile)
	/*err := logsetting.LogInit()
	if err != nil {
		log.Fatal(err)
	}*/

	iotControll := controller.IoTsController{}
	mem := memory.MemoryFmt{}
	iotControll.Init(mem)

	conf, err := config.LoadConfig()
	if err != nil {
		log.Println(err)
		return
	}

	serv_controller := controller.Controller{}
	serv_controller.Init(mem, iotControll)

	serv := server.Server{}
	serv.StartServer(conf, &serv_controller)

	return
}
