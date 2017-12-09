package ip2location

import (
	"io"
	"io/ioutil"
	"os"

	"path/filepath"

	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
	"github.com/fzerorubigd/expand"
	"gopkg.in/fzerorubigd/onion.v3"
)

var (
	fl  onion.String
	pwd string
)

type fileMock struct {
	data []byte
	ln   int
}

func newFileMock() (*fileMock, error) {
	f, err := os.Open(filepath.Join(pwd, fl.String()))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return &fileMock{
		data: data,
		ln:   len(data),
	}, nil
}

func (fm *fileMock) ReadAt(b []byte, off int64) (n int, err error) {
	lb := len(b)
	avail := int64(fm.ln) - off
	if avail < int64(lb) {
		return 0, io.EOF
	}

	copy(b, fm.data[off:off+int64(lb)])
	return lb, nil
}

func init() {
	var err error
	pwd, err = expand.Pwd()
	assert.Nil(err)
	fl = config.RegisterString("services.ip2location.datafile", "IP-COUNTRY-REGION-CITY-ISP.BIN", "location of data file")
}
