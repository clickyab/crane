package ip2location

import "io"

type fileMock struct {
	data []byte
	ln   int
}

func newFileMock() (*fileMock, error) {
	d, err := Asset("IP-COUNTRY-REGION-CITY.BIN")
	if err != nil {
		return nil, err

	}

	return &fileMock{
		data: d,
		ln:   len(d),
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
