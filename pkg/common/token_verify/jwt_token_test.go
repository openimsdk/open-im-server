package token_verify

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ParseToken(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVSUQiOiJvcGVuSU1BZG1pbiIsIlBsYXRmb3JtIjoiQVBhZCIsImV4cCI6MTY3NDYxNTA2MSwibmJmIjoxNjY2ODM4NzYxLCJpYXQiOjE2NjY4MzkwNjF9.l8RiIu6pR4ItwDOpNIDYA9LBzIcpk8r8n6NRtXjqOp8"
	_, err := GetClaimFromToken(token)
	assert.Nil(t, err)
}
