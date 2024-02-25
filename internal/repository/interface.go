package repository

import "context"

type DatabaseHandler interface {
	Create(name string) (int, error)
	GetById(id int) (Price, error)
	GetByName(name string) (Price, error)
	Fetch() (Currencies, error)
	Update(price float64, id int) error
	Close(ctx context.Context)
}
