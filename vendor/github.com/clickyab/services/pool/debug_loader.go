package pool

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"reflect"

	"github.com/clickyab/services/kv"
	"github.com/clickyab/services/xlog"
)

// Pattern is a helper interface to handle loading data for loader
type Pattern interface {
	Value() kv.Serializable
	Key() string
}

// DebugLoaderGenerator is a simple loader generator that append data to real loader,
// useful for debugging
func DebugLoaderGenerator(l Loader, file string, pattern Pattern) Loader {
	return func(c context.Context) (map[string]kv.Serializable, error) {
		res, err := l(c)
		if err != nil {
			xlog.GetWithError(c, err).Error("top level err")
			return nil, err
		}
		f, err := os.Open(file)
		if err != nil {
			xlog.GetWithError(c, err).Errorf("the file %s is not available", file)
			return nil, err
		}
		defer func() {
			if err := f.Close(); err != nil {
				xlog.GetWithError(c, err).Error("can not close file")
			}
		}()
		decoder := json.NewDecoder(f)
		for {
			// a copy of pattern
			cp := reflect.New(reflect.TypeOf(pattern)).Elem().Addr().Interface()
			err := decoder.Decode(cp)
			if err == io.EOF {
				break
			}
			if err != nil {
				xlog.GetWithError(c, err).Error("decode err")
				return nil, err
			}
			res[cp.(Pattern).Key()] = cp.(Pattern).Value()
		}

		return res, nil
	}
}
