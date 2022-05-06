package litefs_test

import (
	"bytes"
	"context"
	"io"
	"io/fs"
	"os"
	"testing"
	"time"

	"go.husin.dev/litefs"
	"go.husin.dev/litefs/config"
)

const dbPath = "file:tmp/tests/litefs.db"

func initLiteFS(tb testing.TB) *litefs.FS {
	lfs, err := litefs.NewFS(&config.Config{
		Database:  dbPath,
		DBTimeout: 5 * time.Second,
	})
	if err != nil {
		tb.Fatal(err)
	}
	tb.Cleanup(purgeFS(tb, lfs))

	return lfs
}

func purgeFS(tb testing.TB, lfs *litefs.FS) func() {
	return func() {
		if err := lfs.Purge(context.Background()); err != nil {
			tb.Fatalf("error: purge litefs: %s", err.Error())
		}
	}
}

func TestFSWalkOpen(t *testing.T) {
	testPath := "my/path/to/file"
	testContent := []byte("foo bar baz")

	lfs := initLiteFS(t)
	if err := lfs.WriteBlob(context.Background(), testPath, testContent); err != nil {
		t.Fatal(err)
	}

	foundPaths := map[string]bool{
		"my":              false,
		"my/path":         false,
		"my/path/to":      false,
		"my/path/to/file": false,
	}
	if err := fs.WalkDir(lfs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			t.Logf("warn: non-nil error for directory '%s': %s", path, err.Error())
		}
		t.Logf("found: %s", path)
		foundPaths[path] = true
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	notFound := 0
	for dir, found := range foundPaths {
		if !found {
			t.Errorf("error: '%s' not found in WalkDir", dir)
			notFound++
		}
	}
	if notFound > 0 {
		t.Fatal("at least one expected path was not walked, check logs above")
	}

	f, err := lfs.Open(testPath)
	if err != nil {
		f.Close()
		t.Fatal(err)
	}

	stat, err := f.Stat()
	if err != nil {
		f.Close()
		t.Fatal(err)
	}

	if stat.IsDir() {
		f.Close()
		t.Fatal("expected file, got directory")
	}

	b, err := io.ReadAll(f)
	f.Close()
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(testContent, b) {
		t.Fatalf("content mismatch!\n\nexpected: %s\nreceived: %s", testContent, b)
	}
}

func TestFSMerge(t *testing.T) {
	sqlDir := os.DirFS(testDirPath)

	lfs := initLiteFS(t)
	if err := lfs.Merge(context.Background(), sqlDir, "."); err != nil {
		t.Fatal(err)
	}

	if err := fs.WalkDir(sqlDir, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			t.Logf("warn: non-nil error for directory '%s': %s", path, err.Error())
			return nil
		}
		t.Logf("walking: %s", path)

		e, err := lfs.GetEntry(path)
		if err != nil {
			t.Errorf("retrieving entry '%s': %s", path, err.Error())
			return nil
		}

		if d.IsDir() {
			if e.IsDir() {
				return nil
			}
			t.Errorf("expected '%s' to be directory, but is not", path)
			return nil
		}

		expectedF, err := sqlDir.Open(path)
		if err != nil {
			t.Errorf("opening file '%s': %s", path, err.Error())
			return nil
		}
		expected, err := io.ReadAll(expectedF)
		expectedF.Close()
		if err != nil {
			t.Errorf("reading file '%s': %s", path, err.Error())
			return nil
		}

		receivedF, err := lfs.Open(path)
		if err != nil {
			t.Errorf("opening entry '%s': %s", path, err.Error())
			return nil
		}
		received, err := io.ReadAll(receivedF)
		receivedF.Close()
		if err != nil {
			t.Errorf("reading entry '%s': %s", path, err.Error())
			return nil
		}

		if !bytes.Equal(expected, received) {
			t.Errorf("content mismatch!\nexpected:\n\t%s\nreceived:\n\t%s", expected, received)
		}

		return nil
	}); err != nil {
		t.Fatal(err)
	}
}
