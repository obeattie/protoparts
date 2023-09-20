package protoparts

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
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

func TestProtoValue(t *testing.T) {
	msg := &testproto.Person{
		Name: p("Ryan Gosling"),
		Address: &testproto.Address{
			StreetAddress: "3532 Hayden Ave",
			City:          "Culver City",
		},
		Boop:          [][]byte{[]byte("bop"), []byte("beat")},
		MaritalStatus: testproto.Person_MARRIED,
		MapStringString: map[string]string{
			"a": "b",
		},
		MapStringLatlng: map[string]*testproto.LatLng{
			"london": {
				Latitude:  51.508125,
				Longitude: -0.128081,
			},
			"berlin": {
				Latitude:  52.511100,
				Longitude: 13.442989,
			},
		},
	}
	parts := split(t, msg.ProtoReflect())
	expectations := map[string]any{
		"name":                      "Ryan Gosling",
		"address/street_address":    "3532 Hayden Ave",
		"map_string_string[0x0161]": "b",
		"map_string_latlng[0x066c6f6e646f6e]/latitude":  51.508125,
		"map_string_latlng[0x066c6f6e646f6e]/longitude": -0.128081,
		"map_string_latlng[0x066265726c696e]/latitude":  52.511100,
		"map_string_latlng[0x066265726c696e]/longitude": 13.442989,
		"marital_status": protoreflect.EnumNumber(4),
		"boop[0]":        []byte("bop"),
		"boop[1]":        []byte("beat"),
	}

	for symbolicPath, expectedVal := range expectations {
		path := DecodeSymbolicPath(symbolicPath, msg.ProtoReflect().Descriptor())
		v, ok := parts.ProtoValue(path)
		if assert.True(t, ok) {
			assert.Equal(t, expectedVal, v.Interface())
		}
	}
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
