package cont

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type multipartUploadID struct {
	Type int    `json:"a,omitempty"`
	ID   string `json:"b,omitempty"`
	Key  string `json:"c,omitempty"`
	Size int64  `json:"d,omitempty"`
	Hash string `json:"e,omitempty"`
}

func newMultipartUploadID(id multipartUploadID) string {
	data, err := json.Marshal(id)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(data)
}

func parseMultipartUploadID(id string) (*multipartUploadID, error) {
	data, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return nil, fmt.Errorf("invalid multipart upload id: %w", err)
	}
	var upload multipartUploadID
	if err := json.Unmarshal(data, &upload); err != nil {
		return nil, fmt.Errorf("invalid multipart upload id: %w", err)
	}
	return &upload, nil
}
