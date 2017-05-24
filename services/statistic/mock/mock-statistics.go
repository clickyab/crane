package mock

import (
	"fmt"
	"strconv"
	"time"

	"clickyab.com/exchange/services/statistic"
)

var (
	pre   = make(map[string]map[string]string)
	store = make(map[string]*Interface)
)

// Interface Interface
type Interface struct {
	MasterKey string
	Data      map[string]string
	Duration  time.Duration
}

// Key return key
func (i *Interface) Key() string {
	return i.MasterKey
}

// IncSubKey increment subKey
func (i *Interface) IncSubKey(s string, a int64) (int64, error) {
	res, ok := i.Data[s]
	if !ok {
		i.Data[s] = fmt.Sprintf("%d", a)
		return a, nil
	}
	oldNum, err := strconv.ParseInt(res, 10, 0)
	if err != nil {
		return 0, err
	}
	return oldNum + a, nil
}

// DecSubKey decrement subKey
func (i *Interface) DecSubKey(s string, a int64) (int64, error) {
	res, ok := i.Data[s]
	if !ok {
		i.Data[s] = fmt.Sprintf("%d", -a)
		return -a, nil
	}
	oldNum, err := strconv.ParseInt(res, 10, 0)
	if err != nil {
		return 0, err
	}
	return oldNum - a, nil
}

// Touch return subKey
func (i *Interface) Touch(s string) (int64, error) {
	res, ok := i.Data[s]
	if !ok {
		i.Data[s] = fmt.Sprintf("%d", 0)
		return 0, nil
	}
	final, err := strconv.ParseInt(res, 10, 0)
	if err != nil {
		return 0, err
	}
	return final, nil
}

// GetAll get all
func (i *Interface) GetAll() (map[string]int64, error) {
	res := i.Data
	final := make(map[string]int64)
	for j := range res {
		num, err := strconv.ParseInt(res[j], 10, 0)
		if err != nil {
			return make(map[string]int64), err
		}
		final[j] = num
	}
	return final, nil
}

// NewMockStatistic generate new mock
func NewMockStatistic(s string, t time.Duration) statistic.Interface {
	if k, ok := store[s]; ok {
		return k
	}
	var (
		data map[string]string
		ok   bool
	)
	if data, ok = pre[s]; !ok {
		data = make(map[string]string)
	}
	m := &Interface{
		MasterKey: s,
		Data:      data,
		Duration:  t,
	}

	store[s] = m
	return m
}
