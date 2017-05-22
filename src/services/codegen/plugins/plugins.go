package plugins

import (
	"services/codegen/annotate"
	"strings"
	"sync"

	"sort"

	"github.com/Sirupsen/logrus"
	"github.com/goraz/humanize"
)

// AnnotationPlug is the base plugin system
type AnnotationPlug interface {
	// GetOrder return the order of this plugin in system
	GetOrder() int
	// GetType return all types that this plugin can operate on
	// for example if the result contain Route then all @Route sections are
	// passed to this plugin
	GetType() []string
	// Finalize is called after all the functions are done. the context is the one from the
	// process
	Finalize(interface{}, humanize.Package) error
}

// AnnotationFunction is the plugin for the functions
type AnnotationFunction interface {
	AnnotationPlug
	// FunctionIsSupported check for a function signature and if the function is supported in this
	// interface
	FunctionIsSupported(humanize.File, humanize.Function) bool
	// ProcessFunction the function with its annotation. any error here means to stop the
	// all process
	// the first argument is the context. if its nil, means its the first run for this package.
	// the result of this function is passed to the plugin next time for the next function
	ProcessFunction(interface{}, humanize.Package, humanize.File, humanize.Function, annotate.Annotate) (interface{}, error)
}

// AnnotationStruct is the plugin for the struct types
type AnnotationStruct interface {
	AnnotationPlug
	// StructureIsSupported check for a function signature and if the function is supported in this
	// interface
	StructureIsSupported(humanize.File, humanize.TypeName) bool
	// ProcessStructure the structure with its annotation. any error here means to stop the
	// all process
	// the first argument is the context. if its nil, means its the first run for this package.
	// the result of this function is passed to the plugin next time for the next function
	ProcessStructure(interface{}, humanize.Package, humanize.File, humanize.TypeName, annotate.Annotate) (interface{}, error)
}

// AnnotationStruct is the plugin for the struct types
type AnnotationType interface {
	AnnotationPlug
	// StructureIsSupported check for a function signature and if the function is supported in this
	// interface
	TypeIsSupported(humanize.File, humanize.TypeName) bool
	// ProcessStructure the structure with its annotation. any error here means to stop the
	// all process
	// the first argument is the context. if its nil, means its the first run for this package.
	// the result of this function is passed to the plugin next time for the next function
	ProcessType(interface{}, humanize.Package, humanize.File, humanize.TypeName, annotate.Annotate) (interface{}, error)
}

type plugin struct {
	p       AnnotationPlug
	context interface{}
	called  bool
}

type pluginList []plugin

func (pl pluginList) Len() int {
	return len(pl)
}

func (pl pluginList) Less(i, j int) bool {
	return pl[i].p.GetOrder() < pl[j].p.GetOrder()
}

func (pl pluginList) Swap(i, j int) {
	pl[i], pl[j] = pl[j], pl[i]
}

var (
	plugins pluginList
	lock    = sync.Mutex{}
)

// Register a new plugin
func Register(p AnnotationPlug) {
	lock.Lock()
	defer lock.Unlock()

	switch p.(type) {
	case AnnotationFunction:
	case AnnotationStruct:
	case AnnotationType:
	default:
		logrus.Panic("Plugin type is not supported")
	}

	plugins = append(plugins, plugin{
		p: p,
	})
	sort.Sort(plugins)
}

func inArray(n string, h ...string) bool {
	for i := range h {
		if n == h[i] {
			return true
		}
	}

	return false
}

// ProcessPackage start the process for a package
func ProcessPackage(p humanize.Package) error {
	for f := range p.Files {
		for t := range p.Files[f].Types {
			if _, ok := p.Files[f].Types[t].Type.(*humanize.StructType); ok {
				a, err := annotate.LoadFromComment(strings.Join(p.Files[f].Types[t].Docs, "\n"))
				if err != nil {
					return err
				}
				err = processStructure(p, *p.Files[f], *p.Files[f].Types[t], a)
				if err != nil {
					return err
				}
			}
			a, err := annotate.LoadFromComment(strings.Join(p.Files[f].Types[t].Docs, "\n"))
			if err != nil {
				return err
			}
			err = processTypes(p, *p.Files[f], *p.Files[f].Types[t], a)
			if err != nil {
				return err
			}
		}
		for fn := range p.Files[f].Functions {
			a, err := annotate.LoadFromComment(strings.Join(p.Files[f].Functions[fn].Docs, "\n"))
			if err != nil {
				return err
			}
			err = processFunction(p, *p.Files[f], *p.Files[f].Functions[fn], a)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Process all plugins against this function
func processFunction(pkg humanize.Package, p humanize.File, f humanize.Function, a annotate.AnnotateGroup) error {
	lock.Lock()
	defer lock.Unlock()
	for _, item := range a {
		for i := range plugins {
			switch plug := plugins[i].p.(type) {
			case AnnotationFunction:
				if inArray(item.Name, plug.GetType()...) && plug.FunctionIsSupported(p, f) {
					c, err := plug.ProcessFunction(
						plugins[i].context,
						pkg,
						p,
						f,
						item,
					)
					if err != nil {
						return err
					}
					plugins[i].context = c
					plugins[i].called = true
				}
			}
		}
	}

	return nil
}

// Process all plugins against this structures
func processStructure(pkg humanize.Package, p humanize.File, f humanize.TypeName, a annotate.AnnotateGroup) error {
	lock.Lock()
	defer lock.Unlock()
	for _, item := range a {
		for i := range plugins {
			switch plug := plugins[i].p.(type) {
			case AnnotationStruct:
				if inArray(item.Name, plug.GetType()...) && plug.StructureIsSupported(p, f) {
					c, err := plug.ProcessStructure(
						plugins[i].context,
						pkg,
						p,
						f,
						item,
					)
					if err != nil {
						return err
					}
					plugins[i].context = c
					plugins[i].called = true
				}
			}
		}
	}

	return nil
}

// Process all plugins against this type
func processTypes(pkg humanize.Package, p humanize.File, f humanize.TypeName, a annotate.AnnotateGroup) error {
	lock.Lock()
	defer lock.Unlock()
	for _, item := range a {
		for i := range plugins {
			switch plug := plugins[i].p.(type) {
			case AnnotationType:
				if inArray(item.Name, plug.GetType()...) && plug.TypeIsSupported(p, f) {
					c, err := plug.ProcessType(
						plugins[i].context,
						pkg,
						p,
						f,
						item,
					)
					if err != nil {
						return err
					}
					plugins[i].context = c
					plugins[i].called = true
				}
			}
		}
	}

	return nil
}

// Finalize all plugins
func Finalize(p humanize.Package) error {
	lock.Lock()
	defer lock.Unlock()

	for i := range plugins {
		if plugins[i].called {
			err := plugins[i].p.Finalize(plugins[i].context, p)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
