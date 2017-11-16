package diskcache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pkg/errors"
)

// Store is an on disk cache, with items cached via calls to Open.
type Store struct {
	// Dir is the directory to cache items
	Dir string

	// Component when set is reported to OpenTracing as the component.
	Component string
}

// File is an os.File, but includes the Path
type File struct {
	*os.File

	// The Path on disk for File
	Path string
}

// Fetcher returns a ReadCloser. It is used by Open if the key is not in the
// cache.
type Fetcher func(context.Context) (io.ReadCloser, error)

// Open will open a file from the local cache with key. If missing, fetcher
// will fill the cache first. Open also performs single-flighting for fetcher.
func (s *Store) Open(ctx context.Context, key string, fetcher Fetcher) (file *File, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Cached Fetch")
	if s.Component != "" {
		ext.Component.Set(span, s.Component)
	}
	defer func() {
		if err != nil {
			ext.Error.Set(span, true)
			span.SetTag("err", err.Error())
		}
		if file != nil {
			// Update modified time. Modified time is used to decide which
			// files to evict from the cache.
			touch(file.Path)
		}
		span.Finish()
	}()

	if s.Dir == "" {
		return nil, errors.New("diskcache.Store.Dir must be set")
	}

	// path uses a sha256 hash of the key since we want to use it for the
	// disk name.
	h := sha256.Sum256([]byte(key))
	path := filepath.Join(s.Dir, hex.EncodeToString(h[:])) + ".zip"
	span.LogKV("key", key, "path", path)

	// First do a fast-path, assume already on disk
	f, err := os.Open(path)
	if err == nil {
		span.SetTag("source", "fast")
		return &File{File: f, Path: path}, nil
	}

	// We have to grab the lock for this key, so we can fetch or wait for
	// someone else to finish fetching.
	urlMu := urlMu(path)
	urlMu.Lock()
	defer urlMu.Unlock()
	span.LogEvent("urlMu acquired")

	// Since we acquired urlMu, someone else may have put the archive onto
	// the disk.
	f, err = os.Open(path)
	if err == nil {
		span.SetTag("source", "other")
		return &File{File: f, Path: path}, nil
	}
	// Just in case we failed due to something bad on the FS, remove
	_ = os.Remove(path)

	// Fetch since we still can't open up the file
	span.SetTag("source", "fetch")
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, errors.Wrap(err, "could not create archive cache dir")
	}

	// We write to a temporary path to prevent another Open finding a
	// partialy written file.
	tmpPath := path + ".part"
	f, err = os.OpenFile(tmpPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create temporary archive cache item")
	}
	defer os.Remove(tmpPath)

	// We are now ready to actually fetch the file. Write it to the
	// partial file and cleanup.
	r, err := fetcher(ctx)
	if err != nil {
		f.Close()
		return nil, errors.Wrap(err, "failed to fetch missing archive cache item")
	}
	err = copyAndClose(f, r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch missing archive cache item")
	}

	// Put the partially written file in the correct place and open
	err = os.Rename(tmpPath, path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to put cache item in place")
	}
	f, err = os.Open(path)
	if err != nil {
		return nil, err
	}
	return &File{File: f, Path: path}, nil
}

// EvictStats is information gathered during Evict.
type EvictStats struct {
	// CacheSize is the size of the cache before evicting.
	CacheSize int64

	// Evicted is the number of items evicted.
	Evicted int
}

// Evict will remove files from Store.Dir until it is smaller than
// maxCacheSizeBytes. It evicts files with the oldest modification time first.
func (s *Store) Evict(maxCacheSizeBytes int64) (stats EvictStats, err error) {
	isZip := func(fi os.FileInfo) bool {
		return strings.HasSuffix(fi.Name(), ".zip")
	}

	list, err := ioutil.ReadDir(s.Dir)
	if err != nil {
		return stats, errors.Wrapf(err, "failed to ReadDir %s", s.Dir)
	}

	// Sum up the total size of all zips
	var size int64
	for _, fi := range list {
		if isZip(fi) {
			size += fi.Size()
		}
	}
	stats.CacheSize = size

	// Nothing to evict
	if size <= maxCacheSizeBytes {
		return stats, nil
	}

	// Keep removing files until we are under the cache size. Remove the
	// oldest first.
	sort.Slice(list, func(i, j int) bool {
		return list[i].ModTime().Before(list[j].ModTime())
	})
	for _, fi := range list {
		if size <= maxCacheSizeBytes {
			break
		}
		if !isZip(fi) {
			continue
		}
		path := filepath.Join(s.Dir, fi.Name())
		err = os.Remove(path)
		if err != nil {
			log.Printf("failed to remove %s: %s", path, err)
			continue
		}
		stats.Evicted++
		size -= fi.Size()
	}

	return stats, nil
}

func copyAndClose(dst io.WriteCloser, src io.ReadCloser) error {
	_, err := io.Copy(dst, src)
	if err1 := src.Close(); err == nil {
		err = err1
	}
	if err1 := dst.Close(); err == nil {
		err = err1
	}
	return err
}

// touch updates the modified time to time.Now(). It is best-effort, and will
// log if it fails.
func touch(path string) {
	t := time.Now()
	if err := os.Chtimes(path, t, t); err != nil {
		log.Printf("failed to touch %s: %s", path, err)
	}
}
