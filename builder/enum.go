package builder

import "errors"

type ClientType struct {
	Cl string
}

const (
	TCPClient = "tcp"
	RTUClient = "rtu"
)

func (t *ClientType) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case `"tcp"`, `tcp`:
		t.Cl = "tcp"
		return nil
	case `"rtu"`, `rtu`:
		t.Cl = "rtu"
		return nil
	default:
		return errors.New("CustomBool: parsing \"" + string(data) + "\": unknown value")
	}
}

func (t *ClientType) String() string {
	return t.Cl
}
