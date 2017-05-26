package mocks

type Slot struct {
	SWidth, SHeight int
	STRackID        string
	SFallback       string
}

func (s Slot) Width() int {
	return s.SWidth
}

func (s Slot) Height() int {
	return s.SHeight
}

func (s Slot) TrackID() string {
	return s.STRackID
}

func (s Slot) Fallback() string {
	return s.SFallback
}
