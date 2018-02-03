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
)

const (
	// AssetTypeImageSubTypeIcon image/icon part
	AssetTypeImageSubTypeIcon int = iota + 1
	// AssetTypeImageSubTypeLogo image/logo part
	AssetTypeImageSubTypeLogo
	// AssetTypeImageSubTypeMain main image
	AssetTypeImageSubTypeMain
	// AssetTypeImageSubTypeBanner is a hack to handle banners
	AssetTypeImageSubTypeBanner int = 501
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
	// AssetTypeTextSubTypeRating rating
	AssetTypeTextSubTypeRating
	// AssetTypeTextSubTypeLikes likes
	AssetTypeTextSubTypeLikes
	// AssetTypeTextSubTypeDownloads downloads
	AssetTypeTextSubTypeDownloads
	// AssetTypeTextSubTypePrice price
	AssetTypeTextSubTypePrice
	// AssetTypeTextSubTypeSalePrice sale price
	AssetTypeTextSubTypeSalePrice
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
	// Duration is only for videos
	Duration int `json:"d,omitempty"`
	// Data is the actual data inside this asset
	Data string `json:"data"`
}

// AssetFilter is a callback function to filter asset
type AssetFilter func(*Asset) bool

// Filter is a full filter of one asset
type Filter struct {
	ID       int
	Type     AssetType
	SubType  int
	Required bool
	Extra    []AssetFilter
}
