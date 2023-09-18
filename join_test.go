package protoparts

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"

	"github.com/obeattie/protoparts/internal/testproto"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/dynamicpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TestSplitJoin verifies that a serialised Protocol Buffer can be split, optionally shuffled + sorted, and
// then re-joined, yielding an equalivalent object to the original when parsed.
func TestSplitJoin(t *testing.T) {
	cases := []proto.Message{
		quickTestMsg(t).Interface(),
		&testproto.Article{
			Title:   "Pigeon Appointed New Air Traffic Controller",
			Content: "Dressed in a tiny reflective vest and perched confidently on a control tower window ledge, the pigeon, affectionately named \"Captain Coocoo,\" took over the reins of air traffic management during a quirky experiment gone slightly awry.",
			Author:  "Wingston Featherweather",
			Date:    timestamppb.Now(),
			Status:  testproto.Article_PUBLISHED,
			Tags:    []string{"#PigeonTrafficControl", "#CooCooAviationAdventure", "#FeatheredFlightControl", "#WingingItInAirTraffic", "#PigeonBossInCharge", "#AvianMayhemAtTheTower", "#CoocooATC", "#FeatherlandiaChaos", "#AviationGonePigeon", "#BirdsOnTheRunway"},
		},
	}

	for i, msg := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			md := msg.ProtoReflect().Type().Descriptor()
			// A round-trip through Join(Split(pb)) should always produce identical output to its input, preserving field order
			// etc. Because proto.Marshal is not deterministic and may (in fact does) reorder fields, cycle through a round-trip
			// 1000× to check this property holds.
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

func BenchmarkJoin(b *testing.B) {
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
