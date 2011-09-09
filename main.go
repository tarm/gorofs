package gorofs

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
)

type raw struct {
	buf []byte
}

func (r *raw) ReadAt(b []byte, off int64) (n int, err os.Error) {
	if off > int64(len(r.buf)) {
		return 0, fmt.Errorf("Offset %v beyond EOF at %v", off, len(r.buf))
	}
	n = copy(b, r.buf[int(off):])
	if n < len(b) {
		err = os.EOF
	}
	return
}

type file struct {
	io.ReadCloser
	zf *zip.File
}

func (f *file) Stat() (fi *os.FileInfo, err os.Error) {
	fi = new(os.FileInfo)
	fi.Mtime_ns = f.zf.Mtime_ns()
	fi.Name = f.zf.Name
	fi.Size = int64(f.zf.UncompressedSize)
	return
}

var defaultFs *Rofs

func Register(buf []byte) (err os.Error) {
	fs, err := NewROFS(buf)
	if err != nil {
		return
	}
	defaultFs = fs
	return
}

func Open(name string) (ReadStatCloser, os.Error) {
	return defaultFs.Open(name)
}

type Rofs struct {
	reader *zip.Reader
}

func NewROFS(buf []byte) (fs *Rofs, err os.Error) {
	data := &raw{buf: buf}
	r, err := zip.NewReader(data, int64(len(data.buf)))
	if err != nil {
		return
	}
	fs = &Rofs{reader: r}
	return
}

func (fs *Rofs) Open(name string) (f ReadStatCloser, err os.Error) {
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
	Stat() (*os.FileInfo, os.Error)
}
