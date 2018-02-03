package filters

import "clickyab.com/crane/demand/entity"

// MinWidth return a minimum width checker
func MinWidth(m int) entity.AssetFilter {
	return func(a *entity.Asset) bool {
		if a.Type != entity.AssetTypeImage {
			return false
		}
		return a.Width >= m
	}
}

// MinHeight return a minimum height checker
func MinHeight(m int) entity.AssetFilter {
	return func(a *entity.Asset) bool {
		if a.Type != entity.AssetTypeImage {
			return false
		}
		return a.Height >= m
	}
}

// ExactWidth return an exact width checker
func ExactWidth(m int) entity.AssetFilter {
	return func(a *entity.Asset) bool {
		if a.Type != entity.AssetTypeImage {
			return false
		}
		return a.Width == m
	}
}

// ExactHeight return an exact height checker
func ExactHeight(m int) entity.AssetFilter {
	return func(a *entity.Asset) bool {
		if a.Type != entity.AssetTypeImage {
			return false
		}
		return a.Height == m
	}
}

// MaxLen return a max length checker (string)
func MaxLen(m int) entity.AssetFilter {
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
func MaxDuration(d int) entity.AssetFilter {
	return func(a *entity.Asset) bool {
		if a.Type != entity.AssetTypeVideo {
			return false
		}
		return a.Duration <= d

	}
}

// MinDuration for videos
func MinDuration(d int) entity.AssetFilter {
	return func(a *entity.Asset) bool {
		if a.Type != entity.AssetTypeVideo {
			return false
		}
		return a.Duration >= d

	}
}

// VastProtocol check for vast 2/3
func VastProtocol(d ...int) entity.AssetFilter {
	return func(a *entity.Asset) bool {
		if a.Type != entity.AssetTypeVideo {
			return false
		}
		for i := range d {
			if d[i] == 3 {
				return true
			}
		}
		return false
	}
}
