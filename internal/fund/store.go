package fund

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// DB defines the database operations required by the store.
type DB interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

// Store implements Repository using a database.
// All methods respect context cancellation and deadlines; callers should
// set appropriate timeouts via context.WithTimeout for production use.
//
// Required indexes (see migrations/006_add_funds_indexes.up.sql):
//   - idx_funds_created_at: (created_at DESC, id DESC) for ordered list queries
//   - idx_funds_name: UNIQUE (name) for duplicate name prevention
type Store struct {
	db DB
}

// NewStore creates a new Store.
// Returns nil if db is nil.
func NewStore(db DB) *Store {
	if db == nil {
		return nil
	}
	return &Store{db: db}
}

// Create persists a new fund to the database.
func (s *Store) Create(ctx context.Context, fund *Fund) error {
	return s.create(ctx, s.db, fund)
}

// CreateTx persists a new fund within the provided transaction.
func (s *Store) CreateTx(ctx context.Context, tx pgx.Tx, fund *Fund) error {
	return s.create(ctx, tx, fund)
}

func (s *Store) create(ctx context.Context, db DB, fund *Fund) error {
	if fund == nil {
		return ErrNilFund
	}

	const query = `
		INSERT INTO funds (id, name, total_units, created_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := db.Exec(ctx, query, fund.ID, fund.Name, fund.TotalUnits, fund.CreatedAt)
	if err != nil {
		// Check for unique constraint violation (PostgreSQL error code 23505)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%w: %s", ErrDuplicateFundName, fund.Name)
		}
		return fmt.Errorf("create fund %s: %w", fund.ID, err)
	}
	return nil
}

// FindByID retrieves a fund by its UUID.
func (s *Store) FindByID(ctx context.Context, id uuid.UUID) (*Fund, error) {
	const query = `
		SELECT id, name, total_units, created_at
		FROM funds
		WHERE id = $1
	`
	var fund Fund
	err := s.db.QueryRow(ctx, query, id).Scan(
		&fund.ID,
		&fund.Name,
		&fund.TotalUnits,
		&fund.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NotFoundError(id)
		}
		return nil, fmt.Errorf("find fund %s: %w", id, err)
	}
	return &fund, nil
}

// List retrieves funds with pagination, ordered by created_at descending.
func (s *Store) List(ctx context.Context, params ListParams) (*ListResult, error) {
	params = params.Normalize()

	// Single query with window function for count and pagination.
	const query = `
		SELECT id, name, total_units, created_at, COUNT(*) OVER() AS total
		FROM funds
		ORDER BY created_at DESC, id DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := s.db.Query(ctx, query, params.Limit, params.Offset)
	if err != nil {
		return nil, fmt.Errorf("list funds: %w", err)
	}
	defer rows.Close()

	funds := make([]*Fund, 0, params.Limit)
	var total int
	for rows.Next() {
		var fund Fund
		if err := rows.Scan(&fund.ID, &fund.Name, &fund.TotalUnits, &fund.CreatedAt, &total); err != nil {
			return nil, fmt.Errorf("scan fund row: %w", err)
		}
		funds = append(funds, &fund)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate fund rows: %w", err)
	}

	// When offset exceeds total rows, window function returns no rows.
	// Fall back to count query to get the actual total.
	if len(funds) == 0 && params.Offset > 0 {
		const countQuery = `SELECT COUNT(*) FROM funds`
		if err := s.db.QueryRow(ctx, countQuery).Scan(&total); err != nil {
			return nil, fmt.Errorf("count funds: %w", err)
		}
	}

	return &ListResult{
		Items:  funds,
		Total:  total,
		Limit:  params.Limit,
		Offset: params.Offset,
	}, nil
}
