package pkg

import (
	"bytes"
	"io"
	"os"
)

type ReadSeekSizeCloser interface {
	io.ReadSeekCloser
	Size() int64
}

func NewReader(r io.Reader, max int64, path string) (ReadSeekSizeCloser, error) {
	buf := make([]byte, max+1)
	n, err := io.ReadFull(r, buf)
	if err == nil {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o666)
		if err != nil {
			return nil, err
		}
		var ok bool
		defer func() {
			if !ok {
				_ = f.Close()
				_ = os.Remove(path)
			}
		}()
		if _, err := f.Write(buf[:n]); err != nil {
			return nil, err
		}
		cn, err := io.Copy(f, r)
		if err != nil {
			return nil, err
		}
		if _, err := f.Seek(0, io.SeekStart); err != nil {
			return nil, err
		}
		ok = true
		return &fileBuffer{
			f: f,
			n: cn + int64(n),
		}, nil
	} else if err == io.EOF || err == io.ErrUnexpectedEOF {
		return &memoryBuffer{
			r: bytes.NewReader(buf[:n]),
		}, nil
	} else {
		return nil, err
	}
}

type fileBuffer struct {
	n int64
	f *os.File
}

func (r *fileBuffer) Read(p []byte) (n int, err error) {
	return r.f.Read(p)
}

func (r *fileBuffer) Seek(offset int64, whence int) (int64, error) {
	return r.f.Seek(offset, whence)
}

func (r *fileBuffer) Size() int64 {
	return r.n
}

func (r *fileBuffer) Close() error {
	name := r.f.Name()
	if err := r.f.Close(); err != nil {
		return err
	}
	return os.Remove(name)
}

type memoryBuffer struct {
	r *bytes.Reader
}

func (r *memoryBuffer) Read(p []byte) (n int, err error) {
	return r.r.Read(p)
}

func (r *memoryBuffer) Seek(offset int64, whence int) (int64, error) {
	return r.r.Seek(offset, whence)
}

func (r *memoryBuffer) Close() error {
	return nil
}

func (r *memoryBuffer) Size() int64 {
	return r.r.Size()
}
