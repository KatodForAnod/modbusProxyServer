package server

import (
	"errors"
	"log"
	"modbusProxyServer/config"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

var (
	proxyServer Server
)

const serverAddr = "127.0.0.1:8080"

type Controller struct {
	ioTs        []config.IotConfig
	iotObserver string
}

func (c *Controller) AddIoTDevice(device config.IotConfig) error {
	if device.DeviceName == "" {
		return errors.New("incorrect device name")
	}

	for _, t := range c.ioTs {
		if t.DeviceName == device.DeviceName {
			return errors.New("device with such name already exist")
		}
	}

	c.ioTs = append(c.ioTs, device)
	return nil
}

func (c *Controller) RmIoTDevice(deviceName string) error {
	for _, t := range c.ioTs {
		if t.DeviceName == deviceName {
			return nil
		}
	}

	return errors.New("not found")
}

func (c *Controller) StopObserveDevice(deviceName string) error {
	if c.iotObserver == "" {
		return errors.New("device not found")
	}

	c.iotObserver = ""
	return nil
}

func (c *Controller) ObserveIoTCoils(deviceName, address, quantity, timeSecondsDuration string) error {
	_, err := strconv.ParseInt(timeSecondsDuration, 10, 64)
	if err != nil {
		log.Println(err)
		return err
	}

	c.iotObserver = deviceName
	return nil
}

func (c *Controller) GetLastNRowsLogs(nRows int) ([]string, error) {
	if nRows < 0 {
		return []string{}, errors.New("wrong count rows")
	}
	return []string{"1 row", "2 row"}, nil
}

func (c *Controller) RemoveIoTDeviceObserve(ioTsConfig []config.IotConfig) error {
	return nil
}

func (c *Controller) NewIotDeviceObserve(iotConfig config.IotConfig) error {
	c.ioTs = append(c.ioTs, iotConfig)
	return nil
}

func (c *Controller) GetInformation(deviceName string) ([]byte, error) {
	for i, t := range c.ioTs {
		if t.DeviceName == deviceName {
			return []byte{byte(i)}, nil
		}
	}

	return []byte{}, errors.New("not found")
}

func TestServerInit(t *testing.T) {
	controller := Controller{}
	go proxyServer.StartServer(config.Config{ProxyServerAddr: serverAddr}, &controller)
}

func TestAddRTUClient(t *testing.T) {
	myReader := strings.NewReader(`{"device_name":"testName",
        "type_client":"rtu",
        "slave_id": 1,
        "com_port":"COM3",
        "baud_rate": 115200,
        "data_bits": 8,
        "stop_bits": 1,
        "parity":"N",
        "timeout_seconds":5}`)

	req := httptest.NewRequest(http.MethodGet, "/device/add", myReader)
	w := httptest.NewRecorder()

	proxyServer.addIoTDevice(w, req)
	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestAddClientFailEmptyDeviceName(t *testing.T) {
	myReader := strings.NewReader(`{"device_name":"",
        "type_client":"rtu",
        "slave_id": 1,
        "com_port":"COM3",
        "baud_rate": 115200,
        "data_bits": 8,
        "stop_bits": 1,
        "parity":"N",
        "timeout_seconds":5}`)

	req := httptest.NewRequest(http.MethodGet, "/device/add", myReader)
	w := httptest.NewRecorder()
	proxyServer.addIoTDevice(w, req)
	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestAddWrongClientType(t *testing.T) {
	myReader := strings.NewReader(`{"device_name":"testName",
        "type_client":"notValid",
        "slave_id": 1,
        "com_port":"COM3",
        "baud_rate": 115200,
        "data_bits": 8,
        "stop_bits": 1,
        "parity":"N",
        "timeout_seconds":5}`)

	req := httptest.NewRequest(http.MethodGet, "/device/add", myReader)
	w := httptest.NewRecorder()

	proxyServer.addIoTDevice(w, req)
	if want, got := http.StatusInternalServerError, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestGetInformationFromIotDevice(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/device/metrics?deviceName=testName", nil)
	w := httptest.NewRecorder()
	proxyServer.getInformationFromIotDevice(w, req)

	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestGetInformationFromEmptyDeviceName(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/device/metrics?deviceName=", nil)
	w := httptest.NewRecorder()
	proxyServer.getInformationFromIotDevice(w, req)

	if want, got := http.StatusInternalServerError, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestGetInformationFromEmptyIotDevice(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/device/metrics", nil)
	w := httptest.NewRecorder()
	proxyServer.getInformationFromIotDevice(w, req)

	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}

	if w.Body.String() == "" {
		t.Fatalf("expected warning msg")
	}
}

func TestRemoveIotDevice(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/device/rm?deviceName=testName", nil)
	w := httptest.NewRecorder()
	proxyServer.rmIoTDevice(w, req)

	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestRemoveIotDeviceFail(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/device/rm?deviceNam", nil)
	w := httptest.NewRecorder()
	proxyServer.rmIoTDevice(w, req)

	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}

	if w.Body.String() == "" {
		t.Fatalf("expected warning msg")
	}
}

func TestGetLogs(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/logs?countLogs=2", nil)
	w := httptest.NewRecorder()
	proxyServer.getLogs(w, req)

	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}

	out := w.Body.String()
	outArr := strings.Split(out, "\n")
	if len(outArr) < 2 {
		t.FailNow()
	}
}

func TestGetLogsFailEmptyCount(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/logs?countLogs=", nil)
	w := httptest.NewRecorder()
	proxyServer.getLogs(w, req)

	if want, got := http.StatusInternalServerError, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestGetLogsFail(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/logs?", nil)
	w := httptest.NewRecorder()
	proxyServer.getLogs(w, req)

	if want, got := http.StatusInternalServerError, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestGetLogsFailNegativeCount(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/logs?countLogs=-1", nil)
	w := httptest.NewRecorder()
	proxyServer.getLogs(w, req)

	if want, got := http.StatusInternalServerError, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestRmIotDevice(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/device/rm?deviceName=testName", nil)
	w := httptest.NewRecorder()
	proxyServer.rmIoTDevice(w, req)

	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestRmIotDeviceFail(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/device/rm?deviceName=", nil)
	w := httptest.NewRecorder()
	proxyServer.rmIoTDevice(w, req)

	if want, got := http.StatusInternalServerError, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestObserveCoils(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet,
		"/device/observer/coils/start?deviceName=testName&address=1&quantity=1&time=1", nil)
	w := httptest.NewRecorder()
	proxyServer.observeDeviceCoils(w, req)

	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestObserveCoilsFailEmptyQuantity(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet,
		"/device/observer/coils/start?deviceName=testName&address=1&time=1", nil)
	w := httptest.NewRecorder()
	proxyServer.observeDeviceCoils(w, req)

	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}

	if w.Body.String() == "" {
		t.Fatalf("expected warning msg")
	}
}

func TestObserveCoilsFailEmptyTime(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet,
		"/device/observer/coils/start?deviceName=testName&quantity=1&address=1", nil)
	w := httptest.NewRecorder()
	proxyServer.observeDeviceCoils(w, req)

	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}

	if w.Body.String() == "" {
		t.Fatalf("expected warning msg")
	}
}

func TestObserveCoilsFailWrongTimeField(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet,
		"/device/observer/coils/start?deviceName=testName&address=1&quantity=1&time=wrong", nil)
	w := httptest.NewRecorder()
	proxyServer.observeDeviceCoils(w, req)

	if want, got := http.StatusInternalServerError, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestObserveCoilsFailEmptyAddress(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet,
		"/device/observer/coils/start?deviceName=testName&quantity=1&time=1", nil)
	w := httptest.NewRecorder()
	proxyServer.observeDeviceCoils(w, req)

	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}

	if w.Body.String() == "" {
		t.Fatalf("expected warning msg")
	}
}

func TestObserveCoilsFailEmptyDeviceName(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet,
		"/device/observer/coils/start?address=1&quantity=1&time=1", nil)
	w := httptest.NewRecorder()
	proxyServer.observeDeviceCoils(w, req)

	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}

	if w.Body.String() == "" {
		t.Fatalf("expected warning msg")
	}
}

func TestStopObserve(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet,
		"/device/observer/stop?deviceName=testName", nil)
	w := httptest.NewRecorder()
	proxyServer.stopObserveDevice(w, req)

	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestStopObserveFail(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet,
		"/device/observer/stop?deviceName", nil)
	w := httptest.NewRecorder()
	proxyServer.stopObserveDevice(w, req)

	if want, got := http.StatusInternalServerError, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}

	if w.Body.String() == "" {
		t.Fatalf("expected warning msg")
	}
}

func TestStopObserveFailEmptyDeviceName(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet,
		"/device/observer/stop?", nil)
	w := httptest.NewRecorder()
	proxyServer.stopObserveDevice(w, req)

	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}

	if w.Body.String() == "" {
		t.Fatalf("expected warning msg")
	}
}
