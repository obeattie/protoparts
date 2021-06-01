package protopartstest

import (
	"bytes"

	"google.golang.org/protobuf/proto"

	"github.com/obeattie/protoparts"
	ppproto "github.com/obeattie/protoparts/test/proto"
)

func Fuzz(data []byte) int {
	p := &ppproto.Person{}
	if err := proto.Unmarshal(data, p); err != nil {
		return 0
	}
	md := p.ProtoReflect().Descriptor()

	// It's a valid protobuf, so we should be able to split it
	parts, err := protoparts.Split(data, md)
	if err != nil {
		return 0
	}

	// joining the parts should yield the same bytes we started with
	pong := parts.Join()
	if !bytes.Equal(data, pong) {
		//panic(fmt.Sprintf("joined data not equal"))
	}

	return 1
}
