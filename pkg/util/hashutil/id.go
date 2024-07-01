package hashutil

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
)

func IdHash(ids []string) uint64 {
	if len(ids) == 0 {
		return 0
	}
	data, _ := json.Marshal(ids)
	sum := md5.Sum(data)
	return binary.BigEndian.Uint64(sum[:])
}
