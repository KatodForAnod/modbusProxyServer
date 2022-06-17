package main

import (
	"github.com/goburrow/serial"
	"log"
	"modbusprottocol/client"
	"modbusprottocol/controller"
	"modbusprottocol/memory"
	"time"
)

func main() {
	log.SetFlags(log.Lshortfile)

	rtuClient := client.RTUClient{}
	iotControll := controller.IoTsController{}
	mem := memory.MemoryFmt{}
	iotControll.Init(mem)

	conf := serial.Config{
		Address:  "COM3",
		BaudRate: 115200,
		DataBits: 8,
		StopBits: 1,
		Parity:   "N",
		Timeout:  5 * time.Second,
		RS485:    serial.RS485Config{},
	}

	rtuClient.Init("testDevice")
	err := rtuClient.Connect(conf, 0x01)
	if err != nil {
		log.Fatal(err)
		return
	}
	iotControll.AddIoTsClients([]client.IoTClient{&rtuClient})

	iotControll.ObserveCoils("testDevice", 0x11, 3, time.Second)
	time.Sleep(time.Second * 10)
	rtuClient.StopObserveInform()
	time.Sleep(time.Second * 10)

	return
	/*fmt.Println(rtuClient.ReadCoils(0x11, 3))
	fmt.Println(rtuClient.WriteSingleCoil(0x11, client.ON))*/
}
