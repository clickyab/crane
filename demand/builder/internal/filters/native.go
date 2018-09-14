package filters

import (
	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/openrtb"
)

// MinWidth return a minimum width checker
func MinWidth(m int32) entity.AssetFilter {
	return func(a *entity.Asset) bool {
		if a.Type != entity.AssetTypeImage {
			return false
		}
		return a.Width >= m
	}
}

// MinHeight return a minimum height checker
func MinHeight(m int32) entity.AssetFilter {
	return func(a *entity.Asset) bool {
		if a.Type != entity.AssetTypeImage {
			return false
		}
		return a.Height >= m
	}
}

// ExactWidth return an exact width checker
func ExactWidth(m int32) entity.AssetFilter {
	return func(a *entity.Asset) bool {
		if a.Type != entity.AssetTypeImage {
			return false
		}
		return a.Width == m
	}
}

// ExactHeight return an exact height checker
func ExactHeight(m int32) entity.AssetFilter {
	return func(a *entity.Asset) bool {
		if a.Type != entity.AssetTypeImage {
			return false
		}
		return a.Height == m
	}
}

// MaxLen return a max length checker (string)
func MaxLen(m int32) entity.AssetFilter {
	return func(a *entity.Asset) bool {
		if a.Type != entity.AssetTypeText {
			return false
		}
		return a.Len <= m
	}
}

// ContainMimeType check for mime type in creative
func ContainMimeType(t ...string) entity.AssetFilter {
	return func(a *entity.Asset) bool {
		for i := range t {
			if a.MimeType == t[i] {
				return true
			}
		}
		return false
	}
}

// MaxDuration for videos
func MaxDuration(d int32) entity.AssetFilter {
	return func(a *entity.Asset) bool {
		if a.Type != entity.AssetTypeVideo {
			return false
		}
		return a.Duration <= d

	}
}

// MinDuration for videos
func MinDuration(d int32) entity.AssetFilter {
	return func(a *entity.Asset) bool {
		if a.Type != entity.AssetTypeVideo {
			return false
		}
		return a.Duration >= d

	}
}

// VastProtocol check for vast 2/3
func VastProtocol(d ...openrtb.Protocol) entity.AssetFilter {
	return func(a *entity.Asset) bool {
		if a.Type != entity.AssetTypeVideo {
			return false
		}
		for i := range d {
			if d[i] == openrtb.Protocol_VASTX3X0 {
				return true
			}
		}
		return false
	}
}
