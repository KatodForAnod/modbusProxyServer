package server

import (
	"fmt"
	"log"
	"modbusprottocol/config"
	"modbusprottocol/controller"
	"net/http"
	"strconv"
	"strings"
)

type Server struct {
	controller controller.Controller
}

func (s *Server) getInformationFromIotDevice(w http.ResponseWriter, r *http.Request) {
	log.Println("handler getInformationFromIotDevice")
	defer r.Body.Close()
	deviceNames := r.URL.Query()["deviceName"]
	if len(deviceNames) == 0 {
		log.Println("device name not found")
		fmt.Fprintf(w, "set device name")
		return
	}
	deviceName := deviceNames[0]

	inf, err := s.controller.GetInformation(deviceName)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(inf)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) getLogs(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	log.Println("handler getLogs")
	countLogsArr := r.URL.Query()["countLogs"]
	if len(countLogsArr) == 0 {
		log.Println("count logs not found")
		fmt.Fprintf(w, "set count logs")
		return
	}
	countLogsStr := countLogsArr[0]
	countLogs, err := strconv.Atoi(countLogsStr)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logs, err := s.controller.GetLastNRowsLogs(countLogs)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	allLogs := strings.Join(logs, "\n")
	_, err = w.Write([]byte(allLogs))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) StartServer(config config.Config, controller controller.Controller) {
	s.controller = controller

	http.HandleFunc("/device/metrics", s.getInformationFromIotDevice)
	http.HandleFunc("/logs", s.getLogs)

	fmt.Println("Server is listening... ", config.ProxyServerAddr)
	log.Fatal(http.ListenAndServe(config.ProxyServerAddr, nil))
}
