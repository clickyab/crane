package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/clickyab/services/assert"
	_ "github.com/clickyab/services/codegen/datatable"    // Datatable
	_ "github.com/clickyab/services/codegen/enums"        // Enums plugin
	_ "github.com/clickyab/services/codegen/gin"          // Gin plugin
	_ "github.com/clickyab/services/codegen/models_mysql" // Models plugin
	"github.com/clickyab/services/codegen/plugins"
	_ "github.com/clickyab/services/codegen/swagger"  // Raml plugin
	_ "github.com/clickyab/services/codegen/validate" // Validateor

	"github.com/goraz/humanize"
	"github.com/ogier/pflag"
	"github.com/sirupsen/logrus"
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
