package entities

import (
	"context"
	"encoding/gob"
	"fmt"
	"io"

	"github.com/clickyab/services/kv"
)

// Brand model for database
type Brand struct {
	ID    int64  `json:"-" db:"id"`
	Brand string `json:"-" db:"brand"`
}

// Encode is the encode function for serialize object in io writer
func (b *Brand) Encode(w io.Writer) error {
	return gob.NewEncoder(w).Encode(b)
}

// Decode try to decode object from io reader
func (b *Brand) Decode(r io.Reader) error {
	return gob.NewDecoder(r).Decode(b)
}

// BrandLoader for caching brands
func BrandLoader(ctx context.Context) (map[string]kv.Serializable, error) {
	var res []Brand

	q := "SELECT ab_id as id, ab_brand as brand FROM apps_brands where ab_show=1"

	_, err := NewManager().GetRDbMap().Select(
		&res,
		q,
	)
	if err != nil {
		return nil, err
	}
	b := make(map[string]kv.Serializable)
	for i := range res {
		b[fmt.Sprint(res[i].ID)] = &res[i]
	}
	return b, nil
}

// Carrier model for database
type Carrier struct {
	ID      int64  `json:"-" db:"id"`
	Carrier string `json:"-" db:"carrier"`
}

// Encode is the encode function for serialize object in io writer
func (b *Carrier) Encode(w io.Writer) error {
	return gob.NewEncoder(w).Encode(b)
}

// Decode try to decode object from io reader
func (b *Carrier) Decode(r io.Reader) error {
	return gob.NewDecoder(r).Decode(b)
}

// CarrierLoader for caching carriers
func CarrierLoader(ctx context.Context) (map[string]kv.Serializable, error) {
	var res []Carrier

	q := "SELECT ac_id as id, ac_carrier as carrier FROM apps_carriers where ac_show=1"

	_, err := NewManager().GetRDbMap().Select(
		&res,
		q,
	)
	if err != nil {
		return nil, err
	}
	b := make(map[string]kv.Serializable)
	for i := range res {
		b[fmt.Sprint(res[i].ID)] = &res[i]
	}
	return b, nil
}

// Network model for database
type Network struct {
	ID      int64  `json:"-" db:"id"`
	Network string `json:"-" db:"network"`
}

// Encode is the encode function for serialize object in io writer
func (b *Network) Encode(w io.Writer) error {
	return gob.NewEncoder(w).Encode(b)
}

// Decode try to decode object from io reader
func (b *Network) Decode(r io.Reader) error {
	return gob.NewDecoder(r).Decode(b)
}

// NetworkLoader for caching networks
func NetworkLoader(ctx context.Context) (map[string]kv.Serializable, error) {
	var res []Network
	q := "SELECT an_id as id, an_network as network FROM apps_networks where an_show=1"
	_, err := NewManager().GetRDbMap().Select(
		&res,
		q,
	)
	if err != nil {
		return nil, err
	}
	b := make(map[string]kv.Serializable)
	for i := range res {
		b[fmt.Sprint(res[i].ID)] = &res[i]
	}
	return b, nil
}
