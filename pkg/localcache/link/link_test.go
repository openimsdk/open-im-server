package link

import (
	"testing"
)

func TestName(t *testing.T) {

	v := New(1)

	//v.Link("a:1", "b:1", "c:1", "d:1")
	v.Link("a:1", "b:1", "c:1")
	v.Link("z:1", "b:1")

	//v.DelKey("a:1")
	v.Del("z:1")

	t.Log(v)

}
