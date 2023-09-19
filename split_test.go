package protoparts

import (
	"fmt"
	"testing"

	"github.com/obeattie/protoparts/internal/testproto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/dynamicpb"
)

func TestSplit(t *testing.T) {
	md := (&testproto.Person{}).ProtoReflect().Descriptor()
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
			testMsg(t, p("Oliver Beattie"), nil, nil, nil, nil, nil),
			[]Path{
				DecodeSymbolicPath("name", md)},
		},
		{
			// An optional field with empty contents is still present
			testMsg(t, p(""), nil, nil, nil, nil, nil),
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

func BenchmarkSplit(b *testing.B) {
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
