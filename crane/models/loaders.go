package models

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

func (b *Brand) Decode(w io.Writer) error {
	return gob.NewEncoder(w).Encode(b)
}

func (b *Brand) Encode(r io.Reader) error {
	return gob.NewDecoder(r).Decode(b)
}

// BrandLoader for caching brands
func BrandLoader(ctx context.Context) (map[string]kv.Serializable, error) {
	var res []Brand

	q := "SELECT ab_id as id, ab_brand as brand FROM apps_brands"

	_, err := NewManager().GetRDbMap().Select(
		&res,
		q,
	)
	if err != nil {
		return nil, err
	}
	b := make(map[string]kv.Serializable, 0)
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

func (b *Carrier) Decode(w io.Writer) error {
	return gob.NewEncoder(w).Encode(b)
}

func (b *Carrier) Encode(r io.Reader) error {
	return gob.NewDecoder(r).Decode(b)
}

// CarrierLoader for caching carriers
func CarrierLoader(ctx context.Context) (map[string]kv.Serializable, error) {
	var res []Carrier

	q := "SELECT ac_id as id, ac_carrier as carrier FROM apps_carriers"

	_, err := NewManager().GetRDbMap().Select(
		&res,
		q,
	)
	if err != nil {
		return nil, err
	}
	b := make(map[string]kv.Serializable, 0)
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

func (b *Network) Decode(w io.Writer) error {
	return gob.NewEncoder(w).Encode(b)
}

func (b *Network) Encode(r io.Reader) error {
	return gob.NewDecoder(r).Decode(b)
}

// NetworkLoader for caching networks
func NetworkLoader(ctx context.Context) (map[string]kv.Serializable, error) {
	var res []Network
	q := "SELECT an_id as id, an_network as network FROM apps_networks"
	_, err := NewManager().GetRDbMap().Select(
		&res,
		q,
	)
	if err != nil {
		return nil, err
	}
	b := make(map[string]kv.Serializable, 0)
	for i := range res {
		b[fmt.Sprint(res[i].ID)] = &res[i]
	}
	return b, nil
}
