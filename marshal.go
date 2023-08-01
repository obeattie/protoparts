package protoparts

import (
	"fmt"

	"google.golang.org/protobuf/proto"
)

var marshalOpts = proto.MarshalOptions{
	Deterministic: true,
}

// Marshal is a convenient way to construct Parts from a Protocol Buffer message struct, without the need for the caller
// to do their own marshaling before calling Split().
func Marshal(m proto.Message) (Parts, error) {
	pb, err := marshalOpts.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("could not marshal proto message: %w", err)
	}
	return Split(pb, m.ProtoReflect().Descriptor())
}
