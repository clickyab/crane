package gmaps

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/clickyab/services/assert"
)

var (
	data = []byte{0x00, 0x0e, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x1b, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00}
)

func intToByte(i int64) [4]byte {
	var res [4]byte
	pos := 4
	for i > 0 && pos > 0 {
		pos--
		res[pos] = byte(i % 256)
		i /= 256
	}
	return res
}

func byteToFloat(b1, b2, b3, b4 byte) float64 {
	return float64(b1) + float64(b2)*256 + float64(b3)*256*256 + float64(b4)*256*256*256
}

// LockUp try to lockup location
func LockUp(mcc, mnc, lac, cid int64) (float64, float64, error) {
	dataCP := data
	bMCC := intToByte(mcc)
	for i := range bMCC {
		dataCP[55-34+i] = bMCC[i]
		dataCP[55-12+i] = bMCC[i]
	}
	bMNC := intToByte(mnc)
	for i := range bMNC {
		dataCP[55-16+i] = bMNC[i]
		dataCP[55-38+i] = bMNC[i]
	}
	bLAC := intToByte(lac)
	for i := range bLAC {
		dataCP[55-20+i] = bLAC[i]
	}
	bCID := intToByte(cid)
	for i := range bCID {
		dataCP[55-24+i] = bCID[i]
	}

	if cid > 0xffff && mcc != 0 && mnc != 0 {
		dataCP[55-27] = 5
	} else {
		dataCP[55-24] = 0
		dataCP[55-23] = 0
	}
	buf := bytes.NewBuffer(dataCP)
	r, err := http.Post("http://www.google.com/glm/mmap", "application/binary", buf)
	if err != nil {
		return 0, 0, err
	}

	resp, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return 0, 0, err
	}
	if len(resp) > 14 {
		return byteToFloat(resp[10], resp[9], resp[8], resp[7]) / 1000000, byteToFloat(resp[14], resp[13], resp[12], resp[11]) / 1000000, nil
	}
	return 0, 0, errors.New("invalid response")
}

func init() {
	assert.True(len(data) == 55)
}
