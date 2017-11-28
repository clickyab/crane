package cyslot

import (
	"fmt"

	"strings"
)

var sizes = map[string]int{
	"120x600":  1,
	"160x600":  2,
	"300x250":  3,
	"336x280":  4,
	"468x60":   5,
	"728x90":   6,
	"120x240":  7,
	"320x50":   8,
	"800x440":  9,
	"300x600":  11,
	"970x90":   12,
	"970x250":  13,
	"250x250":  14,
	"300x1050": 15,
	"320x480":  16,
	"48x320":   17,
	"128x128":  18,
}

// GetSize return the size of a banner in clickyab std
func GetSize(size string) (int, error) {
	s, ok := sizes[size]
	if ok {
		return s, nil
	}

	return 0, fmt.Errorf("size %s is not valid", size)
}

func ValidWebSlotSize(size int) bool {
	for i := range sizes {
		if sizes[i] == size {
			return true
		}
	}
	return false
}

// GetSizeByNum return the size
func GetSizeByNum(num int) (string, string) {
	for i := range sizes {
		if sizes[i] == num {
			if dimentions := strings.Split(i, "x"); len(dimentions) == 2 {
				return dimentions[0], dimentions[1]
			}
			return "", ""
		}
	}
	return "", ""
}
