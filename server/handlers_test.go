package server

import (
	"errors"
	"log"
	"modbusProxyServer/config"
	"net/http"
	"net/http/httptest"
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

func Init() {
	controller := Controller{}
	proxyServer.StartServer(config.Config{ProxyServerAddr: serverAddr}, &controller)
}

func TestServer_addIotDevice(t *testing.T) {
	go Init()
	log.SetFlags(log.Lshortfile)
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

func TestServer_addIotDeviceFails(t *testing.T) {
	log.SetFlags(log.Lshortfile)
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

func TestServer_getInformationFromIotDevice(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/device/metrics?deviceName=testName", nil)
	w := httptest.NewRecorder()
	proxyServer.getInformationFromIotDevice(w, req)

	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestServer_getInformationFromIotDeviceFail(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/device/metrics?deviceName=", nil)
	w := httptest.NewRecorder()
	proxyServer.getInformationFromIotDevice(w, req)

	if want, got := http.StatusInternalServerError, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestServer_getInformationFromIotDeviceFail2(t *testing.T) {
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

func TestServer_removeIotDevice(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/device/rm?deviceName=testName", nil)
	w := httptest.NewRecorder()
	proxyServer.rmIoTDevice(w, req)

	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestServer_removeIotDeviceFail(t *testing.T) {
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

func TestServer_getLogs(t *testing.T) {
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

func TestServer_getLogsFail(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/logs?countLogs=", nil)
	w := httptest.NewRecorder()
	proxyServer.getLogs(w, req)

	if want, got := http.StatusInternalServerError, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestServer_getLogsFail2(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/logs?countLogs=-1", nil)
	w := httptest.NewRecorder()
	proxyServer.getLogs(w, req)

	if want, got := http.StatusInternalServerError, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestServer_rmIotDevice(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/device/rm?deviceName=testName", nil)
	w := httptest.NewRecorder()
	proxyServer.rmIoTDevice(w, req)

	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestServer_rmIotDeviceFail(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/device/rm?deviceName=", nil)
	w := httptest.NewRecorder()
	proxyServer.rmIoTDevice(w, req)

	if want, got := http.StatusInternalServerError, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestServer_observeCoils(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet,
		"/device/observer/coils/start?deviceName=testName&address=1&quantity=1&time=1", nil)
	w := httptest.NewRecorder()
	proxyServer.observeDeviceCoils(w, req)

	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestServer_observeCoilsFail(t *testing.T) {
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

func TestServer_stopObserve(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet,
		"/device/observer/stop?deviceName=testName", nil)
	w := httptest.NewRecorder()
	proxyServer.stopObserveDevice(w, req)

	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestServer_stopObserveFail(t *testing.T) {
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

func TestServer_stopObserveFail2(t *testing.T) {
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
