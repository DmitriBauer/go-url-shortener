package urlrep

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	insertURLPostgreQuery = "INSERT INTO urls (id, session_id, url) VALUES ($1, $2, $3) ON CONFLICT (url) DO UPDATE SET url = $3 RETURNING id"
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

	_, err = conn.Exec(ctx, "CREATE TABLE IF NOT EXISTS urls (id VARCHAR(8) PRIMARY KEY, session_id VARCHAR(36), url VARCHAR(255), removed BOOLEAN DEFAULT FALSE)")
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
	var removed bool
	r.pool.QueryRow(ctx, "SELECT url, removed FROM urls WHERE id = $1", id).Scan(&url, &removed)
	return url, removed
}

func (r *postgreURLRepo) Save(ctx context.Context, url string, sessionID string) (string, error) {
	var recID string
	id := r.GenerateID(url)
	r.pool.QueryRow(ctx, insertURLPostgreQuery, id, sessionID, url).Scan(&recID)
	if recID != id {
		return recID, ErrDuplicateURL
	}
	return id, nil
}

func (r *postgreURLRepo) SaveList(ctx context.Context, urls []string, sessionID string) ([]string, error) {
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
		_, err := tx.Exec(ctx, stmt.Name, id, sessionID, url)
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

func (r *postgreURLRepo) RemoveList(ctx context.Context, ids []string, sessionID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	stmt, err := tx.Prepare(ctx, "set_url_removed", "UPDATE urls SET removed = TRUE WHERE id = $1 AND session_id = $2")
	if err != nil {
		return err
	}

	for _, id := range ids {
		_, err := tx.Exec(ctx, stmt.Name, id, sessionID)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *postgreURLRepo) GenerateID(url string) string {
	return r.urlIDGenerator(url)
}
