package litefs

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/xid"

	"go.husin.dev/litefs/internal/db"
	"go.husin.dev/litefs/internal/log"
)

func (f *FS) txListEntries(ctx context.Context, dbtx *db.Queries, path string) ([]db.LiteFSEntry, error) {
	if path == "." {
		return dbtx.ListRootEntries(ctx)
	}

	e, err := f.txEntryLookup(ctx, dbtx, path)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []db.LiteFSEntry{}, nil
		}
		return nil, fmt.Errorf("looking up entry '%s': %w", path, err)
	}

	return dbtx.ListEntries(ctx, sql.NullString{String: e.ID, Valid: true})
}

func (f *FS) txEntryLookup(ctx context.Context, dbtx *db.Queries, path string) (*db.LiteFSEntry, error) {
	if path == "." {
		return &db.LiteFSEntry{
			ID:      "",
			Name:    "_root",
			Modtime: "",
			Content: nil,
		}, nil
	}

	entries := strings.Split(path, string(filepath.Separator))
	parentID := sql.NullString{Valid: false}

	var entry db.LiteFSEntry
	for _, name := range entries {
		var e db.LiteFSEntry
		var err error
		if parentID.Valid {
			e, err = dbtx.LookupEntry(ctx, db.LookupEntryParams{
				ParentID: parentID,
				Name:     name,
			})
		} else {
			e, err = dbtx.LookupRootEntry(ctx, name)
		}
		if err != nil {
			log.Err(err).Str("path", path).Str("name", name).Send()
			return nil, fmt.Errorf("lookup of '%s' in '%s': %w", name, path, err)
		}

		entry = e
		parentID.Valid = true
		parentID.String = e.ID
	}

	return &entry, nil
}

func (f *FS) txMkAll(ctx context.Context, dbtx *db.Queries, entries []string) (string, error) {
	parentID := sql.NullString{Valid: false}
	now := time.Now().Format(time.RFC3339)

	for _, name := range entries {
		if name == "" {
			log.Warn().Strs("path", entries).Msg("skipping empty string")
			continue
		}
		var e db.LiteFSEntry
		var err error
		if parentID.Valid {
			e, err = dbtx.LookupEntry(ctx, db.LookupEntryParams{
				ParentID: parentID,
				Name:     name,
			})
		} else {
			e, err = dbtx.LookupRootEntry(ctx, name)
		}
		if err == nil {
			parentID.String = e.ID
			parentID.Valid = true
			continue
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("querying for existing row: %w", err)
		}

		id := xid.New().String()
		err = dbtx.CreateEntry(ctx, db.CreateEntryParams{
			ID:       id,
			ParentID: parentID,
			Name:     name,
			Modtime:  now,
			Content:  nil,
		})
		if err != nil {
			log.Debug().Err(err).Strs("path", entries).Str("name", name).Msg("creating entry")
			return "", err
		}
		log.Debug().Strs("path", entries).Str("name", name).Str("id", id).Send()

		parentID.Valid = true
		parentID.String = id
	}

	return parentID.String, nil
}

func (f *FS) txWrite(ctx context.Context, dbtx *db.Queries, path string, content []byte) (string, error) {
	entries := strings.Split(path, string(filepath.Separator))

	parentID, err := f.txMkAll(ctx, dbtx, entries[:len(entries)-1])
	if err != nil {
		return "", fmt.Errorf("creating parent directories: %w", err)
	}

	eID := xid.New().String()
	err = dbtx.CreateEntry(ctx, db.CreateEntryParams{
		ID:       eID,
		ParentID: sql.NullString{String: parentID, Valid: true},
		Name:     entries[len(entries)-1],
		Modtime:  time.Now().Format(time.RFC3339),
		Content:  content,
	})
	if err != nil {
		return "", fmt.Errorf("writing blob content: %w", err)
	}

	return eID, nil
}
