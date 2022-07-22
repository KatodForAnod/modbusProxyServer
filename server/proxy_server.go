package server

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"modbusProxyServer/config"
	"modbusProxyServer/controller"
	"net/http"
	"strconv"
	"strings"
)

type Server struct {
	controller controller.ServerController
}

func (s *Server) addIoTDevice(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var iotDev config.IotConfig
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&iotDev); err != nil {
		log.Errorln(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := s.controller.AddIoTDevice(iotDev); err != nil {
		log.Errorln(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) rmIoTDevice(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	deviceNames := r.URL.Query()["deviceName"]
	if len(deviceNames) == 0 {
		http.Error(w, "set device name", http.StatusBadRequest)
		return
	}
	deviceName := deviceNames[0]

	if err := s.controller.RmIoTDevice(deviceName); err != nil {
		log.Errorln(err)
		http.Error(w, "wrong device name", http.StatusInternalServerError)
		return
	}
}

func (s *Server) stopObserveDevice(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	deviceNames := r.URL.Query()["deviceName"]
	if len(deviceNames) == 0 {
		http.Error(w, "set device name", http.StatusBadRequest)
		return
	}
	deviceName := deviceNames[0]

	if err := s.controller.StopObserveDevice(deviceName); err != nil {
		log.Errorln(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) getInformationFromIotDevice(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	deviceNames := r.URL.Query()["deviceName"]
	if len(deviceNames) == 0 {
		http.Error(w, "set device name", http.StatusBadRequest)
		return
	}
	deviceName := deviceNames[0]

	inf, err := s.controller.GetInformation(deviceName)
	if err != nil {
		log.Errorln(err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(inf)
	if err != nil {
		log.Errorln(err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) getLogs(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	countLogsArr := r.URL.Query()["countLogs"]
	if len(countLogsArr) == 0 {
		http.Error(w, "set device name", http.StatusBadRequest)
		return
	}
	countLogsStr := countLogsArr[0]
	countLogs, err := strconv.Atoi(countLogsStr)
	if err != nil {
		log.Errorln(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logs, err := s.controller.GetLastNRowsLogs(countLogs)
	if err != nil {
		log.Errorln(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	allLogs := strings.Join(logs, "\n")
	_, err = w.Write([]byte(allLogs))
	if err != nil {
		log.Errorln(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) observeDeviceCoils(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	for _, s2 := range []string{"deviceName", "address", "quantity", "time"} {
		if !r.URL.Query().Has(s2) {
			http.Error(w, "parameter is missing: "+s2, http.StatusBadRequest)
			return
		}
	}

	deviceName := r.URL.Query().Get("deviceName")
	address := r.URL.Query().Get("address")
	quantity := r.URL.Query().Get("quantity")
	time := r.URL.Query().Get("time")

	err := s.controller.ObserveIoTCoils(deviceName, address, quantity, time)
	if err != nil {
		log.Errorln(err)
		http.Error(w, "wrong params", http.StatusInternalServerError)
		return
	}
}

func (s *Server) StartServer(config config.Config, controller controller.ServerController) {
	s.controller = controller

	http.HandleFunc("/device/metrics", s.getInformationFromIotDevice)
	http.HandleFunc("/logs", s.getLogs)
	http.HandleFunc("/device/add", s.addIoTDevice)
	http.HandleFunc("/device/rm", s.rmIoTDevice)
	http.HandleFunc("/device/observer/stop", s.stopObserveDevice)
	http.HandleFunc("/device/observer/coils/start", s.observeDeviceCoils)

	fmt.Println("Server is listening... ", config.ProxyServerAddr)
	log.Fatal(http.ListenAndServe(config.ProxyServerAddr, nil))
}
