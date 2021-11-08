package gbt19056

import (
	"fmt"
	"time"

	"github.com/leegggg/GBT19056-2020Gen/utils/bcd"
)

// DateTime ...
type DateTime struct {
	time.Time
}

// DumpData DateTime
func (e *DateTime) DumpData() ([]byte, error) {
	bs := make([]byte, 6)
	if e.Year() < 2000 || e.Year() > 2099 {
		err := fmt.Errorf("year %d not in range [2000-2099]", e.Year())
		return bs, err
	}

	var nb uint8
	nb = uint8(e.Year() - 2000)
	bs[0] = bcd.FromUint8(nb)

	nb = uint8(e.Month())
	bs[1] = bcd.FromUint8(nb)

	nb = uint8(e.Day())
	bs[2] = bcd.FromUint8(nb)

	nb = uint8(e.Hour())
	bs[3] = bcd.FromUint8(nb)

	nb = uint8(e.Minute())
	bs[4] = bcd.FromUint8(nb)

	nb = uint8(e.Second())
	bs[5] = bcd.FromUint8(nb)
	return bs, nil
}

// LoadBinary RealTime Table A.8, Code 0x02
func (e *DateTime) LoadBinary(buffer []byte) {
	year := bcd.ToUint8(buffer[0])
	month := bcd.ToUint8(buffer[1])
	day := bcd.ToUint8(buffer[2])
	hour := bcd.ToUint8(buffer[3])
	min := bcd.ToUint8(buffer[4])
	sec := bcd.ToUint8(buffer[5])
	e.Time = time.Date(
		int(year)+2000, time.Month(int(month)), int(day), int(hour), int(min), int(sec), 0, time.UTC)
}

// LoadBinaryShort RealTime Table A.14
func (e *DateTime) LoadBinaryShort(buffer []byte) {
	year := bcd.ToUint8(buffer[0])
	month := bcd.ToUint8(buffer[1])
	day := bcd.ToUint8(buffer[2])

	e.Time = time.Date(
		int(year)+2000, time.Month(int(month)), int(day), 0, 0, 0, 0, time.UTC)
}

// DumpDataShort DateTime
func (e *DateTime) DumpDataShort() ([]byte, error) {
	bs := make([]byte, 3)
	if e.Year() < 2000 || e.Year() > 2099 {
		err := fmt.Errorf("Year %d not in range [2000-2099]", e.Year())
		return bs, err
	}

	var nb uint8
	nb = uint8(e.Year() - 2000)
	bs[0] = bcd.FromUint8(nb)

	nb = uint8(e.Month())
	bs[1] = bcd.FromUint8(nb)

	nb = uint8(e.Day())
	bs[2] = bcd.FromUint8(nb)

	return bs, nil
}
