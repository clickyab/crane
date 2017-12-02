package models

import (
	"clickyab.com/crane/crane/models/internal/entities"
	"github.com/clickyab/services/pool"
)

var networks pool.Interface

// GetNetworkID return the network object based on the network name
func GetNetworkID(b string) (int64, error) {
	bs := &entities.Network{}
	data, err := networks.Get(b, bs)
	if err != nil {
		return 0, err
	}
	return data.(*entities.Network).ID, nil
}
