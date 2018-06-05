package entities

import (
	"encoding/gob"
	"io"

	"clickyab.com/crane/demand/entity"
)

//CreativeStatistics struct type for add network creatives statistics
type CreativeStatistics struct {
	AdType int64   `db:"ad_type" json:"ad_type"`
	Count  int64   `db:"count" json:"count"`
	Minctr float64 `db:"min_ctr" json:"min_ctr"`
	Maxctr float64 `db:"max_ctr" json:"max_ctr"`
	Avgctr float64 `db:"avg_ctr" json:"avg_ctr"`
}

// CreativeType return the type of ad
func (c *CreativeStatistics) CreativeType() entity.AdType {
	return entity.AdType(c.AdType)
}

// TotalCount return total count of creatives base on type
func (c *CreativeStatistics) TotalCount() int64 {
	return c.Count
}

// MinCTR return network creatives statistics
func (c *CreativeStatistics) MinCTR() float64 {
	return c.Minctr
}

// MaxCTR return network creatives statistics
func (c *CreativeStatistics) MaxCTR() float64 {
	return c.Maxctr
}

// AvgCTR return network creatives statistics
func (c *CreativeStatistics) AvgCTR() float64 {
	return c.Avgctr
}

// Encode is the encode function for serialize object in io writer
func (c *CreativeStatistics) Encode(w io.Writer) error {
	g := gob.NewEncoder(w)
	return g.Encode(c)
}

// Decode try to decode object from io reader
func (c *CreativeStatistics) Decode(r io.Reader) error {
	g := gob.NewDecoder(r)
	return g.Decode(c)
}
