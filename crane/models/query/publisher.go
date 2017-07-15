package query

import (
	"clickyab.com/crane/crane/entity"
	"github.com/clickyab/services/assert"
)

var qp entity.QPublisher

// Publisher return queryable object
func Publisher() entity.QPublisher {
	assert.NotNil(qp)
	return qp
}

// Register will save queryable object
func Register(n entity.QPublisher) {
	qp = n
}
