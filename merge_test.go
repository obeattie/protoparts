package protoparts

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/dynamicpb"
)

func TestMerge(t *testing.T) {
	tc := [][3]*dynamicpb.Message{
		{
			// 0: Top-level no-op
			testMsg(t, s("Oliver Beattie"), nil, nil, nil, nil, nil),
			testMsg(t, nil, nil, nil, nil, nil, nil),
			testMsg(t, s("Oliver Beattie"), nil, nil, nil, nil, nil)},
		{
			// 1: Top-level field replacement
			testMsg(t, s("Oliver Beattie"), nil, nil, nil, nil, nil),
			testMsg(t, s("Lindy Bishop"), nil, nil, nil, nil, nil),
			testMsg(t, s("Lindy Bishop"), nil, nil, nil, nil, nil)},
		{
			// 2: Top-level field zeroing, with a field that has explicit presence
			testMsg(t, s("Oliver Beattie"), nil, nil, nil, nil, nil),
			testMsg(t, s(""), nil, nil, nil, nil, nil),
			testMsg(t, s(""), nil, nil, nil, nil, nil)},
		{
			// 3: Nested field replacement
			testMsg(t, nil, nil, s("London"), nil, nil, nil),
			testMsg(t, nil, nil, s("Los Angeles"), nil, nil, nil),
			testMsg(t, nil, nil, s("Los Angeles"), nil, nil, nil)},
		{
			// 4: Nested no-op
			testMsg(t, nil, nil, s("London"), nil, nil, nil),
			testMsg(t, nil, nil, nil, nil, nil, nil),
			testMsg(t, nil, nil, s("London"), nil, nil, nil)},
		{
			// 5: List no-op
			testMsg(t, nil, nil, nil, []string{"a"}, nil, nil),
			testMsg(t, nil, nil, nil, nil, nil, nil),
			testMsg(t, nil, nil, nil, []string{"a"}, nil, nil)},
		{
			// 6: List no-op
			testMsg(t, nil, nil, nil, []string{"a"}, nil, nil),
			// From the Protobuf docs: "Repeated fields and maps do not track presence: there is no distinction between
			// an empty and a not-present repeated field." So, we expect an empty list to the be treated the same as nil
			testMsg(t, nil, nil, nil, []string{}, nil, nil),
			testMsg(t, nil, nil, nil, []string{"a"}, nil, nil)},
		{
			// 7: List replacement (second list is longer than the first)
			testMsg(t, nil, nil, nil, []string{"a"}, nil, nil),
			testMsg(t, nil, nil, nil, []string{"b", "c"}, nil, nil),
			testMsg(t, nil, nil, nil, []string{"b", "c"}, nil, nil)},
		{
			// 8: Partial list replacement (second list is shorter than the first)
			testMsg(t, nil, nil, nil, []string{"a", "b", "c"}, nil, nil),
			testMsg(t, nil, nil, nil, []string{"aa"}, nil, nil),
			testMsg(t, nil, nil, nil, []string{"aa", "b", "c"}, nil, nil)},
		{
			// 9: Adding a key to a map
			testMsg(t, nil, nil, nil, nil, nil, map[string]string{"a": "a", "b": "b"}),
			testMsg(t, nil, nil, nil, nil, nil, map[string]string{"c": "c"}),
			testMsg(t, nil, nil, nil, nil, nil, map[string]string{"a": "a", "b": "b", "c": "c"})},
		{
			// 10: Overwriting a key in a map
			testMsg(t, nil, nil, nil, nil, nil, map[string]string{"a": "a", "b": "b"}),
			testMsg(t, nil, nil, nil, nil, nil, map[string]string{"a": "🤑"}),
			testMsg(t, nil, nil, nil, nil, nil, map[string]string{"a": "🤑", "b": "b"})},
		{
			// 11: Overwriting + adding keys to a map
			testMsg(t, nil, nil, nil, nil, nil, map[string]string{"a": "a", "b": "b", "c": "c"}),
			testMsg(t, nil, nil, nil, nil, nil, map[string]string{"a": "🤑", "d": "d"}),
			testMsg(t, nil, nil, nil, nil, nil, map[string]string{"a": "🤑", "b": "b", "c": "c", "d": "d"})},
	}

	for i, c := range tc {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			before, mutation, expected := c[0], c[1], c[2]
			beforeParts, mutationParts := split(t, before), split(t, mutation)
			expectedB := marshalProto(t, expected)
			afterParts := Merge(beforeParts, mutationParts)
			after := afterParts.ProtoMessage()
			require.NotNil(t, after)
			afterB := marshalProto(t, after.Interface())
			if !assert.Equal(t, expectedB, afterB) {
				t.Logf("       == %d detail", i)
				t.Logf("  Before: %v", beforeParts)
				t.Logf(" (bytes):\n%s", hex.Dump((marshalProto(t, before))))
				t.Logf("Mutation: %v", mutationParts)
				t.Logf(" (bytes):\n%s", hex.Dump(marshalProto(t, mutation)))
				t.Logf("   (msg): %v", mutation)
				t.Logf("Expected: %v", split(t, expected))
				t.Logf(" (bytes):\n%s", hex.Dump(marshalProto(t, expected)))
				t.Logf("  Actual: %v", afterParts)
				t.Logf(" (bytes):\n%s", hex.Dump(afterB))
			}
		})
	}
}
