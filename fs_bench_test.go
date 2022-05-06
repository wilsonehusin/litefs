package litefs_test

import (
	"context"
	"io/fs"
	"os"
	"testing"
)

const testDirPath = "./sql"

func BenchmarkLiteFSMerge(b *testing.B) {
	lfs := initLiteFS(b)

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
}

func BenchmarkProfLiteFSWalkNoop(b *testing.B) {
	lfs := initLiteFS(b)
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

func BenchmarkProfOSFSWalkNoop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := fs.WalkDir(os.DirFS(testDirPath), ".", func(path string, d fs.DirEntry, err error) error {
			return err
		}); err != nil {
			b.Error(err)
		}
	}
}
