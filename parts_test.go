package protoparts

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/dynamicpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	ppproto "github.com/obeattie/protoparts/test/proto"
)

func TestSplitProto(t *testing.T) {
	md := (&ppproto.Person{}).ProtoReflect().Descriptor()
	type tc struct {
		msg           *dynamicpb.Message
		expectedPaths []Path
	}
	cases := []tc{
		{
			// Empty message
			testMsg(t, nil, nil, nil, nil, nil, nil),
			[]Path{},
		},
		{
			// Simple, one field populated
			testMsg(t, s("Oliver Beattie"), nil, nil, nil, nil, nil),
			[]Path{
				DecodeSymbolicPath("name", md)},
		},
		{
			// An optional field with empty contents is still present
			testMsg(t, s(""), nil, nil, nil, nil, nil),
			[]Path{
				DecodeSymbolicPath("name", md)},
		},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			parts := split(t, c.msg)
			var seenPaths []Path
			for _, part := range parts {
				seenPaths = append(seenPaths, part.Path)
			}
			assert.ElementsMatch(t, c.expectedPaths, seenPaths)
		})
	}
}

func TestMergeProtoParts(t *testing.T) {
	tc := [][3]*dynamicpb.Message{
		{
			// Top-level no-op
			testMsg(t, s("Oliver Beattie"), nil, nil, nil, nil, nil),
			testMsg(t, nil, nil, nil, nil, nil, nil),
			testMsg(t, s("Oliver Beattie"), nil, nil, nil, nil, nil)},
		{
			// Top-level field replacement
			testMsg(t, s("Oliver Beattie"), nil, nil, nil, nil, nil),
			testMsg(t, s("Lindy Bishop"), nil, nil, nil, nil, nil),
			testMsg(t, s("Lindy Bishop"), nil, nil, nil, nil, nil)},
		{
			// Top-level field zeroing, with a field that has explicit presence
			testMsg(t, s("Oliver Beattie"), nil, nil, nil, nil, nil),
			testMsg(t, s(""), nil, nil, nil, nil, nil),
			testMsg(t, s(""), nil, nil, nil, nil, nil)},
		{
			// Nested field replacement
			testMsg(t, nil, nil, s("London"), nil, nil, nil),
			testMsg(t, nil, nil, s("Los Angeles"), nil, nil, nil),
			testMsg(t, nil, nil, s("Los Angeles"), nil, nil, nil)},
		{
			// Nested no-op
			testMsg(t, nil, nil, s("London"), nil, nil, nil),
			testMsg(t, nil, nil, nil, nil, nil, nil),
			testMsg(t, nil, nil, s("London"), nil, nil, nil)},
		{
			// List no-op
			testMsg(t, nil, nil, nil, []string{"a"}, nil, nil),
			testMsg(t, nil, nil, nil, nil, nil, nil),
			testMsg(t, nil, nil, nil, []string{"a"}, nil, nil)},
		{
			// List no-op
			testMsg(t, nil, nil, nil, []string{"a"}, nil, nil),
			// From the Protobuf docs: "Repeated fields and maps do not track presence: there is no distinction between
			// an empty and a not-present repeated field." So, we expect an empty list to the be treated the same as nil
			testMsg(t, nil, nil, nil, []string{}, nil, nil),
			testMsg(t, nil, nil, nil, []string{"a"}, nil, nil)},
		{
			// List replacement (second list is longer than the first)
			testMsg(t, nil, nil, nil, []string{"a"}, nil, nil),
			testMsg(t, nil, nil, nil, []string{"b", "c"}, nil, nil),
			testMsg(t, nil, nil, nil, []string{"b", "c"}, nil, nil)},
		{
			// Partial list replacement (second list is shorter than the first)
			testMsg(t, nil, nil, nil, []string{"a", "b", "c"}, nil, nil),
			testMsg(t, nil, nil, nil, []string{"aa"}, nil, nil),
			testMsg(t, nil, nil, nil, []string{"aa", "b", "c"}, nil, nil)},
		{
			// Adding a key to a map
			testMsg(t, nil, nil, nil, nil, nil, map[string]string{"a": "a", "b": "b"}),
			testMsg(t, nil, nil, nil, nil, nil, map[string]string{"c": "c"}),
			testMsg(t, nil, nil, nil, nil, nil, map[string]string{"a": "a", "b": "b", "c": "c"})},
		{
			// Overwriting a key in a map
			testMsg(t, nil, nil, nil, nil, nil, map[string]string{"a": "a", "b": "b"}),
			testMsg(t, nil, nil, nil, nil, nil, map[string]string{"a": "🤑"}),
			testMsg(t, nil, nil, nil, nil, nil, map[string]string{"a": "🤑", "b": "b"})},
		{
			// Overwriting + adding keys to a map
			testMsg(t, nil, nil, nil, nil, nil, map[string]string{"a": "a", "b": "b", "c": "c"}),
			testMsg(t, nil, nil, nil, nil, nil, map[string]string{"a": "🤑", "d": "d"}),
			testMsg(t, nil, nil, nil, nil, nil, map[string]string{"a": "🤑", "b": "b", "c": "c", "d": "d"})},
	}

	for i, c := range tc {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			before, mutation, expected := c[0], c[1], c[2]
			beforeParts, mutationParts := split(t, before), split(t, mutation)
			expectedB := marshalProto(t, expected)
			afterParts := MergeProtoParts(beforeParts, mutationParts)
			after := afterParts.ProtoMessage()
			require.NotNil(t, after)
			afterB := marshalProto(t, after.Interface())
			if !assert.Equal(t, expectedB, afterB) {
				t.Logf("          == %d detail", i)
				t.Logf("  Before: %v", beforeParts)
				mutationB := marshalProto(t, mutation)
				t.Logf("Mutation: %v", mutationParts)
				t.Logf(" (bytes): %x", mutationB)
				t.Logf("   (msg): %v", mutation)
				t.Logf("   After: %v", afterParts)
			}
		})
	}
}

// TestSplitJoinProto verifies that a serialised Protocol Buffer can be split, optionally shuffled + sorted, and
// then re-joined, yielding an equalivalent object to the original when parsed.
func TestSplitJoinProto(t *testing.T) {
	cases := []proto.Message{
		quickTestMsg(t).Interface(),
		&ppproto.Article{
			Title:   "Pigeon Appointed New Air Traffic Controller",
			Content: "Dressed in a tiny reflective vest and perched confidently on a control tower window ledge, the pigeon, affectionately named \"Captain Coocoo,\" took over the reins of air traffic management during a quirky experiment gone slightly awry.",
			Author:  "Wingston Featherweather",
			Date:    timestamppb.Now(),
			Status:  ppproto.Article_PUBLISHED,
			Tags:    []string{"#PigeonTrafficControl", "#CooCooAviationAdventure", "#FeatheredFlightControl", "#WingingItInAirTraffic", "#PigeonBossInCharge", "#AvianMayhemAtTheTower", "#CoocooATC", "#FeatherlandiaChaos", "#AviationGonePigeon", "#BirdsOnTheRunway"},
		},
	}

	for i, msg := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			md := msg.ProtoReflect().Type().Descriptor()
			// A round-trip through Join(Split(pb)) should always produce identical output to its input, preserving field order
			// etc. Because proto.Marshal is not deterministic and may (in fact does) reorder fields, cycle through a round-trip
			// 1000x to check this property holds.
			for i := 0; i < 1000; i++ {
				pb, err := proto.Marshal(msg)
				require.NoError(t, err)
				parts, err := Split(pb, md)
				require.NoError(t, err)
				pb2 := parts.Join()
				require.Equal(t, pb, pb2)
			}

			// Now, test that if the parts are shuffled, they are still reassembled in a way that yields an equivalent message
			// Note: this does not mean the output is byte-wise identical to the input as we assert above – we can't expect
			// that property to hold since proto.Marshal can reorder things arbitrarily – simply that they can be unmarshalled
			// successfully are considered equal by proto.Equal.
			for i := 0; i < 100; i++ {
				ps := split(t, msg.ProtoReflect())
				rand.Shuffle(len(ps), func(i, j int) {
					ps[i], ps[j] = ps[j], ps[i]
				})
				// if i == 0 {
				// 	t.Logf("unsorted: %v", ps)
				// }
				sort.Sort(ps) // Since we shuffled, this is crucial to get a well-formed result
				// if i == 0 {
				// 	t.Logf("  sorted: %v", ps)
				// }
				pb2 := ps.Join()
				msg2 := dynamicpb.NewMessage(md)
				require.NoError(t, proto.Unmarshal(pb2, msg2))
				require.True(t, proto.Equal(msg, msg2))
			}
		})
	}

}
