package fund

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresRepository implements Repository using PostgreSQL.
type PostgresRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresRepository creates a new PostgresRepository.
func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

// Create persists a new fund to the database.
func (r *PostgresRepository) Create(ctx context.Context, fund *Fund) error {
	return r.create(ctx, r.pool, fund)
}

// CreateTx persists a new fund within the provided transaction.
func (r *PostgresRepository) CreateTx(ctx context.Context, tx pgx.Tx, fund *Fund) error {
	return r.create(ctx, tx, fund)
}

// querier abstracts pgxpool.Pool and pgx.Tx for shared query logic.
type querier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

func (r *PostgresRepository) create(ctx context.Context, q querier, fund *Fund) error {
	const query = `
		INSERT INTO funds (id, name, total_units, created_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := q.Exec(ctx, query, fund.ID, fund.Name, fund.TotalUnits, fund.CreatedAt)
	return err
}

// FindByID retrieves a fund by its UUID.
func (r *PostgresRepository) FindByID(ctx context.Context, id uuid.UUID) (*Fund, error) {
	const query = `
		SELECT id, name, total_units, created_at
		FROM funds
		WHERE id = $1
	`
	var fund Fund
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&fund.ID,
		&fund.Name,
		&fund.TotalUnits,
		&fund.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NotFoundError(id)
		}
		return nil, err
	}
	return &fund, nil
}

// List retrieves funds with pagination, ordered by created_at descending.
func (r *PostgresRepository) List(ctx context.Context, params ListParams) (*ListResult, error) {
	params = params.Normalize()

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM funds`
	if err := r.pool.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, err
	}

	// Get paginated results
	const query = `
		SELECT id, name, total_units, created_at
		FROM funds
		ORDER BY created_at DESC, id DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.pool.Query(ctx, query, params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	funds := make([]*Fund, 0)
	for rows.Next() {
		var fund Fund
		if err := rows.Scan(&fund.ID, &fund.Name, &fund.TotalUnits, &fund.CreatedAt); err != nil {
			return nil, err
		}
		funds = append(funds, &fund)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &ListResult{
		Items:  funds,
		Total:  total,
		Limit:  params.Limit,
		Offset: params.Offset,
	}, nil
}
