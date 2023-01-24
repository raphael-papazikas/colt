package colt

import (
	"math/rand"
	"time"
)

func MockSetup() *Database {
	rand.Seed(time.Now().UnixNano())
	var mockDb = Database{}
	mockDb.Connect("mongodb://localhost:27017/colt?readPreference=primary&directConnection=true&ssl=false", "colt")
	mockDb.db.Drop(DefaultContext())
	return &mockDb
}
