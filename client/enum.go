package client

import "errors"

type CoilType [2]byte

var ON = CoilType{0xFF, 0x00}
var OFF = CoilType{0x00, 0x00}

func (t *CoilType) ToByteSlice() []byte {
	return []byte{t[0], t[1]}
}

func (t *CoilType) ToUint16() (coil uint16) {
	if *t == ON {
		return 0xFF00
	} else {
		return 0x0000
	}
}

func (t *CoilType) String() (str string) {
	if *t == ON {
		str = "ON"
	} else {
		str = "OFF"
	}

	return str
}

func BytesToCoilType(bytes []byte) (CoilType, error) {
	if len(bytes) != 2 {
		return CoilType{}, errors.New("wront type")
	}

	if bytes[0] == 255 && bytes[1] == 0 {
		return ON, nil
	} else if bytes[0] == 0 && bytes[1] == 0 {
		return OFF, nil
	}

	return CoilType{}, errors.New("wront type")
}

type ClientType struct {
	Cl string
}

const (
	TCPClientType = "tcp"
	RTUClientType = "rtu"
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
