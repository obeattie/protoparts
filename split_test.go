package protoparts

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/dynamicpb"
)

func TestSplitJoinProto(t *testing.T) {
	person := quickTestMsg(t)
	personMd := person.Type().Descriptor()

	// A round-trip through Join(Split(pb)) should always produce identical output to its input, preserving field order
	// etc. Because proto.Marshal is not deterministic and may (in fact does) reorder fields, cycle through a round-trip
	// 1000x to check this property holds.
	for i := 0; i < 1000; i++ {
		pb, err := proto.Marshal(person)
		require.NoError(t, err)
		parts, err := Split(pb, personMd)
		require.NoError(t, err)
		// if i == 0 {
		// 	t.Logf("%x", pb)
		// 	t.Log(parts)
		// }
		pb2 := parts.Join()
		require.Equal(t, pb, pb2)
	}

	// Now, test that if the parts are shuffled, they are still reassembled in a way that yields an equivalent message
	// Note: this does not mean the output is byte-wise identical to the input as we assert above – we can't expect
	// that property to hold since proto.Marshal can reorder things arbitrarily – simply that they can be unmarshaled
	// successfully are considered equal by proto.Equal.
	for i := 0; i < 100; i++ {
		ps := split(t, person)
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
		person2 := dynamicpb.NewMessage(person.Descriptor())
		require.NoError(t, proto.Unmarshal(pb2, person2))
		require.True(t, proto.Equal(person, person2))
	}
}

func BenchmarkSplitProto(b *testing.B) {
	msg := benchmarkMsg(b)
	md := msg.Descriptor()
	pb, err := proto.Marshal(msg)
	require.NoError(b, err)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Split(pb, md)
	}
}

func BenchmarkProtoPartsSort(b *testing.B) {
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

func BenchmarkProtoPartsJoin(b *testing.B) {
	msg := benchmarkMsg(b)
	md := msg.Descriptor()
	pb, err := proto.Marshal(msg)
	require.NoError(b, err)
	ps, err := Split(pb, md)
	require.NoError(b, err)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ps.Join()
	}
}
