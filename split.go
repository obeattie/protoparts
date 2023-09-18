package protoparts

import (
	"fmt"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func splitPb(pb []byte, md, originalMd protoreflect.MessageDescriptor, prefix Path) (Parts, error) {
	v := Parts{}
	indices := map[protowire.Number]int{} // Tracks the last seen index in repeated fields
	for len(pb) > 0 {
		num, typ, typLength := protowire.ConsumeTag(pb)
		if typLength < 0 {
			return nil, fmt.Errorf("error parsing tag: %w", protowire.ParseError(typLength))
		}
		valueLength := protowire.ConsumeFieldValue(num, typ, pb[typLength:])
		if valueLength < 0 {
			return nil, fmt.Errorf("error parsing field value: %w", protowire.ParseError(valueLength))
		}
		value := pb[:typLength+valueLength]
		pb = pb[len(value):]

		// Ignore unknown fields
		fd := md.Fields().ByNumber(num)
		if fd == nil {
			continue
		}

		// Compile the path for the value
		idx := -1
		if fd.Cardinality() == protoreflect.Repeated && !fd.IsMap() {
			idx = indices[num]
			indices[num]++
		}
		path := append(prefix, PathTerm{
			Tag:   num,
			Index: idx,
		})

		// If we have a nested message, it's encoded recursively as bytes
		if typ == protowire.BytesType && fd.Kind() == protoreflect.MessageKind {
			nested, n := protowire.ConsumeBytes(value[typLength:])
			if n < 0 {
				return nil, fmt.Errorf("error consuming map bytes: %w", protowire.ParseError(n))
			}
			if len(nested) > 0 {
				children, err := splitPb(nested, fd.Message(), originalMd, path)
				if err != nil {
					return nil, err
				}
				if len(children) > 0 && fd.IsMap() {
					// Special treatment is needed for maps. Their key/value pairs are encoded as repeated nested
					// messages, each with two fields (key (1) and value (2)). But, to allow values to be addressed by
					// their keys, we want the value to be represented as a 'top level' Part, with the key forming part
					// of the path.
					//
					// So, we need to manipulate the values' paths: at this point they look like  ….parent.2.…, but
					// they need to be changed to ….parent[key].…
					key, values := children[0], children[1:]
					keyTermIndex, keyTerm := len(prefix), PathTerm{
						Tag:   num,
						Index: idx,
						Key:   key.Bytes,
					}
					for i := range values {
						// Add the key to the "parent" term
						p := values[i].Path
						p[keyTermIndex] = keyTerm
						// Snip the superfluous field number term (https://github.com/golang/go/wiki/SliceTricks#delete)
						p = append(p[:keyTermIndex+1], p[keyTermIndex+2:]...)
						values[i].Path = p
						values[i].KeyType = key.Type
					}
					children = values // Toss the key as a dedicated element now it's part of the path of each value
				}
				v = append(v, children...)
				continue
			}
		}

		v = append(v, Part{
			Path:  path,
			Type:  typ,
			Bytes: value[typLength:],
			Md:    originalMd,
		})
	}

	return v, nil
}

// Split explodes a serialised Protocol Buffer into parts of its constituent fields and those within nested
// messages. The resulting parts are returned in the order they appear in, so recombining the parts (unless they
// are reordered either manually or by sorting them) yields a byte-wise identical message.
func Split(pb []byte, md protoreflect.MessageDescriptor) (Parts, error) {
	return splitPb(pb, md, md, nil)
}
