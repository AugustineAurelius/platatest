package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/platatest/pkg/repository/connection"
)

func New(url string) (*pgx.Conn, error) {
	con, err := pgx.Connect(context.Background(), url)
	if err != nil {
		return nil, fmt.Errorf("%w : %w", connection.PostgreConErr, err)
	}
	return con, err
}
