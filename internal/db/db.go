// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0

package db

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.createEntryStmt, err = db.PrepareContext(ctx, createEntry); err != nil {
		return nil, fmt.Errorf("error preparing query CreateEntry: %w", err)
	}
	if q.deleteEntryStmt, err = db.PrepareContext(ctx, deleteEntry); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteEntry: %w", err)
	}
	if q.getEntryStmt, err = db.PrepareContext(ctx, getEntry); err != nil {
		return nil, fmt.Errorf("error preparing query GetEntry: %w", err)
	}
	if q.listEntriesStmt, err = db.PrepareContext(ctx, listEntries); err != nil {
		return nil, fmt.Errorf("error preparing query ListEntries: %w", err)
	}
	if q.listRootEntriesStmt, err = db.PrepareContext(ctx, listRootEntries); err != nil {
		return nil, fmt.Errorf("error preparing query ListRootEntries: %w", err)
	}
	if q.lookupEntryStmt, err = db.PrepareContext(ctx, lookupEntry); err != nil {
		return nil, fmt.Errorf("error preparing query LookupEntry: %w", err)
	}
	if q.lookupRootEntryStmt, err = db.PrepareContext(ctx, lookupRootEntry); err != nil {
		return nil, fmt.Errorf("error preparing query LookupRootEntry: %w", err)
	}
	if q.renameEntryStmt, err = db.PrepareContext(ctx, renameEntry); err != nil {
		return nil, fmt.Errorf("error preparing query RenameEntry: %w", err)
	}
	if q.updateEntryBlobStmt, err = db.PrepareContext(ctx, updateEntryBlob); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateEntryBlob: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.createEntryStmt != nil {
		if cerr := q.createEntryStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createEntryStmt: %w", cerr)
		}
	}
	if q.deleteEntryStmt != nil {
		if cerr := q.deleteEntryStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteEntryStmt: %w", cerr)
		}
	}
	if q.getEntryStmt != nil {
		if cerr := q.getEntryStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getEntryStmt: %w", cerr)
		}
	}
	if q.listEntriesStmt != nil {
		if cerr := q.listEntriesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listEntriesStmt: %w", cerr)
		}
	}
	if q.listRootEntriesStmt != nil {
		if cerr := q.listRootEntriesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listRootEntriesStmt: %w", cerr)
		}
	}
	if q.lookupEntryStmt != nil {
		if cerr := q.lookupEntryStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing lookupEntryStmt: %w", cerr)
		}
	}
	if q.lookupRootEntryStmt != nil {
		if cerr := q.lookupRootEntryStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing lookupRootEntryStmt: %w", cerr)
		}
	}
	if q.renameEntryStmt != nil {
		if cerr := q.renameEntryStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing renameEntryStmt: %w", cerr)
		}
	}
	if q.updateEntryBlobStmt != nil {
		if cerr := q.updateEntryBlobStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateEntryBlobStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                  DBTX
	tx                  *sql.Tx
	createEntryStmt     *sql.Stmt
	deleteEntryStmt     *sql.Stmt
	getEntryStmt        *sql.Stmt
	listEntriesStmt     *sql.Stmt
	listRootEntriesStmt *sql.Stmt
	lookupEntryStmt     *sql.Stmt
	lookupRootEntryStmt *sql.Stmt
	renameEntryStmt     *sql.Stmt
	updateEntryBlobStmt *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                  tx,
		tx:                  tx,
		createEntryStmt:     q.createEntryStmt,
		deleteEntryStmt:     q.deleteEntryStmt,
		getEntryStmt:        q.getEntryStmt,
		listEntriesStmt:     q.listEntriesStmt,
		listRootEntriesStmt: q.listRootEntriesStmt,
		lookupEntryStmt:     q.lookupEntryStmt,
		lookupRootEntryStmt: q.lookupRootEntryStmt,
		renameEntryStmt:     q.renameEntryStmt,
		updateEntryBlobStmt: q.updateEntryBlobStmt,
	}
}
