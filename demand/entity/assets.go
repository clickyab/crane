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
	Type     AssetType `json:"t"`
	SubType  int       `json:"st"`
	Width    int32     `json:"w,omitempty"` // if asset support width and height
	Height   int32     `json:"h,omitempty"`
	Len      int32     `json:"l,omitempty"` // if asset support len
	MimeType string    `json:"mt"`          // MimeType is the mime-type checker
	Duration int32     `json:"d,omitempty"` // Duration is only for videos
	Data     string    `json:"data"`        // Data is the actual data inside this asset

}

// AssetFilter is a callback function to filter asset
type AssetFilter func(*Asset) bool

// Filter is a full filter of one asset
type Filter struct {
	ID       int32
	Type     AssetType
	SubType  int
	Required bool
	Extra    []AssetFilter
}
