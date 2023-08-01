package protoparts

import (
	"testing"

	ppproto "github.com/obeattie/protoparts/test/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestUnmarshal(t *testing.T) {
	pb := &ppproto.Person{}
	assert.NoError(t, proto.Unmarshal([]byte("2\x80\x00"), pb))
}
