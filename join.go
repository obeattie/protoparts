package protoparts

import (
	"slices"

	"google.golang.org/protobuf/encoding/protowire"
)

// msgPrefix returns the tag+length prefix needed to signal a nested message
func msgPrefix(tag protowire.Number, length int) []byte {
	b := make([]byte, 0, protowire.SizeTag(tag)+protowire.SizeVarint(uint64(length)))
	b = protowire.AppendTag(b, tag, protowire.BytesType)
	b = protowire.AppendVarint(b, uint64(length))
	return b
}

func join(ps Parts, prefixLen int) []byte {
	buf := bufferPool.Get()
	defer bufferPool.Put(buf)

	parts := ps
	for len(parts) > 0 {
		part := parts[0]
		isMapVal := part.Path[prefixLen].Key != nil
		isNestedMsg := len(part.Path) > prefixLen+1

		// Extract all descendants of the prefix to serialise together
		prefix, run := part.Path[:prefixLen+1], parts
		for ii, candidate := range run[1:] {
			if len(candidate.Path) <= len(prefix) || !candidate.Path.HasPrefix(prefix) {
				run = run[:ii+1] // +1 because the [1:] above means the range indices are off by one
				break
			}
		}

		bb := part.Bytes
		if len(run) > 1 {
			bb = join(run, prefixLen+1)
		} else if isNestedMsg || !isMapVal {
			// Prepend field tag
			c := protowire.AppendTag(nil, part.Path.Last().Tag, part.Type)
			bb = append(c, bb...)
		}

		if isMapVal {
			// We have a map value, so we need to wrap everything in an enclosing entry message with key, value fields
			entry := protowire.AppendTag(nil, 1, part.KeyType)
			entry = append(entry, part.Path[prefixLen].Key...)
			if isNestedMsg {
				entry = append(entry, msgPrefix(2, len(bb))...)
			} else {
				entry = protowire.AppendTag(entry, 2, part.Type)
			}
			entry = append(entry, bb...)
			bb = msgPrefix(prefix.Last().Tag, len(entry))
			bb = append(bb, entry...)
		} else if isNestedMsg {
			// We have nested elements, so we need to stick on a tag + length prefix
			bb = append(msgPrefix(prefix.Last().Tag, len(bb)), bb...)
		}

		buf.Write(bb)
		parts = parts[len(run):]
	}

	// We need to copy the bytes of the buffer; otherwise they can/will get reused
	return slices.Clone(buf.Bytes())
}
