package protoparts

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/dynamicpb"

	"github.com/obeattie/protoparts/internal/testproto"
)

// TestProject verifies that we can pull individual Parts out of the Protocol Buffer and unmarshal them without issue.
func TestProject(t *testing.T) {
	msg := quickTestMsg(t)
	md := msg.ProtoReflect().Type().Descriptor()
	parts := split(t, msg)

	// Pull all the pieces out individually
	for _, p := range parts {
		fieldMsg := dynamicpb.NewMessage(md)
		require.NoError(t, proto.Unmarshal(parts.Value(p.Path), fieldMsg), p.Path.String())
	}

	// Now, pull random combinations of the pieces out (by shuffling then taking a random number of parts from the
	// head) and check that any combination
	for i := 0; i < 5000; i++ {
		rand.Shuffle(len(parts), func(i, j int) { parts[i], parts[j] = parts[j], parts[i] })
		head := parts[:rand.Intn(len(parts))]
		sort.Sort(head) // ensure there's a valid order (which may be different to the original order)
		headMsg := dynamicpb.NewMessage(md)
		require.NoError(t, proto.Unmarshal(head.Join(), headMsg))
	}
}

// TestRearrange asserts that we can manipulate the Paths of Parts, and when they are joined get a message that's
// identical to what we get if we construct a protobuf in the target state directly.
func TestRearrange(t *testing.T) {
	msg := &testproto.Address{
		StreetAddress: "foobar",
	}
	parts, err := Marshal(msg)
	require.NoError(t, err)
	parts[0].Path[0].Tag = 2 // this changes the "foobar" to be in the city field, not the street address field

	expected := marshalProto(t, &testproto.Address{
		City: "foobar",
	})
	assert.Equal(t, expected, parts.Join())
}

func BenchmarkSort(b *testing.B) {
	msg := benchmarkMsg(b)
	md := msg.Descriptor()
	pb, err := proto.Marshal(msg)
	require.NoError(b, err)
	ps, err := Split(pb, md)
	require.NoError(b, err)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sort.Sort(ps)
	}
}
