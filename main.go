package main

import (
	"fmt"
	"github.com/goburrow/serial"
	"log"
	"modbusprottocol/client"
	"time"
)

func main() {
	log.SetFlags(log.Lshortfile)

	rtuClient := client.RTUClient{}
	conf := serial.Config{
		Address:  "COM3",
		BaudRate: 115200,
		DataBits: 8,
		StopBits: 1,
		Parity:   "N",
		Timeout:  5 * time.Second,
		RS485:    serial.RS485Config{},
	}

	rtuClient.Connect(conf, 1)
	fmt.Println(rtuClient.ReadCoils(17, 3))
	fmt.Println(rtuClient.WriteSingleCoil(17, client.OFF))
}
