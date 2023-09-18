package protoparts

import (
	"slices"
	"sort"
)

// Merge merges p1 and p2, with items in p2 taking precedence over p1. The resulting Parts are sorted.
func Merge(p1, p2 Parts) Parts {
	// Make a copy of p1, into which we will merge p2
	v := slices.Clone(p1)
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
