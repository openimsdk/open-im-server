package msggateway

import (
	"testing"
)

func TestGobEncoder_Encode(t *testing.T) {
	encoder := NewGobEncoder()

	// 测试用例1: 编码 []byte 数据
	inputData := []byte("test data")
	encodedData, err := encoder.Encode(inputData)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if string(encodedData) != string(inputData) {
		t.Fatalf("expected encoded data to be '%s', got '%s'", inputData, encodedData)
	}

	// 测试用例2: 编码非 []byte 数据
	nonByteData := "string data"
	_, err = encoder.Encode(nonByteData)
	if err == nil {
		t.Fatalf("expected an error when encoding non-byte data, got none")
	}
}

func TestGobEncoder_Decode(t *testing.T) {
	encoder := NewGobEncoder()

	// 测试用例1: 解码到 []byte 数据
	encodedData := []byte("test data")
	var decodedData []byte
	err := encoder.Decode(encodedData, &decodedData)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if string(decodedData) != string(encodedData) {
		t.Fatalf("expected decoded data to be '%s', got '%s'", encodedData, decodedData)
	}

	// 测试用例2: 解码到非 []byte 数据
	var nonByteData string
	err = encoder.Decode(encodedData, &nonByteData)
	if err == nil {
		t.Fatalf("expected an error when decoding to non-byte data, got none")
	}
}
