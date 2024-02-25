package factory

import (
	"errors"

	"github.com/platatest/internal/repository"
	dblayer "github.com/platatest/internal/repository/db_layer"
)

type DBTYPE string

const (
	POSTGRE = "POSTGRE"
)

func NewPersistenceLayer(url string, options DBTYPE) (repository.DatabaseHandler, error) {
	switch options {
	case POSTGRE:
		return dblayer.NewPostgres(url)
	}
	return nil, errors.New("No such database type requared")
}
