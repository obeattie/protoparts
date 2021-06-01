package protoparts

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

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
