package entity

// MinimalSupplier is used in supplier side
// TODO : use this supplier in supplier side
type MinimalSupplier interface {
	// Name of Supplier
	Name() string
	// AllowCreate indicated if this supplier can create publisher on demand
	AllowCreate() bool
	// DefaultMinBid return the default min bid for this supplier
	DefaultMinBid() int64
}

// Supplier is the ad-network interface
type Supplier interface {
	MinimalSupplier
	// Token of this for web request
	Token() string
	// SoftFloorCPM for this supplier by request sub type and publisher type
	// example : web banner,web vast, app native , ...
	SoftFloorCPM(string, string) int64
	// SoftFloorCPC for this supplier by request sub type and publisher type
	// example : web banner,web vast, app native , ...
	SoftFloorCPC(string, string) int64
	// DefaultCTR for this supplier by request sub type and publisher type
	// example : web banner,web vast, app native , ...
	DefaultCTR(string, string) float64
	// TinyMark means we can add our mark to it
	TinyMark() bool
	// TinyLogo will be the url to the logo (ex: //clickyab.com/tiny.png)
	TinyLogo() string
	// TinyURL is the link of ancher tag of tiny (ex: http://clickyab.com/?ref=tiny)
	TinyURL() string
	// ShowDomain is a domain that all links are generated against it
	ShowDomain() string
	// UserID return user id of supplier
	UserID() int64
	// Rate return ratio currency conversion to IRR
	Rate() int
	// Share is a percentage of the minbid reported to the rtb module. if share is 100 means the
	// min bid reported with no change, if less than 100, means reported less than actual value (not correct normally!)
	// greater than 100 means reported more than its actual value
	Share() int
	// Strategy return default strategy of supplier (cpm, cpc)
	Strategy() Strategy
}
