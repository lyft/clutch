//go:build !withAssets

package assets

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	pathpkg "path"
	"time"
)

const warning = `
Server was not generated with embedded assets. Only API endpoints are available.
For development, you probably want to use the development server.
For production, either embed assets using the generator and building with the withAssets tag or serve your frontend separately from a CDN or a different server.
`

// Below is a simple asset server based on vfsgen for purpose of serving the warning. Do not modify.

// assets statically implements the virtual filesystem provided to vfsgen.
var VirtualFS = func() http.FileSystem {
	fs := vfsgen۰FS{
		"/": &vfsgen۰DirInfo{
			name:    "/",
			modTime: time.Time{},
		},
		"/index.html": &vfsgen۰FileInfo{
			name:    "index.html",
			modTime: time.Time{},
			content: []byte(warning),
		},
	}
	fs["/"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/index.html"].(os.FileInfo),
	}

	return fs
}()

type vfsgen۰FS map[string]interface{}

func (fs vfsgen۰FS) Open(path string) (http.File, error) {
	path = pathpkg.Clean("/" + path)
	f, ok := fs[path]
	if !ok {
		return nil, &os.PathError{Op: "open", Path: path, Err: os.ErrNotExist}
	}

	switch f := f.(type) {
	case *vfsgen۰FileInfo:
		return &vfsgen۰File{
			vfsgen۰FileInfo: f,
			Reader:          bytes.NewReader(f.content),
		}, nil
	case *vfsgen۰DirInfo:
		return &vfsgen۰Dir{
			vfsgen۰DirInfo: f,
		}, nil
	default:
		// This should never happen because we generate only the above types.
		panic(fmt.Sprintf("unexpected type %T", f))
	}
}

// vfsgen۰FileInfo is a static definition of an uncompressed file (because it's not worth gzip compressing).
type vfsgen۰FileInfo struct {
	name    string
	modTime time.Time
	content []byte
}

func (f *vfsgen۰FileInfo) Readdir(count int) ([]os.FileInfo, error) {
	return nil, fmt.Errorf("cannot Readdir from file %s", f.name)
}
func (f *vfsgen۰FileInfo) Stat() (os.FileInfo, error) { return f, nil }

func (f *vfsgen۰FileInfo) NotWorthGzipCompressing() {}

func (f *vfsgen۰FileInfo) Name() string       { return f.name }
func (f *vfsgen۰FileInfo) Size() int64        { return int64(len(f.content)) }
func (f *vfsgen۰FileInfo) Mode() os.FileMode  { return 0444 }
func (f *vfsgen۰FileInfo) ModTime() time.Time { return f.modTime }
func (f *vfsgen۰FileInfo) IsDir() bool        { return false }
func (f *vfsgen۰FileInfo) Sys() interface{}   { return nil }

// vfsgen۰File is an opened file instance.
type vfsgen۰File struct {
	*vfsgen۰FileInfo
	*bytes.Reader
}

func (f *vfsgen۰File) Close() error {
	return nil
}

// vfsgen۰DirInfo is a static definition of a directory.
type vfsgen۰DirInfo struct {
	name    string
	modTime time.Time
	entries []os.FileInfo
}

func (d *vfsgen۰DirInfo) Read([]byte) (int, error) {
	return 0, fmt.Errorf("cannot Read from directory %s", d.name)
}
func (d *vfsgen۰DirInfo) Close() error               { return nil }
func (d *vfsgen۰DirInfo) Stat() (os.FileInfo, error) { return d, nil }

func (d *vfsgen۰DirInfo) Name() string       { return d.name }
func (d *vfsgen۰DirInfo) Size() int64        { return 0 }
func (d *vfsgen۰DirInfo) Mode() os.FileMode  { return 0755 | os.ModeDir }
func (d *vfsgen۰DirInfo) ModTime() time.Time { return d.modTime }
func (d *vfsgen۰DirInfo) IsDir() bool        { return true }
func (d *vfsgen۰DirInfo) Sys() interface{}   { return nil }

// vfsgen۰Dir is an opened dir instance.
type vfsgen۰Dir struct {
	*vfsgen۰DirInfo
	pos int // Position within entries for Seek and Readdir.
}

func (d *vfsgen۰Dir) Seek(offset int64, whence int) (int64, error) {
	if offset == 0 && whence == io.SeekStart {
		d.pos = 0
		return 0, nil
	}
	return 0, fmt.Errorf("unsupported Seek in directory %s", d.name)
}

func (d *vfsgen۰Dir) Readdir(count int) ([]os.FileInfo, error) {
	if d.pos >= len(d.entries) && count > 0 {
		return nil, io.EOF
	}
	if count <= 0 || count > len(d.entries)-d.pos {
		count = len(d.entries) - d.pos
	}
	e := d.entries[d.pos : d.pos+count]
	d.pos += count
	return e, nil
}
