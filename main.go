package gorofs

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"time"
)

type raw struct {
	buf []byte
}

func (r *raw) ReadAt(b []byte, off int64) (n int, err error) {
	if off > int64(len(r.buf)) {
		return 0, fmt.Errorf("Offset %v beyond EOF at %v", off, len(r.buf))
	}
	n = copy(b, r.buf[int(off):])
	if n < len(b) {
		err = io.EOF
	}
	return
}

type file struct {
	io.ReadCloser
	zf *zip.File
}

func (f *file) Name() string {
	return f.zf.Name
}

func (f *file) ModTime() time.Time {
	return f.zf.FileHeader.ModTime()
}

func (f *file) Mode() os.FileMode {
	m, _ := f.zf.FileHeader.Mode()
	return m
}

func (f *file) IsDir() bool {
	return false
}

func (f *file) Size() int64 {
	return int64(f.zf.UncompressedSize)
}

func (f *file) Stat() (fi os.FileInfo, err error) {
	return f, nil
}

var defaultFs *Rofs

func Register(buf []byte) (err error) {
	fs, err := NewROFS(buf)
	if err != nil {
		return
	}
	defaultFs = fs
	return
}

func Open(name string) (ReadStatCloser, error) {
	return defaultFs.Open(name)
}

type Rofs struct {
	reader *zip.Reader
}

func NewROFS(buf []byte) (fs *Rofs, err error) {
	data := &raw{buf: buf}
	r, err := zip.NewReader(data, int64(len(data.buf)))
	if err != nil {
		return
	}
	fs = &Rofs{reader: r}
	return
}

func (fs *Rofs) Open(name string) (f ReadStatCloser, err error) {
	f, err = os.Open(name)
	if err == nil {
		return
	}
	if fs == nil {
		return nil, fmt.Errorf("No fs registered")
	}
	for _, fz := range fs.reader.File {
		if fz.Name == name {
			fi, err := fz.Open()
			if err != nil {
				continue
			}
			f = &file{fi, fz}
			return f, nil
		}
	}
	return nil, fmt.Errorf("File not found")
}

type ReadStatCloser interface {
	io.ReadCloser
	Stat() (os.FileInfo, error)
}
