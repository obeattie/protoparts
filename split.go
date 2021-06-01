package protoparts

import (
	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func splitPb(pb []byte, md, originalMd protoreflect.MessageDescriptor, prefix Path) (Parts, error) {
	v := Parts{}
	repeatedIndices := map[protowire.Number]int{} // incremented each time a repeated field is seen
	for len(pb) > 0 {
		num, typ, typLength := protowire.ConsumeTag(pb)
		if typLength < 0 {
			return nil, errors.Wrap(protowire.ParseError(typLength), "error parsing tag")
		}
		valueLength := protowire.ConsumeFieldValue(num, typ, pb[typLength:])
		if valueLength < 0 {
			return nil, errors.Wrap(protowire.ParseError(typLength), "error parsing field value")
		}
		value := pb[:typLength+valueLength]
		pb = pb[len(value):]

		// Compile the path for the value
		idx, fd := -1, md.Fields().ByNumber(num)
		if fd.Cardinality() == protoreflect.Repeated && !fd.IsMap() {
			idx = repeatedIndices[num]
			repeatedIndices[num]++
		}
		path := append(prefix, PathTerm{
			Tag:   num,
			Index: idx,
		})

		// If we have a nested message on our hands, it's encoded recursively as bytes.
		if typ == protowire.BytesType && fd.Kind() == protoreflect.MessageKind {
			_, varintLength := protowire.ConsumeVarint(value)
			if varintLength < 0 {
				return nil, protowire.ParseError(varintLength)
			}
			nested := value[typLength+varintLength:]
			if len(nested) > 0 {
				children, err := splitPb(nested, fd.Message(), originalMd, path)
				if err != nil {
					return nil, err
				}
				if fd.IsMap() {
					// Special treatment is needed for maps. Its key/value pairs are encoded as repeated nested
					// messages, each with two fields for the key and value. But, to allow manipulation and addressing
					// of values by their keys, we want the value to represented as a 'top level' Part, with the
					// key forming part of the path.
					//
					// Thus, we need to manipulate the values' paths: at this point they look like  ….parent.2.…, but
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
					}
					children = values // Toss the key as a dedicated element now it's part of the path of each value
				}
				v = append(v, children...)
				continue
			}
		}

		v = append(v, Part{
			Path:  path,
			Bytes: value,
			Md:    originalMd,
		})
	}

	return v, nil
}

// Split explodes a serialised Protocol Buffer into parts of its constituent fields and those within nested
// messages. The resulting parts are returned in the order they appear in pb – so recombining the parts (unless they
// are reordered either manually or by sorting them) yields a byte-wise identical message.
func Split(pb []byte, md protoreflect.MessageDescriptor) (Parts, error) {
	return splitPb(pb, md, md, nil)
}
