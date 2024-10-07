package service

import (
	"errors"
	"math/rand"
)

// DataProviderImpl is the concrete implementation of the DataProvider interface.
type DataProviderImpl struct{}

// GetRandomNumber simulates fetching a random number between 0 and id.
func (d *DataProviderImpl) GetRandomNumber(id int) (int, error) {
	if id < 0 {
		return 0, errors.New("InvalidId")
	}
	// Simulate fetching a random number between 0 and id
	return rand.Intn(id + 1), nil
}
