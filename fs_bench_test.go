package litefs_test

import (
	"context"
	"io/fs"
	"os"
	"testing"
	"time"

	"go.husin.dev/litefs"
	"go.husin.dev/litefs/config"
)

const testDirPath = "./sql"

func BenchmarkLiteFSMerge(b *testing.B) {
	lfs, err := litefs.NewFS(&config.Config{
		BlobPath:  b.TempDir(),
		Database:  dbPath,
		DBTimeout: 5 * time.Second,
	})
	if err != nil {
		b.Fatal(err)
	}

	b.StopTimer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		purgeFS(b, lfs)()
		b.StartTimer()
		err := lfs.Merge(context.Background(), os.DirFS(testDirPath), ".")
		b.StopTimer()
		if err != nil {
			b.Error(err)
		}
	}

	purgeFS(b, lfs)()
}

func BenchmarkLiteFSWalkNoop(b *testing.B) {
	lfs, err := litefs.NewFS(&config.Config{
		BlobPath:  b.TempDir(),
		Database:  dbPath,
		DBTimeout: 5 * time.Second,
	})
	if err != nil {
		b.Fatal(err)
	}
	purgeFS(b, lfs)()
	b.Cleanup(purgeFS(b, lfs))

	if err := lfs.Merge(context.Background(), os.DirFS(testDirPath), "."); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := fs.WalkDir(lfs, ".", func(path string, d fs.DirEntry, err error) error {
			return err
		}); err != nil {
			b.Error(err)
		}
	}
	b.StopTimer()
}

func BenchmarkOSFSWalkNoop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := fs.WalkDir(os.DirFS(testDirPath), ".", func(path string, d fs.DirEntry, err error) error {
			return err
		}); err != nil {
			b.Error(err)
		}
	}
}
