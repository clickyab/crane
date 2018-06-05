package statistics

import (
	"clickyab.com/crane/demand/entity"
	"github.com/clickyab/services/pool"
)

var creativesStatisticsPool pool.Interface

//GetNetworkCreativesStatistics get total statistics of all creatives in network per type
func GetNetworkCreativesStatistics() []entity.CreativeStatistics {
	data := creativesStatisticsPool.All()
	all := make([]entity.CreativeStatistics, len(data))

	var c int
	for i := range data {
		all[c] = data[i].(entity.CreativeStatistics)
		c++
	}

	return all
}

// GetTypeStatistics try to get network statistics based on its creative type
func GetTypeStatistics(crType entity.AdType) entity.CreativeStatistics {
	data := creativesStatisticsPool.All()

	return data[string(crType)].(entity.CreativeStatistics)
}
