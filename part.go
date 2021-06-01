package protoparts

import (
	"bytes"
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"
)

// A Part represents a field as serialised within a Protocol Buffer message.
type Part struct {
	Path  Path
	Bytes []byte
	// Md contains the message descriptor that the Path can be found in
	Md protoreflect.MessageDescriptor
}

func (p Part) String() string {
	return fmt.Sprintf("%v=%x", p.Path, p.Bytes)
}

func (p Part) fd() protoreflect.FieldDescriptor {
	return fieldDescriptorInMessage(p.Md, p.Path)
}

// Equal returns whether this and the passed Part are equivalent.
func (p Part) Equal(p2 Part) bool {
	return p.Path.Equal(p2.Path) &&
		bytes.Equal(p.Bytes, p2.Bytes) &&
		p.Md.FullName() == p2.Md.FullName()
}
