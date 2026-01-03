package ownership

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresRepository implements Repository using PostgreSQL.
// All methods respect context cancellation and deadlines; callers should
// set appropriate timeouts via context.WithTimeout for production use.
type PostgresRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresRepository creates a new PostgresRepository.
// Returns nil if pool is nil.
func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	if pool == nil {
		return nil
	}
	return &PostgresRepository{pool: pool}
}

// querier abstracts pgxpool.Pool and pgx.Tx for shared query logic.
type querier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

// Create persists a new cap table entry to the database.
func (r *PostgresRepository) Create(ctx context.Context, entry *Entry) error {
	return r.create(ctx, r.pool, entry)
}

// CreateTx persists a new cap table entry within the provided transaction.
func (r *PostgresRepository) CreateTx(ctx context.Context, tx pgx.Tx, entry *Entry) error {
	return r.create(ctx, tx, entry)
}

func (r *PostgresRepository) create(ctx context.Context, q querier, entry *Entry) error {
	if entry == nil {
		return ErrNilEntry
	}

	const query = `
		INSERT INTO cap_table_entries (id, fund_id, owner_name, units, acquired_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := q.Exec(ctx, query, entry.ID, entry.FundID, entry.OwnerName, entry.Units, entry.AcquiredAt, entry.UpdatedAt, entry.DeletedAt)
	if err != nil {
		// Check for unique constraint violation (PostgreSQL error code 23505).
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("owner %q already exists in fund %s: %w", entry.OwnerName, entry.FundID, err)
		}
		return fmt.Errorf("create cap table entry %s: %w", entry.ID, err)
	}
	return nil
}

// FindByFundID retrieves all cap table entries for a fund with pagination.
func (r *PostgresRepository) FindByFundID(ctx context.Context, fundID uuid.UUID, params ListParams) (*CapTableView, error) {
	params = params.Normalize()

	// Single query with window function for count and pagination.
	// Ordered by units descending to show largest owners first.
	const query = `
		SELECT id, fund_id, owner_name, units, acquired_at, updated_at, deleted_at, COUNT(*) OVER() AS total
		FROM cap_table_entries
		WHERE fund_id = $1 AND deleted_at IS NULL
		ORDER BY units DESC, owner_name ASC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.pool.Query(ctx, query, fundID, params.Limit, params.Offset)
	if err != nil {
		return nil, fmt.Errorf("find cap table entries for fund %s: %w", fundID, err)
	}
	defer rows.Close()

	entries := make([]*Entry, 0, params.Limit)
	var total int
	for rows.Next() {
		var entry Entry
		if err := rows.Scan(&entry.ID, &entry.FundID, &entry.OwnerName, &entry.Units, &entry.AcquiredAt, &entry.UpdatedAt, &entry.DeletedAt, &total); err != nil {
			return nil, fmt.Errorf("scan cap table entry row: %w", err)
		}
		entries = append(entries, &entry)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate cap table entry rows: %w", err)
	}

	// When offset exceeds total rows, window function returns no rows.
	// Fall back to count query to get the actual total.
	if len(entries) == 0 && params.Offset > 0 {
		const countQuery = `SELECT COUNT(*) FROM cap_table_entries WHERE fund_id = $1 AND deleted_at IS NULL`
		if err := r.pool.QueryRow(ctx, countQuery, fundID).Scan(&total); err != nil {
			return nil, fmt.Errorf("count cap table entries: %w", err)
		}
	}

	return &CapTableView{
		FundID:     fundID,
		Entries:    entries,
		TotalCount: total,
		Limit:      params.Limit,
		Offset:     params.Offset,
	}, nil
}

// FindByFundAndOwner retrieves a single cap table entry by fund and owner name.
func (r *PostgresRepository) FindByFundAndOwner(ctx context.Context, fundID uuid.UUID, ownerName string) (*Entry, error) {
	const query = `
		SELECT id, fund_id, owner_name, units, acquired_at, updated_at, deleted_at
		FROM cap_table_entries
		WHERE fund_id = $1 AND owner_name = $2 AND deleted_at IS NULL
	`
	var entry Entry
	err := r.pool.QueryRow(ctx, query, fundID, ownerName).Scan(
		&entry.ID,
		&entry.FundID,
		&entry.OwnerName,
		&entry.Units,
		&entry.AcquiredAt,
		&entry.UpdatedAt,
		&entry.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, OwnerNotFoundError(fundID, ownerName)
		}
		return nil, fmt.Errorf("find owner %q in fund %s: %w", ownerName, fundID, err)
	}
	return &entry, nil
}

// Upsert creates or updates a cap table entry.
func (r *PostgresRepository) Upsert(ctx context.Context, entry *Entry) error {
	return r.upsert(ctx, r.pool, entry)
}

// UpsertTx creates or updates a cap table entry within the provided transaction.
func (r *PostgresRepository) UpsertTx(ctx context.Context, tx pgx.Tx, entry *Entry) error {
	return r.upsert(ctx, tx, entry)
}

func (r *PostgresRepository) upsert(ctx context.Context, q querier, entry *Entry) error {
	if entry == nil {
		return ErrNilEntry
	}

	// Use ON CONFLICT to upsert. Preserve acquired_at on update.
	// The EXCLUDED pseudo-table refers to the values we're trying to insert.
	const query = `
		INSERT INTO cap_table_entries (id, fund_id, owner_name, units, acquired_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (fund_id, owner_name) DO UPDATE SET
			units = EXCLUDED.units,
			updated_at = EXCLUDED.updated_at,
			deleted_at = EXCLUDED.deleted_at
		RETURNING id, acquired_at
	`
	var returnedID uuid.UUID
	err := q.QueryRow(ctx, query, entry.ID, entry.FundID, entry.OwnerName, entry.Units, entry.AcquiredAt, entry.UpdatedAt, entry.DeletedAt).Scan(&returnedID, &entry.AcquiredAt)
	if err != nil {
		return fmt.Errorf("upsert cap table entry for owner %q in fund %s: %w", entry.OwnerName, entry.FundID, err)
	}

	// Update the entry ID with the actual ID from the database.
	// This handles the case where an existing row was updated.
	entry.ID = returnedID
	return nil
}
