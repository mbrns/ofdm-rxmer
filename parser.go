package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"time"
)

type Frequency uint64

func (f Frequency) Int() int {
	return int(f)
}

const (
	Hertz     Frequency = 1
	Kilohertz           = 1000 * Hertz
)

type scMer struct {
	Idx  int       `json:"sc_index"`
	Freq Frequency `json:"freq"`
	Mer  float64   `json:"mer"`
}

type PnmDsMerFile struct {
	FileType           []byte    `json:"file_type"`
	CaptureTime        time.Time `json:"capture_time"`
	DsChanId           int       `json:"ds_chan_id"`
	CmMacAddress       string    `json:"cm_mac_address"`
	ScZeroFreq         Frequency `json:"sc_zero_freq"`
	FirstActiveScIndex int       `json:"first_active_sc_index"`
	ScSpacing          Frequency `json:"sc_spacing"`
	RxmerDataLenth     int       `json:"rxmer_data_lenth"`
	RxmerData          []scMer   `json:"rxmer_data"`
}

func (r PnmDsMerFile) GetStartingFrequency() Frequency {
	plus := r.FirstActiveScIndex * int(r.ScSpacing)
	return r.ScZeroFreq + Frequency(plus)
}

func parseRxmerFile(data []byte) *PnmDsMerFile {
	pnmFile := PnmDsMerFile{}
	var merData []byte = nil

	fileType := data[:4] // first 4 bytes
	pnmFile.FileType = fileType
	// fmt.Printf("File Type: %x\n", pnmFile.FileType)

	// true if file type matches. Per spec it should be '504e4e04' or '504e4d04'.
	var pnmDsMerFileType1 = []byte{80, 78, 77, 4}
	var pnmDsMerFileType2 = []byte{80, 78, 78, 4}
	if bytes.Equal(fileType, pnmDsMerFileType1) {
		log.Println("Found match for type1", pnmDsMerFileType1)
		merData = data[4:]
	} else if bytes.Equal(fileType, pnmDsMerFileType2) {
		log.Println("Found match for type2", pnmDsMerFileType2)
		merData = data[6:]

	} else {
		log.Fatalf("ERROR: File is not PNMDsMer. Type %x does not match required type %x or %x\n", fileType, pnmDsMerFileType1, pnmDsMerFileType2)
	}

	captureTime := merData[:4] // 4 bytes
	pnmFile.CaptureTime = time.Unix(int64(makeInt(captureTime)), 0)
	// fmt.Println("Capture Time: ", pnmFile.CaptureTime)

	dsChanId := merData[4:5] // 1 byte
	pnmFile.DsChanId = makeInt(dsChanId)
	// fmt.Printf("DS Channel ID: %d\n", pnmFile.DsChanId)

	var cmMacAddress net.HardwareAddr = merData[5:11] // 6 bytes
	pnmFile.CmMacAddress = cmMacAddress.String()
	// fmt.Printf("CM MAC Address: %s\n", pnmFile.CmMacAddress)

	scZeroFreq := merData[11:15] // 4 bytes
	pnmFile.ScZeroFreq = Frequency(makeInt(scZeroFreq))
	// fmt.Printf("Subcarrier zero frequency in Hz: %d\n", pnmFile.ScZeroFreq)

	firstActiveSc := merData[15:17] // 2 bytes
	pnmFile.FirstActiveScIndex = makeInt(firstActiveSc)
	// fmt.Println("First Active Subcarrier Index: ", pnmFile.FirstActiveScIndex)

	scSpacing := merData[17:18] // 1 byte ; in kHz
	pnmFile.ScSpacing = Frequency(makeInt(scSpacing)) * Kilohertz
	// fmt.Printf("Subcarrier Spacing in Hz: %d\n", pnmFile.ScSpacing)

	rxMerDataLength := merData[18:22] // 4 bytes
	pnmFile.RxmerDataLenth = makeInt(rxMerDataLength)
	// fmt.Printf("Length in bytes of RxMerData that follows: %d\n", pnmFile.RxmerDataLenth)

	rxMerRaw := merData[22:] // rest of file
	// fmt.Printf("Subcarrier RxMER Data: %x\n", rxMerRaw)

	mers := make([]scMer, 0, pnmFile.RxmerDataLenth)
	for k, v := range rxMerRaw {
		m := scMer{
			Idx:  k + pnmFile.FirstActiveScIndex,
			Freq: pnmFile.GetStartingFrequency() + Frequency(k*pnmFile.ScSpacing.Int()),
			Mer:  float64(v) * .25,
		}
		mers = append(mers, m)
	}
	pnmFile.RxmerData = mers
	//fmt.Println(mers)

	return &pnmFile
}

func makeInt(xb []byte) int {
	buf := bytes.NewBuffer(xb)
	switch len(xb) {
	case 1:
		var b uint8
		_ = binary.Read(buf, binary.BigEndian, &b)
		return int(b)
	case 2:
		var b uint16
		_ = binary.Read(buf, binary.BigEndian, &b)
		return int(b)
	case 4:
		var b uint32
		_ = binary.Read(buf, binary.BigEndian, &b)
		return int(b)
	}
	return -1
}
