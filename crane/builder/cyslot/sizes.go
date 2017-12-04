package cyslot

import (
	"fmt"
)

type size struct {
	Width,
	Height int
}

var (
	// Sizes contain all allowed size
	sizesModel = map[int]*size{
		1:  {Width: 120, Height: 600},
		2:  {Width: 160, Height: 600},
		3:  {Width: 300, Height: 250},
		4:  {Width: 336, Height: 280},
		5:  {Width: 468, Height: 60},
		6:  {Width: 728, Height: 90},
		7:  {Width: 120, Height: 240},
		8:  {Width: 320, Height: 50},
		9:  {Width: 800, Height: 440},
		11: {Width: 300, Height: 600},
		12: {Width: 970, Height: 90},
		13: {Width: 970, Height: 250},
		14: {Width: 250, Height: 250},
		15: {Width: 300, Height: 1050},
		16: {Width: 320, Height: 480},
		17: {Width: 48, Height: 320},
		18: {Width: 128, Height: 128},
	}
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

// ValidWebSlotSize return true if the size is valid
func ValidWebSlotSize(s int) bool {
	if _, ok := sizesModel[s]; ok {
		return true
	}
	return false
}

// GetSizeByNum return the size (order: width, height)
func GetSizeByNum(num int) (int, int) {
	if v, ok := sizesModel[num]; ok {
		return v.Width, v.Height
	}
	panic("not valid size")

}
