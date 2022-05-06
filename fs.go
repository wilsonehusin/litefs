package litefs

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path/filepath"
	"time"

	//_ "github.com/mattn/go-sqlite3"
	//_ "modernc.org/sqlite"
	_ "github.com/tailscale/sqlite"

	"go.husin.dev/litefs/config"
	"go.husin.dev/litefs/internal/db"
	"go.husin.dev/litefs/internal/log"
)

//go:embed sql/schema/schema.sql
var schema string

const sqliteDriver = "sqlite3"

// FS implements fs.FS, which provides structure in accessing Blob items.
type FS struct {
	root string
	db   *sql.DB

	timeout time.Duration
}

func NewFS(cfg *config.Config) (*FS, error) {
	db, err := sql.Open(sqliteDriver, cfg.Database+"?_fk=true")
	if err != nil {
		return nil, fmt.Errorf("opening '%s' database on '%s': %w", sqliteDriver, cfg.Database, err)
	}
	root, err := filepath.Abs(cfg.BlobPath)
	if err != nil {
		return nil, fmt.Errorf("calculating absolute path of '%s': %w", cfg.BlobPath, err)
	}

	return &FS{
		root:    root,
		db:      db,
		timeout: cfg.DBTimeout,
	}, nil
}

func NewMemFS(cfg *config.Config) (*FS, error) {
	db, err := sql.Open(sqliteDriver, cfg.Database+"?mode=memory&_fk=true")
	if err != nil {
		return nil, fmt.Errorf("opening '%s' database on '%s': %w", sqliteDriver, cfg.Database, err)
	}
	rows, err := db.Query(schema)
	if err != nil {
		return nil, fmt.Errorf("initializing database schema: %w", err)
	}
	rows.Close()

	root, err := filepath.Abs(cfg.BlobPath)
	if err != nil {
		return nil, fmt.Errorf("calculating absolute path of '%s': %w", cfg.BlobPath, err)
	}

	return &FS{
		root:    root,
		db:      db,
		timeout: cfg.DBTimeout,
	}, nil
}

func (f *FS) Merge(ctx context.Context, src fs.FS, root string) error {
	tx, err := f.db.Begin()
	if err != nil {
		return fmt.Errorf("begin database transaction: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck

	dbtx := db.New(tx)

	err = fs.WalkDir(src, root, func(path string, d fs.DirEntry, err error) error {
		if err := ctx.Err(); err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		file, err := src.Open(path)
		if err != nil {
			return &fs.PathError{Op: "open", Path: path, Err: err}
		}

		content, err := io.ReadAll(file)
		file.Close()
		if err != nil {
			return &fs.PathError{Op: "read", Path: path, Err: err}
		}

		_, err = f.txWrite(ctx, dbtx, path, content)
		return err
	})
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (f *FS) ReadDir(path string) ([]fs.DirEntry, error) {
	if !fs.ValidPath(path) {
		return nil, &fs.PathError{Op: "readdir", Path: path, Err: fs.ErrInvalid}
	}

	path = filepath.Clean(path)
	log.Debug().Str("path", path).Str("op", "readdir").Send()

	dbtx := db.New(f.db)
	ctx, cancel := context.WithTimeout(context.Background(), f.timeout)
	defer cancel()

	dbEntries, err := f.txListEntries(ctx, dbtx, path)
	if err != nil {
		log.Debug().Err(err).Str("path", path).Send()
		return nil, &fs.PathError{Op: "list", Path: path, Err: fs.ErrNotExist}
	}

	result := []fs.DirEntry{}
	for _, e := range dbEntries {
		log.Debug().Str("path", path).Str("name", e.Name).Send()
		result = append(result, entryFromDB(e))
	}
	return result, nil
}

func (f *FS) WriteBlob(ctx context.Context, path string, content []byte) error {
	if !fs.ValidPath(path) {
		return &fs.PathError{Op: "write", Path: path, Err: fs.ErrInvalid}
	}

	path = filepath.Clean(path)

	tx, err := f.db.Begin()
	if err != nil {
		return fmt.Errorf("begin database transaction: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck

	dbtx := db.New(tx)
	eID, err := f.txWrite(ctx, dbtx, path, content)
	if err != nil {
		return fmt.Errorf("writing file: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing database transaction: %w", err)
	}

	log.Debug().Str("path", path).Str("id", eID).Send()
	return nil
}

func (f *FS) GetEntry(path string) (*Entry, error) {
	if !fs.ValidPath(path) {
		log.Debug().Str("path", path).Msg("invalid path to open")
		return nil, &fs.PathError{Op: "open", Path: path, Err: fs.ErrInvalid}
	}

	path = filepath.Clean(path)
	log.Debug().Str("path", path).Str("op", "open").Send()

	dbtx := db.New(f.db)
	ctx, cancel := context.WithTimeout(context.Background(), f.timeout)
	defer cancel()

	e, err := f.txEntryLookup(ctx, dbtx, path)
	if err != nil {
		return nil, err
	}

	return entryFromDB(*e), nil
}

func (f *FS) Open(path string) (fs.File, error) {
	e, err := f.GetEntry(path)
	if err != nil {
		log.Debug().Err(err).Str("path", path).Send()
		if errors.Is(err, &fs.PathError{}) {
			return nil, err
		}

		return nil, &fs.PathError{Op: "open", Path: path, Err: fs.ErrNotExist}
	}
	return e, nil
}

func (f *FS) Purge(ctx context.Context) error {
	tx, err := f.db.Begin()
	if err != nil {
		return fmt.Errorf("begin database transaction: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck

	dbtx := db.New(tx)
	dbEntries, err := f.txListEntries(ctx, dbtx, ".")
	if err != nil {
		return err
	}

	for _, e := range dbEntries {
		if err := dbtx.DeleteEntry(ctx, e.ID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (f *FS) Remove(path string) error {
	if !fs.ValidPath(path) {
		log.Debug().Str("path", path).Msg("invalid path to open")
		return &fs.PathError{Op: "open", Path: path, Err: fs.ErrInvalid}
	}

	path = filepath.Clean(path)
	log.Debug().Str("path", path).Str("op", "delete").Send()

	dbtx := db.New(f.db)
	ctx, cancel := context.WithTimeout(context.Background(), f.timeout)
	defer cancel()

	e, err := f.txEntryLookup(ctx, dbtx, path)
	if err != nil {
		return &fs.PathError{Op: "open", Path: path, Err: fs.ErrNotExist}
	}

	return dbtx.DeleteEntry(ctx, e.ID)
}

func (f *FS) ServeHTTP(w http.ResponseWriter, req *http.Request) {

}
