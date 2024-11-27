// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package msggateway

import (
	"crypto/rand"
	"sync"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
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
