package main

import (
	"os"
	"path/filepath"
	"runtime"
	"services/assert"
	_ "services/codegen/datatable" // Datatable
	_ "services/codegen/enums"     // Enums plugin
	_ "services/codegen/gin"       // Gin plugin
	_ "services/codegen/models"    // Models plugin
	"services/codegen/plugins"
	_ "services/codegen/swagger"  // Raml plugin
	_ "services/codegen/validate" // Validateor
	"strings"

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
