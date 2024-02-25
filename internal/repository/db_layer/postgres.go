package dblayer

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/platatest/internal/repository"
	"github.com/platatest/pkg/repository/connection/postgres"
)

type Postgres struct {
	con *pgx.Conn
}

func NewPostgres(url string) (repository.DatabaseHandler, error) {
	con, err := postgres.New(url)
	if err != nil {
		return nil, err
	}
	return &Postgres{con: con}, err
}

func (p *Postgres) Create(name string) (int, error) {
	tx, err := p.con.Begin(context.Background())
	if err != nil {
		return -1, fmt.Errorf("%w : %w", repository.TxErr, err)
	}
	defer tx.Rollback(context.Background())

	query := `insert into  currency (name)
	values ($1)
	on conflict (name) 
	do update set name = excluded.name
	returning id;`

	var id int
	err = tx.QueryRow(context.Background(), query, name).Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("%w %s : %w", repository.CreateErr, name, err)
	}

	priceQuery := `insert into price (currency_id) 
	values ($1)
	on conflict (currency_id)
	do update set  currency_id = excluded.currency_id;`
	_, err = tx.Exec(context.Background(), priceQuery, id)
	if err != nil {
		return -1, fmt.Errorf("%w %s : %w", repository.CreateErr, name, err)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return -1, fmt.Errorf("%w %s : %w", repository.TxCommitErr, name, err)
	}
	return id, err

}
func (p *Postgres) GetById(id int) (repository.Price, error) {

	query := `select price, date from price where currency_id = $1`

	var cur repository.Price
	err := p.con.QueryRow(context.Background(), query, id).Scan(&cur.Value, &cur.Date)
	if err != nil {
		return cur, fmt.Errorf("%w %d : %w", repository.GetByIdErr, id, err)
	}

	return cur, err
}

func (p *Postgres) GetByName(name string) (repository.Price, error) {

	query := `select price, date from price
	join public.currency c on c.id = price.currency_id
	where name = $1`

	var cur repository.Price
	err := p.con.QueryRow(context.Background(), query, name).Scan(&cur.Value, &cur.Date)
	if err != nil {
		return cur, fmt.Errorf("%w %s : %w", repository.GetByIdErr, name, err)
	}

	return cur, err
}
func (p *Postgres) Fetch() (repository.Currencies, error) {

	var curs repository.Currencies
	query := `select id, name from currency`

	rows, err := p.con.Query(context.Background(), query)
	if err != nil {
		return curs, fmt.Errorf("%w : %w", repository.FetchErr, err)
	}

	defer rows.Close()

	for rows.Next() {
		var temp repository.Currency
		err = rows.Scan(&temp.Id, &temp.Name)
		if err != nil {
			return curs, fmt.Errorf("%w : %w", repository.ScanErr, err)
		}
		curs = append(curs, temp)
	}

	return curs, err
}
func (p *Postgres) Update(price float64, id int) error {

	query := `update price set price = $1, date = $2 where currency_id = $3`
	_, err := p.con.Exec(context.Background(), query, price, time.Now(), id)
	if err != nil {
		return fmt.Errorf("%w : %w", repository.UpdateErr, err)
	}

	return err
}
func (p *Postgres) Close(ctx context.Context) {
	p.con.Close(ctx)
}
