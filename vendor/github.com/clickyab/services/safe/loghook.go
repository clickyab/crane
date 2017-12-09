package safe

import (
	"net/http"
	"net/http/httputil"

	"github.com/sirupsen/logrus"
)

type logPrinter struct {
}

func (logPrinter) Recover(err error, stack []byte, extra ...interface{}) {
	logrus.Errorf("Error '%s' is recovered, stack was :\n\n%s", err, string(stack))

	for i := range extra {
		if t, ok := extra[i].(*http.Request); ok {
			if b, err := httputil.DumpRequest(t, true); err == nil {
				logrus.Errorf("the https request dump : \n\n%s", string(b))
				continue
			}
		}

		logrus.Errorf("Extra data :\n %T => %+v", extra[i], extra[i])
	}
}

func init() {
	Register(&logPrinter{})
}
