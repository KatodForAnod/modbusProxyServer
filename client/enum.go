package client

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
