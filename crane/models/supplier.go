package models

import (
	"clickyab.com/crane/crane/entity"
	"clickyab.com/crane/crane/models/internal/entities"
	"github.com/clickyab/services/pool"
)

var (
	suppliers       pool.Interface
	suppliersByName pool.Interface
)

// GetSupplierByToken is a function to get the supplier by its key
func GetSupplierByToken(key string) (entity.Supplier, error) {
	d := &entities.Supplier{}
	res, err := suppliers.Get(key, d)
	if err != nil {
		return nil, err
	}

	return res.(entity.Supplier), nil
}

// GetSupplierByName is a function to get the supplier by its nbame
func GetSupplierByName(name string) (entity.Supplier, error) {
	d := &entities.Supplier{}
	res, err := suppliersByName.Get(name, d)
	if err != nil {
		return nil, err
	}

	return res.(entity.Supplier), nil
}
