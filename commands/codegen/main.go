package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"clickyab.com/exchange/services/assert"
	_ "clickyab.com/exchange/services/codegen/datatable" // Datatable
	_ "clickyab.com/exchange/services/codegen/enums"     // Enums plugin
	_ "clickyab.com/exchange/services/codegen/gin"       // Gin plugin
	_ "clickyab.com/exchange/services/codegen/models"    // Models plugin
	"clickyab.com/exchange/services/codegen/plugins"
	_ "clickyab.com/exchange/services/codegen/swagger"  // Raml plugin
	_ "clickyab.com/exchange/services/codegen/validate" // Validateor

	"github.com/Sirupsen/logrus"
	"github.com/goraz/humanize"
	"github.com/ogier/pflag"
)

var (
	pkg = pflag.StringP("package", "p", "", "the package to scan for gin controller")
)

func main() {
	pflag.Parse()

	if *pkg == "" {
		path, err := os.Getwd()
		assert.Nil(err)
		// load the current directory as package
		gopath := append(strings.Split(os.Getenv("GOPATH"), ":"), runtime.GOROOT())
		for i := range gopath {
			p := filepath.Join(gopath[i], "src")
			if strings.HasPrefix(path, p) {
				*pkg = strings.Trim(strings.TrimPrefix(path, p), "/")
				break
			}
		}
	}
	//fmt.Println(os.Environ())
	p, err := humanize.ParsePackage(*pkg)
	if err != nil {
		logrus.Fatal(err)
	}

	err = plugins.ProcessPackage(*p)
	if err != nil {
		logrus.Fatal(err)
	}

	err = plugins.Finalize(*p)
	if err != nil {
		logrus.Fatal(err)
	}
}
