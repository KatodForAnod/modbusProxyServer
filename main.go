package main

import (
	"github.com/KatodForAnod/modbusProxyServer/config"
	"github.com/KatodForAnod/modbusProxyServer/controller"
	"github.com/KatodForAnod/modbusProxyServer/memory"
	"github.com/KatodForAnod/modbusProxyServer/server"
	log "github.com/sirupsen/logrus"
)

func main() {
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
