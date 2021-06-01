package protoparts

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

// fieldDescriptorInMessage returns the field descriptor for the field at the given path, or nil if the path is invalid.
func fieldDescriptorInMessage(md protoreflect.MessageDescriptor, p Path) protoreflect.FieldDescriptor {
	if len(p) == 0 {
		return nil
	}
	head, tail := p[0], p[1:]
	fd := md.Fields().ByNumber(head.Tag)
	if head.Key != nil && fd != nil {
		if !fd.IsMap() {
			return nil // not a map; path is incorrect
		}
		fd = fd.MapValue()
	}
	if len(tail) > 0 {
		if fd != nil && fd.Kind() == protoreflect.MessageKind { // Nested message: recurse
			return fieldDescriptorInMessage(fd.Message(), tail)
		}
		return nil // We have a tail but no nested message; path is incorrect
	}
	return fd
}

// valueFromMessage returns the value at the given path within the message. By recursively following enclosed messages
// it is useful to extract deeply-nested attributes.
//
// The bool returned indicates whether the value was found, and if it was, whether it was populated.
func valueFromMessage(m protoreflect.Message, p Path) (protoreflect.Value, bool) {
	if p == nil {
		return protoreflect.Value{}, false
	}
	v, has := protoreflect.ValueOfMessage(m), true
	for _, term := range p {
		m := v.Message()
		if m == nil { // We have a tail but no nested message; path is incorrect
			return protoreflect.Value{}, false
		}
		fd := m.Descriptor().Fields().ByNumber(term.Tag)
		if fd == nil {
			return protoreflect.Value{}, false
		}
		v, has = m.Get(fd), m.Has(fd)
		if term.Index != -1 && fd.Cardinality() == protoreflect.Repeated {
			l := v.List()
			if term.Index >= l.Len() {
				return protoreflect.Value{}, false // out of range
			}
			v = l.Get(term.Index)
		}
		if term.Key != nil {
			if !fd.IsMap() {
				return protoreflect.Value{}, false // not a map; path is incorrect
			}
			entry := dynamicpb.NewMessage(fd.Message())
			if err := proto.Unmarshal(term.Key, entry); err != nil {
				return protoreflect.Value{}, false // invalid key
			}
			v = v.Map().Get(entry.Get(fd.MapKey()).MapKey())
			has = v.IsValid()
		}
	}
	return v, has
}
