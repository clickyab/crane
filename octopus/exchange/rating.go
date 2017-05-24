package exchange

// Rate is an arbitrary number that defines and ad can be send for publisher or not
type Rate string

const (
	// RateA is an arbitrary number
	RateA Rate = "A"
	// RateB is an arbitrary number
	RateB Rate = "B"
	// RateC is an arbitrary number
	RateC Rate = "C"
	// RateD is an arbitrary number
	RateD Rate = "D"
	// RateE is an arbitrary number
	RateE Rate = "E"
	// RateF is an arbitrary number
	RateF Rate = "F"
	// RateG is an arbitrary number
	RateG Rate = "G"
	// RateH is an arbitrary number
	RateH Rate = "H"
	// RateI is an arbitrary number
	RateI Rate = "I"
	// RateJ is an arbitrary number
	RateJ Rate = "J"
)

// Rater is interface
type Rater interface {
	// Rates return  slice of rate that can be represent rate of ad or publisher rate blacklist
	Rates() []Rate
}
