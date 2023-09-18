package protoparts

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

// Parts are the combination of all the (potentially nested) fields within a serialised Protocol Buffer message.
//
// Fields are ordered, and the order is important when encoding because protobuf can only decode messages when, for
// example, repeated elements are encoded in-order, the fields of nested messages are serialised together in a 'run'
// (but not necessarily in an order within that run), etc. This means there can be many valid encoding orders for a
// single message, but not all conceivable orders are valid. To always ensure a valid serialisation will result,
// Parts can be sorted to produce a canonical order. However, this may unnecessarily change an already-valid
// field order, so it is unnecessary to sort if the order is known to be good – for example if the Parts were
// produced by Split.
type Parts []Part

// Paths returns all the contained field Paths
func (ps Parts) Paths() []Path {
	return mapSlice(ps, func(p Part) Path {
		return p.Path
	})
}

// Value returns the Protocol Buffer for the given Path, or nil if it could not be found which is possible if:
//
// - the path is invalid within the Protocol Buffer
//
// - the field was not set within the original Protocol Buffer
//
// - the field is contained within a sub-message, and no message descriptor was provided
func (ps Parts) Value(p Path) []byte {
	return ps.Select(p).Join()
}

// ProtoValue returns (protoreflect.Value, true) for the given Path, or (protoreflect.Value{}, false) if it could not
// be found or unmarshalled.
func (ps Parts) ProtoValue(p Path) (protoreflect.Value, bool) {
	selected := ps.Select(p)
	if len(selected) == 0 {
		return protoreflect.Value{}, false
	}
	pb := selected.Join()
	md := selected[0].Md
	// The index (if there is one) will always be 0 in pb since Select() eliminates all other siblings
	if len(p) > 0 && p.Last().Index != -1 {
		p[len(p)-1].Index = 0
	}
	msg := dynamicpb.NewMessage(md)
	if err := proto.Unmarshal(pb, msg); err != nil {
		return protoreflect.Value{}, false
	}
	return valueFromMessage(msg, p)
}

// ProtoMessage returns the unmarshalled protoreflect.Message for the entire Parts assembly. If the message cannot
// be unmarshalled, returns nil.
func (ps Parts) ProtoMessage() protoreflect.Message {
	pv, ok := ps.ProtoValue(Path{})
	if !ok {
		return nil
	}
	return pv.Message()
}

// selectIndices returns the (sorted) indices of all parts that are prefixed with any of the passed paths
func (ps Parts) selectIndices(prefixes ...Path) []int {
	var indices []int
	// @TODO: An index could speed this up, as could a binary search if we know the parts are sorted
	for i, candidate := range ps {
		for _, prefix := range prefixes {
			if candidate.Path.HasPrefix(prefix) {
				indices = append(indices, i)
				break
			}
		}
	}
	return indices
}

// Select returns all the parts that are prefixed with any of the passed paths. The order of selected parts will be the
// same as their order before selection.
func (ps Parts) Select(prefixes ...Path) Parts {
	indices := ps.selectIndices(prefixes...)
	selected := make(Parts, len(indices))
	for i, ii := range indices {
		selected[i] = ps[ii]
	}
	return selected
}

// Exclude returns all the parts that are not prefixed with any of the passed paths. The order of the resulting parts
// will be the same as their order before exclusion.
func (ps Parts) Exclude(prefixes ...Path) Parts {
	indices := ps.selectIndices(prefixes...)
	remaining := make(Parts, 0, len(ps)-len(indices))
	// To build a new list of Parts that excludes the indexes returned by selectIndices, we:
	//
	// - maintain a cursor of where we have iterated to in ps
	// - iterate over the slice of _indices_ and copy all of the items between the cursor and the excluded index
	// - set the cursor to just past the excluded item on each iteration
	// - finally, copy the included bits past the end of the excluded indices
	//
	// This slightly convoluted process means this can be done in one pass over the list, and without unnecessarily
	// copying items to the new list we'll later discard.
	cursor := 0
	for _, ii := range indices {
		for i := cursor; i < ii; i++ {
			remaining = append(remaining, ps[i])
		}
		cursor = ii + 1
	}
	for i := cursor; i < len(ps); i++ {
		remaining = append(remaining, ps[i])
	}
	return remaining
}

// Join stitches the parts back together as a serialised Protocol Buffer message.
//
// Warning: If the order of fields is not valid, Join() may produce malformed output which cannot be parsed as a
// Protocol Buffer. If the order of parts is not known to be valid, sort them first.
func (ps Parts) Join() []byte {
	return join(ps, 0)
}

func (ps Parts) Len() int           { return len(ps) }
func (ps Parts) Swap(x, y int)      { ps[x], ps[y] = ps[y], ps[x] }
func (ps Parts) Less(x, y int) bool { return ps[x].Path.Compare(ps[y].Path) == -1 }
