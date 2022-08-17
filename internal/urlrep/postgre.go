package urlrep

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	insertURLPostgreQuery = "INSERT INTO urls (id, url) VALUES ($1, $2) ON CONFLICT (url) DO UPDATE SET url = $2 RETURNING id"
)

type postgreURLRepo struct {
	Address        string
	pool           *pgxpool.Pool
	urlIDGenerator func(url string) string
}

// NewPostgre returns new PostgreSQL URL repository.
func NewPostgre(address string, urlIDGenerator func(url string) string) (URLRepo, error) {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, address)
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)

	_, err = conn.Exec(ctx, "CREATE TABLE IF NOT EXISTS urls (id VARCHAR(8) PRIMARY KEY, url VARCHAR(255))")
	if err != nil {
		return nil, err
	}

	_, err = conn.Exec(ctx, "CREATE UNIQUE INDEX IF NOT EXISTS unique_url ON urls (url)")
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.Connect(ctx, address)
	if err != nil {
		return nil, err
	}

	if urlIDGenerator == nil {
		urlIDGenerator = func(url string) string {
			return uuid.New().String()[:8]
		}
	}

	return &postgreURLRepo{
		Address:        address,
		pool:           pool,
		urlIDGenerator: urlIDGenerator,
	}, nil
}

func (r *postgreURLRepo) URLByID(ctx context.Context, id string) (string, bool) {
	var url string
	r.pool.QueryRow(ctx, "SELECT url FROM urls WHERE id = $1", id).Scan(&url)
	return url, url != ""
}

func (r *postgreURLRepo) Save(ctx context.Context, url string) (string, error) {
	var recID string
	id := r.GenerateID(url)
	r.pool.QueryRow(ctx, insertURLPostgreQuery, id, url).Scan(&recID)
	if recID != id {
		return recID, ErrDuplicateURL
	}
	return id, nil
}

func (r *postgreURLRepo) SaveList(ctx context.Context, urls []string) ([]string, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	stmt, err := tx.Prepare(ctx, "insert_url", insertURLPostgreQuery)
	if err != nil {
		return nil, err
	}

	ids := make([]string, len(urls))
	for i, url := range urls {
		id := r.GenerateID(url)
		ids[i] = id
		_, err := tx.Exec(ctx, stmt.Name, id, url)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *postgreURLRepo) GenerateID(url string) string {
	return r.urlIDGenerator(url)
}
