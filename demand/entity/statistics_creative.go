package entity

//CreativeStatistics struct type for add network creatives statistics
type CreativeStatistics struct {
	AdType int64   `db:"ad_type" json:"ad_type"`
	Count  int64   `db:"count" json:"count"`
	MinCTR float64 `db:"min_ctr" json:"min_ctr"`
	MaxCTR float64 `db:"max_ctr" json:"max_ctr"`
	AvgCTR float64 `db:"avg_ctr" json:"avg_ctr"`
}
