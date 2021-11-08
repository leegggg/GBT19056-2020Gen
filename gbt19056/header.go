package gbt19056

import "encoding/binary"

// RecoderID ..
type Header struct {
	dataBlockMeta
	Year         uint16    `json:"year"`
	Filename     string    `json:"filename"`
	CreateTime   DateTime  `json:"create_time"`
	Size         uint32    `json:"size"`
	Description  string    `json:"description"`
	RecoderID    RecoderID `json:"recoder_id"`
	Plate        string    `json:"plate"`
	PlateType    string    `json:"plate_type"`
	VIN          string    `json:"vin"`
	Sn           uint64    `json:"sn"`
	PulseFactor  uint16    `json:"pulse_factor"`
	Installation DateTime  `json:"dop,string"`
	HashSum      []uint8   `json:"hash_sum"`
}

// DumpData RecoderID
func (e *RecoderID) DumpData() ([]byte, error) {
	var err error

	// ASCII ID should be safe to be copy directly
	ccc := make([]byte, 7)
	copy(ccc, []byte(e.CCC))

	// TODO: Check length
	version := make([]byte, 16)
	copy(version, []byte(e.Version))

	var dop []byte
	dop, err = e.Dop.DumpDataShort()

	sn := make([]byte, 4)
	binary.BigEndian.PutUint32(sn, e.Sn)

	// TODO： Unfinished comment logic
	comment := make([]byte, 5)
	copy(comment, []byte(e.Comment))

	bs := append(ccc, version...)
	bs = append(bs, dop...)
	bs = append(bs, sn...)
	bs = append(bs, comment...)

	bs, err = e.linkDumpedData(bs)
	return bs, err
}

// LoadBinary RecoderID Table A.11, Code 0x08
func (e *RecoderID) LoadBinary(buffer []byte, meta dataBlockMeta) {
	e.dataBlockMeta = meta
	e.CCC = bytesToStr(buffer[0:7])
	e.Version = bytesToStr(buffer[7:23])
	e.Dop.LoadBinaryShort(buffer[23:26])
	e.Sn = binary.BigEndian.Uint32(buffer[26:30])
	// e.Comment = bytesToStr(buffer[30:35])
	e.Comment = buffer[30:35]
	return
}
