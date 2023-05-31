package controller

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/unrelation"
)

func NewOfficeDatabase(mgo *unrelation.Mongo) OfficeDatabase {
	return &officeDatabase{mgo: mgo}
}

type OfficeDatabase interface {
	// table.unrelation.office.go
	// unrelation.office.go
}

type officeDatabase struct {
	mgo *unrelation.Mongo
}
