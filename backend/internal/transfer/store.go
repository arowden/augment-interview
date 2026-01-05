package transfer

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DB interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

type Store struct {
	db DB
}

func NewStore(db DB) *Store {
	if db == nil {
		return nil
	}
	return &Store{db: db}
}

func (s *Store) Create(ctx context.Context, transfer *Transfer) error {
	return s.create(ctx, s.db, transfer)
}

func (s *Store) CreateTx(ctx context.Context, tx pgx.Tx, transfer *Transfer) error {
	return s.create(ctx, tx, transfer)
}

func (s *Store) create(ctx context.Context, db DB, transfer *Transfer) error {
	if transfer == nil {
		return ErrNilTransfer
	}

	const query = `
		INSERT INTO transfers (id, fund_id, from_owner, to_owner, units, idempotency_key, transferred_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		RETURNING transferred_at
	`
	err := db.QueryRow(ctx, query,
		transfer.ID,
		transfer.FundID,
		transfer.FromOwner,
		transfer.ToOwner,
		transfer.Units,
		transfer.IdempotencyKey,
	).Scan(&transfer.TransferredAt)
	if err != nil {
		return fmt.Errorf("create transfer %s: %w", transfer.ID, err)
	}
	return nil
}

func (s *Store) FindByFundID(ctx context.Context, fundID uuid.UUID, params ListParams) (*TransferList, error) {
	params = params.Normalize()

	const query = `
		SELECT id, fund_id, from_owner, to_owner, units, idempotency_key, transferred_at, COUNT(*) OVER() AS total
		FROM transfers
		WHERE fund_id = $1
		ORDER BY transferred_at ASC, id ASC
		LIMIT $2 OFFSET $3
	`
	rows, err := s.db.Query(ctx, query, fundID, params.Limit, params.Offset)
	if err != nil {
		return nil, fmt.Errorf("find transfers for fund %s: %w", fundID, err)
	}
	defer rows.Close()

	transfers := make([]*Transfer, 0, params.Limit)
	var total int
	for rows.Next() {
		var t Transfer
		if err := rows.Scan(
			&t.ID,
			&t.FundID,
			&t.FromOwner,
			&t.ToOwner,
			&t.Units,
			&t.IdempotencyKey,
			&t.TransferredAt,
			&total,
		); err != nil {
			return nil, fmt.Errorf("scan transfer row: %w", err)
		}
		transfers = append(transfers, &t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate transfer rows: %w", err)
	}

	if len(transfers) == 0 && params.Offset > 0 {
		const countQuery = `SELECT COUNT(*) FROM transfers WHERE fund_id = $1`
		if err := s.db.QueryRow(ctx, countQuery, fundID).Scan(&total); err != nil {
			return nil, fmt.Errorf("count transfers: %w", err)
		}
	}

	return &TransferList{
		Transfers:  transfers,
		TotalCount: total,
		Limit:      params.Limit,
		Offset:     params.Offset,
	}, nil
}

func (s *Store) FindByIdempotencyKey(ctx context.Context, tx pgx.Tx, key uuid.UUID) (*Transfer, error) {
	const query = `
		SELECT id, fund_id, from_owner, to_owner, units, idempotency_key, transferred_at
		FROM transfers
		WHERE idempotency_key = $1
	`
	var t Transfer
	err := tx.QueryRow(ctx, query, key).Scan(
		&t.ID,
		&t.FundID,
		&t.FromOwner,
		&t.ToOwner,
		&t.Units,
		&t.IdempotencyKey,
		&t.TransferredAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find transfer by idempotency key %s: %w", key, err)
	}
	return &t, nil
}
