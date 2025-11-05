package msggateway

import (
	"crypto/rand"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"unsafe"
)

func mockRandom() []byte {
	bs := make([]byte, 50)
	rand.Read(bs)
	return bs
}

func TestCompressDecompress(t *testing.T) {

	compressor := NewGzipCompressor()

	for i := 0; i < 2000; i++ {
		src := mockRandom()

		// compress
		dest, err := compressor.CompressWithPool(src)
		if err != nil {
			t.Log(err)
		}
		assert.Equal(t, nil, err)

		// decompress
		res, err := compressor.DecompressWithPool(dest)
		if err != nil {
			t.Log(err)
		}
		assert.Equal(t, nil, err)

		// check
		assert.EqualValues(t, src, res)
	}
}

func TestCompressDecompressWithConcurrency(t *testing.T) {
	wg := sync.WaitGroup{}
	compressor := NewGzipCompressor()

	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			src := mockRandom()

			// compress
			dest, err := compressor.CompressWithPool(src)
			if err != nil {
				t.Log(err)
			}
			assert.Equal(t, nil, err)

			// decompress
			res, err := compressor.DecompressWithPool(dest)
			if err != nil {
				t.Log(err)
			}
			assert.Equal(t, nil, err)

			// check
			assert.EqualValues(t, src, res)

		}()
	}
	wg.Wait()
}

func BenchmarkCompress(b *testing.B) {
	src := mockRandom()
	compressor := NewGzipCompressor()

	for i := 0; i < b.N; i++ {
		_, err := compressor.Compress(src)
		assert.Equal(b, nil, err)
	}
}

func BenchmarkCompressWithSyncPool(b *testing.B) {
	src := mockRandom()

	compressor := NewGzipCompressor()
	for i := 0; i < b.N; i++ {
		_, err := compressor.CompressWithPool(src)
		assert.Equal(b, nil, err)
	}
}

func BenchmarkDecompress(b *testing.B) {
	src := mockRandom()

	compressor := NewGzipCompressor()
	comdata, err := compressor.Compress(src)

	assert.Equal(b, nil, err)

	for i := 0; i < b.N; i++ {
		_, err := compressor.DeCompress(comdata)
		assert.Equal(b, nil, err)
	}
}

func BenchmarkDecompressWithSyncPool(b *testing.B) {
	src := mockRandom()

	compressor := NewGzipCompressor()
	comdata, err := compressor.Compress(src)
	assert.Equal(b, nil, err)

	for i := 0; i < b.N; i++ {
		_, err := compressor.DecompressWithPool(comdata)
		assert.Equal(b, nil, err)
	}
}

func TestName(t *testing.T) {
	t.Log(unsafe.Sizeof(Client{}))

}
