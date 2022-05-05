package litefs

import (
	"bytes"
	"io/fs"
	"os"
	"sync"
	"time"

	"go.husin.dev/litefs/internal/db"
)

const defaultDirectorySize = 4096

// Entry implements fs.File, fs.FileInfo, and fs.DirEntry.
// Represents the content of litefs.FS.
type Entry struct {
	name string
	mode fs.FileMode

	modtime    time.Time
	modtimeStr string

	size        int64
	modtimeOnce sync.Once

	content *bytes.Buffer
}

func entryFromDB(dbE db.LiteFSEntry) *Entry {
	size := len(dbE.Content)
	mode := fs.FileMode(os.O_RDONLY)
	var content *bytes.Buffer

	if dbE.Content == nil {
		size = defaultDirectorySize
		mode |= fs.ModeDir
	} else {
		content = bytes.NewBuffer(dbE.Content)
	}

	return &Entry{
		name:       dbE.Name,
		mode:       mode,
		modtimeStr: dbE.Modtime,
		size:       int64(size),
		content:    content,
	}
}

func (e *Entry) Name() string {
	return e.name
}

func (e *Entry) Size() int64 {
	return e.size
}

func (e *Entry) Mode() fs.FileMode {
	return e.mode
}

func (e *Entry) Type() fs.FileMode {
	return e.mode
}

func (e *Entry) ModTime() time.Time {
	e.modtimeOnce.Do(func() {
		t, err := time.Parse(time.RFC3339, e.modtimeStr)
		if err != nil {
			t = time.Unix(0, 0)
		}
		e.modtime = t
	})
	return e.modtime
}

func (e *Entry) IsDir() bool {
	return e.mode.IsDir()
}

func (e *Entry) Sys() any {
	return nil
}

func (e *Entry) Info() (fs.FileInfo, error) {
	return e, nil
}

func (e *Entry) Stat() (fs.FileInfo, error) {
	return e, nil
}

func (e *Entry) Read(dst []byte) (int, error) {
	return e.content.Read(dst)
}

func (e *Entry) Close() error {
	return nil
}
