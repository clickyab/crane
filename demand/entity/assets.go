package entity

type (
	// AssetType is the main asset types
	AssetType int
)

const (
	// AssetTypeImage is all image asset type
	AssetTypeImage AssetType = iota
	// AssetTypeVideo is all video asset type
	AssetTypeVideo
	// AssetTypeText is all text type (from data normally)
	AssetTypeText
	// AssetTypeNumber is all number from data field
	AssetTypeNumber
)

const (
	// AssetTypeImageSubTypeIcon image/icon part
	AssetTypeImageSubTypeIcon int = iota + 1
	// AssetTypeImageSubTypeLogo image/logo part
	AssetTypeImageSubTypeLogo
	// AssetTypeImageSubTypeMain main image
	AssetTypeImageSubTypeMain
	// AssetTypeImageSubTypeBanner is a hack to handle banners
	AssetTypeImageSubTypeBanner int = 500
)

const (
	// AssetTypeVideoSubTypeMain is the video only type
	AssetTypeVideoSubTypeMain = iota + 1
)

const (
	// AssetTypeTextSubTypeTitle title type, normally in a separate request (ortb spec)
	AssetTypeTextSubTypeTitle int = iota
	// AssetTypeTextSubTypeSponsored sponsored text
	AssetTypeTextSubTypeSponsored
	// AssetTypeTextSubTypeDesc main desceiption
	AssetTypeTextSubTypeDesc
	// AssetTypeNumberSubTypeRating rating
	AssetTypeNumberSubTypeRating
	// AssetTypeNumberSubTypeLikes likes
	AssetTypeNumberSubTypeLikes
	// AssetTypeNumberSubTypeDownloads downloads
	AssetTypeNumberSubTypeDownloads
	// AssetTypeNumberSubTypePrice price
	AssetTypeNumberSubTypePrice
	// AssetTypeNumberSubTypeSalePrice sale price
	AssetTypeNumberSubTypeSalePrice
	// AssetTypeTextSubTypePhone phone
	AssetTypeTextSubTypePhone
	// AssetTypeTextSubTypeAddress address
	AssetTypeTextSubTypeAddress
	// AssetTypeTextSubTypeDesc2 2nd description
	AssetTypeTextSubTypeDesc2
	// AssetTypeTextSubTypeDisplayURL display url
	AssetTypeTextSubTypeDisplayURL
	// AssetTypeTextSubTypeCTAText Call to action text
	AssetTypeTextSubTypeCTAText
)

// Asset is a structure to handle single asset inside creative
type Asset struct {
	Type    AssetType `json:"t"`
	SubType int       `json:"st"`
	// if asset support width and height
	Width  int `json:"w,omitempty"`
	Height int `json:"h,omitempty"`
	// if asset support len
	Len int `json:"l,omitempty"`
	// MimeType is the mime-type checker
	MimeType string `json:"mt"`
	// Data is the actual data inside this asset
	Data string `json:"data"`
}

// AssetFilter is a callback function to filter asset
type AssetFilter func(*Asset) bool

// MinWidth return a minimum width checker
func MinWidth(m int) AssetFilter {
	return func(a *Asset) bool {
		if a.Type != AssetTypeImage {
			return false
		}
		return a.Width >= m
	}
}

// MinHeight return a minimum height checker
func MinHeight(m int) AssetFilter {
	return func(a *Asset) bool {
		if a.Type != AssetTypeImage {
			return false
		}
		return a.Height >= m
	}
}

// ExactWidth return an exact width checker
func ExactWidth(m int) AssetFilter {
	return func(a *Asset) bool {
		if a.Type != AssetTypeImage {
			return false
		}
		return a.Width == m
	}
}

// ExactHeight return an exact height checker
func ExactHeight(m int) AssetFilter {
	return func(a *Asset) bool {
		if a.Type != AssetTypeImage {
			return false
		}
		return a.Height == m
	}
}

// MaxLen return a max length checker (string)
func MaxLen(m int) AssetFilter {
	return func(a *Asset) bool {
		if a.Type != AssetTypeText {
			return false
		}
		return a.Len <= m
	}
}

// ContainMimeType check for mime type in creative
func ContainMimeType(t ...string) AssetFilter {
	return func(a *Asset) bool {
		for i := range t {
			if a.MimeType == t[i] {
				return true
			}
		}
		return false
	}
}
