package protoparts

import (
	"sort"

	"google.golang.org/protobuf/encoding/protowire"
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

// MergeProtoParts merges p1 and p2, with items in p2 taking precedence over p1. The resulting Parts are sorted.
func MergeProtoParts(p1, p2 Parts) Parts {
	// Make a copy of p1, into which we will merge p2
	v := make(Parts, len(p1))
	copy(v, p1)
	for _, part := range p2 {
		// Remove any existing part from v which is prefixed with the new part's path. Preserving order is
		// unimportant as we will sort below.
		indicesToRemove := v.selectIndices(part.Path)
		for ii := len(indicesToRemove) - 1; ii >= 0; ii-- {
			i := indicesToRemove[ii]
			v[i] = v[len(v)-1] // https://github.com/golang/go/wiki/SliceTricks#delete-without-preserving-order
			v[len(v)-1] = Part{}
			v = v[:len(v)-1]
		}
		// Add the new part to the end
		v = append(v, part)
	}
	// The new list may now have a completely nonsensical ordering, so sort
	sort.Sort(v)
	return v
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
	return ps.join(0)
}

func (ps Parts) join(prefixLen int) []byte {
	b := bufferPool.Get()
	defer bufferPool.Put(b)

	// wrapNested wraps the passed buffer with the tag + length prefix needed to denote that it's a nested message
	wrapNested := func(bb []byte, tag protowire.Number) []byte {
		envelope := make([]byte, 0, protowire.SizeTag(tag)+protowire.SizeBytes(len(bb)))
		envelope = protowire.AppendTag(envelope, tag, protowire.BytesType)
		envelope = protowire.AppendVarint(envelope, uint64(len(bb)))
		return append(envelope, bb...)
	}

	parts := ps
	for len(parts) > 0 {
		part := parts[0]
		bb := part.Bytes
		prefix, run := part.Path[:prefixLen+1], parts
		for ii, candidate := range run { // extract all descendants of the prefix to serialise together
			if !candidate.Path.HasPrefix(prefix) {
				run = parts[:ii]
				break
			}
		}
		if len(run) > 1 {
			bb = run.join(prefixLen + 1)
		}

		if part.Path[prefixLen].Key != nil {
			// We have a map value, so we need to wrap everything in an enclosing entry message with key, value fields
			key := part.Path[prefixLen].Key
			if len(part.Path) > prefixLen+1 {
				// If the value is itself a nested message, it needs a tag + length prefix
				bb = wrapNested(bb, 2)
			}
			entry := append(key, bb...)
			bb = wrapNested(entry, prefix.Last().Tag)
		} else if len(part.Path) > prefixLen+1 {
			// We have nested elements, so we need to stick on a tag + length prefix
			bb = wrapNested(bb, prefix.Last().Tag)
		}

		b.Write(bb)
		parts = parts[len(run):]
		continue // skip over all the now-marshaled descendants
	}

	// We need to copy the bytes of the buffer; otherwise they can/will get reused
	result := make([]byte, b.Len())
	copy(result, b.Bytes())
	return result
}

func (ps Parts) Len() int           { return len(ps) }
func (ps Parts) Swap(x, y int)      { ps[x], ps[y] = ps[y], ps[x] }
func (ps Parts) Less(x, y int) bool { return ps[x].Path.Compare(ps[y].Path) == -1 }