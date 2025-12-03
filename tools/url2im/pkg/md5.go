package pkg

import (
	"crypto/md5"
	"encoding/hex"
	"hash"
	"io"
)

func NewMd5Reader(r io.Reader) *Md5Reader {
	return &Md5Reader{h: md5.New(), r: r}
}

type Md5Reader struct {
	h hash.Hash
	r io.Reader
}

func (r *Md5Reader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	if err == nil && n > 0 {
		r.h.Write(p[:n])
	}
	return
}

func (r *Md5Reader) Md5() string {
	return hex.EncodeToString(r.h.Sum(nil))
}
