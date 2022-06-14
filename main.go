package main

import (
	"fmt"
	"github.com/goburrow/modbus"
	"log"
	"time"
)

func main() {
	handler := modbus.NewRTUClientHandler("COM3")

	handler.BaudRate = 115200
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.SlaveId = 1
	handler.Timeout = 5 * time.Second

	err := handler.Connect()
	if err != nil {
		return
	}
	defer handler.Close()

	client := modbus.NewClient(handler)
	results, err := client.WriteSingleCoil(16, 1)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(results)


	
}
