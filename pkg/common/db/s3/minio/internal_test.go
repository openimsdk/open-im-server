package minio

import (
	"testing"
)

func TestName(t *testing.T) {
	//u, err := makeTargetURL(&minio.Client{}, "openim", "test.png", "", false, nil)
	//if err != nil {
	//	panic(err)
	//}
	//u.String()
	//t.Log(percentEncodeSlash("1234"))
	//
	//t.Log(FastRand())
	t.Log(makeTargetURL(nil, "", "", "", false, nil))
	//t.Log(privateNew("", nil))

}
