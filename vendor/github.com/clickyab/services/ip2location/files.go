package ip2location

import (
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
	file *os.File
}

func newFileMock() (*fileMock, error) {
	f, err := os.Open(filepath.Join(pwd, fl.String()))
	if err != nil {
		return nil, err
	}

	return &fileMock{
		file: f,
	}, nil
}

func (fm *fileMock) ReadAt(b []byte, off int64) (n int, err error) {
	return fm.file.ReadAt(b, off)
}

func init() {
	var err error
	pwd, err = expand.Pwd()
	assert.Nil(err)
	fl = config.RegisterString("services.ip2location.datafile", "IP-COUNTRY-REGION-CITY-ISP.BIN", "location of data file")
}
