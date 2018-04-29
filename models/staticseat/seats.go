package staticseat

import (
	"clickyab.com/crane/demand/entity"
	entities2 "clickyab.com/crane/models/internal/entities"
	"clickyab.com/crane/supplier/layers/entities"
	"github.com/clickyab/services/pool"
)

var seats pool.Interface

// GetStaticSeats return all ads in system
func GetStaticSeats() []entities.StaticSeat {
	data := seats.All()
	all := make([]entities.StaticSeat, len(data))
	var c int
	for i := range data {
		all[c] = data[i].(entities.StaticSeat)
		c++
	}

	return all
}

// GetStaticSeat try to get website. do not use it in initializer
func GetStaticSeat(pub entity.Publisher, typ, position string) (entities.StaticSeat, error) {
	d := &entities2.StaticSeat{}
	res, err := seats.Get(pub.Name()+"/"+pub.Supplier().Name()+"/"+typ+"/"+position, d)
	if err != nil {
		return nil, err
	}
	d = res.(*entities2.StaticSeat)
	return d, nil
}
