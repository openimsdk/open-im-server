package service

type DataProvider interface {
	GetRandomNumber(id int) (int, error)
}
