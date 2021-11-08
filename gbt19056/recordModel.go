package gbt19056

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/leegggg/GBT19056-2020Gen/utils/encode"
)

// LengthMetadata ...
const LengthMetadata = 6

// PositionMultiplier ...
const PositionMultiplier = 10000.0 * 60

// HexUint8 ...
type HexUint8 uint8

// UnmarshalJSON HexUint8 ...
func (sd *HexUint8) UnmarshalJSON(input []byte) error {
	strInput := encode.BytesToStr(input)
	strInput = strings.Trim(strInput, `"`)
	res, err := strconv.ParseUint(strInput, 0, 8)
	if err != nil {
		return err
	}

	*sd = HexUint8(res)
	return nil
}

// MarshalJSON HexUint8
func (sd *HexUint8) MarshalJSON() ([]byte, error) {
	//do your serializing here
	stamp := "\"0x00\""
	if uint8(*sd) != 0x00 {
		stamp = fmt.Sprintf("\"0x%02x\"", uint8(*sd))
	}
	return []byte(stamp), nil
}

// Position ...
type Position struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Elevation float64 `json:"elevation"`
}

// LengthPosition ...
const LengthPosition = 10

// DumpData Position
func (e *Position) DumpData() ([]byte, error) {
	var longitude, latitude int32
	var elevation int16

	bs := make([]byte, 10)

	if e.Latitude == 0 && e.Longitude == 0 && e.Elevation == 0 {
		longitude = 0x7FFFFFFF
		latitude = 0x7FFFFFFF
		elevation = 0x7FFF
	} else {
		if math.Abs(e.Latitude) > 90 {
			latitude = 0x7FFFFFFF
		} else {
			latitude = int32(math.Round(e.Longitude * PositionMultiplier))
		}

		if math.Abs(e.Longitude) > 180 {
			longitude = 0x7FFFFFFF
		} else {
			longitude = int32(math.Round(e.Latitude * PositionMultiplier))
		}

		if math.Abs(e.Elevation) > 10000 {
			elevation = 0x7FFF
		} else {
			elevation = int16(e.Elevation)
		}

	}

	copy(bs[0:4], encode.Int32ToBytes(latitude))
	copy(bs[4:8], encode.Int32ToBytes(longitude))
	copy(bs[8:], encode.Int16ToBytes(elevation))
	return bs, nil
}

// LoadBinary ...
func (e *Position) LoadBinary(buffer []byte) {
	// Table A.20
	e.Longitude = float64(encode.BytesToInt32(buffer[0:4])) / PositionMultiplier
	e.Latitude = float64(encode.BytesToInt32(buffer[4:8])) / PositionMultiplier
	e.Elevation = float64(encode.BytesToInt16(buffer[8:10]))
	return
}

// dataBlockMeta ...
type dataBlockMeta struct {
	SynB1 HexUint8 `json:"syn_b1"` // need decode from 0x01 ...
	SynB2 HexUint8 `json:"syn_b2"` // need decode from 0x01 ...
	FmtM  HexUint8 `json:"mfmt"`   // need decode from 0x01 ...
	FmtS  HexUint8 `json:"sfmt"`   // need decode from 0x01 ...
	Size  uint32   `json:"size"`
}

// DumpDate ...
func (e *dataBlockMeta) DumpData() ([]byte, error) {
	bs := make([]byte, 6)
	bs[0] = (byte)(e.SynB1)
	bs[1] = (byte)(e.SynB2)
	bs[2] = (byte)(e.FmtM)
	bs[3] = (byte)(e.FmtS)
	size := e.Size / 16
	if e.Size%16 != 0 {
		return nil, errors.New("size is not in 16B, padding maybe needed")
	}

	binary.BigEndian.PutUint16(bs[4:], (uint16)(size))
	return bs, nil
}

// LoadDate return datablock size
func (e *dataBlockMeta) LoadBinary(buffer []byte) (int, error) {
	// Table B.2
	var err error
	e.SynB1 = HexUint8(buffer[0])
	e.SynB2 = HexUint8(buffer[1])
	e.FmtM = HexUint8(buffer[2])
	e.FmtS = HexUint8(buffer[3])
	e.Size = (uint32)(binary.BigEndian.Uint16(buffer[4:6])) * 16
	return int(e.Size), err
}

func addPaddingBytes(body []byte) ([]byte, error) {
	totalLength := len(body) + LengthMetadata
	paddingLength := (totalLength/16+1)*16 - totalLength
	return append(body, make([]byte, paddingLength)...), nil
}

func calcXorSum(body []byte) (byte, error) {
	var sum byte = 97
	for _, v := range body {
		sum ^= v
	}
	return sum, nil
}

// linkDumpedData
func (e dataBlockMeta) linkDumpedData(body []byte) ([]byte, error) {
	body, _ = addPaddingBytes(body)
	sum, _ := calcXorSum(body)
	body[len(body)-1] = sum
	e.Size = (uint32)(len(body))
	meta, err := e.DumpData()
	bs := append(meta, body...)
	return bs, err
}
