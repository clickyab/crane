package entity

// Target is the target of this campaign
type Target int

const (
	// TargetWeb is the normal banner in system
	TargetWeb Target = iota
	// TargetApp is the app targeted campaign
	TargetApp
	// TargetVast is the vast target
	TargetVast
)

// Campaign is the single campaign in ssytem
type Campaign interface {
	// ID return the campaign id
	ID() int64
	// Name is the campaign name
	Name() string
	// MaxBID get the campaign max bid
	MaxBID() int64
	// Make sure the result is >= 1
	Frequency() int
	// Target return the target of this campaign
	Target() []Target
}
