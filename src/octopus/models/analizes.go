package models

// SupSrcDem supplier_source_demand
type SupplierSourceDemand struct {
	Supplier   string `json:"supplier" db:"supplier"`
	Demand     string `json:"demand" db:"demand"`
	Source     string `json:"source" db:"source"`
	Time       int    `json:"time" db:"time"`
	Request    int    `json:"request" db:"request"`
	Impression int    `json:"impression" db:"impression"`
	Show       int    `json:"show" db:"show"`
	ImpBid     int    `json:"imp_bid" db:"imp_bid"`
	ShowBid    int    `json:"show_bid" db:"show_bid"`
	Win        int    `json:"win" db:"win"`
}

// SupSrc supplier_source
type SupplierSource struct {
	Supplier   string `json:"supplier" db:"supplier"`
	Source     string `json:"source" db:"source"`
	Time       int    `json:"time" db:"time"`
	Request    int    `json:"request" db:"request"`
	Impression int    `json:"impression" db:"impression"`
	Show       int    `json:"show" db:"show"`
	ImpBid     int    `json:"imp_bid" db:"imp_bid"`
	ShowBid    int    `json:"show_bid" db:"show_bid"`
}

// DemSrc demand_source
type DemandSource struct {
	Demand     string `json:"demand" db:"demand"`
	Source     string `json:"source" db:"source"`
	Time       int    `json:"time" db:"time"`
	Request    int    `json:"request" db:"request"`
	Impression int    `json:"impression" db:"impression"`
	Show       int    `json:"show" db:"show"`
	ImpBid     int    `json:"imp_bid" db:"imp_bid"`
	ShowBid    int    `json:"show_bid" db:"show_bid"`
}
