package protoparts

import (
	"math/rand"
	"testing"

	"github.com/obeattie/protoparts/testproto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protowire"
)

func TestDecodePath(t *testing.T) {
	for i := 0; i < 1000; i++ {
		depth := rand.Intn(100)
		vp := make(Path, depth)
		for i := 0; i < depth; i++ {
			vp[i].Tag = protowire.Number(rand.Int31())
			if rand.Float32() > 0.5 {
				vp[i].Index = rand.Int()
			}
			if rand.Float32() > 0.5 {
				key := make([]byte, rand.Intn(100))
				rand.Read(key)
				vp[i].Key = key
			}
		}
		vp2 := DecodePath(vp.String())
		assert.True(t, vp.Equal(vp2), `path: "%v"`, vp)
		assert.Equal(t, vp.String(), vp2.String())
	}
}

func TestDecodeSymbolicPath(t *testing.T) {
	md := (&testproto.Person{}).ProtoReflect().Descriptor()
	type tc struct {
		string
		Path
	}
	cases := []tc{
		{"address", Path{{2, -1, nil}}},
		{"address/city", Path{{2, -1, nil}, {2, -1, nil}}},
		{"moar_addresses", Path{{3, -1, nil}}},
		{"moar_addresses[2]", Path{{3, 2, nil}}},
		{"moar_addresses[2]/city", Path{{3, 2, nil}, {2, -1, nil}}},
		{"map_string_string[0x68656c6c6f]", Path{{10, -1, []byte("hello")}}},
		{"map_string_string[2][0x68656c6c6f]", Path{{10, 2, []byte("hello")}}},
		{"address.city", nil}, // invalid paths = nil
		{"boopboopboop", nil},
		{"map_string_string[0x68656c6c6f][2]", nil},
	}
	for _, c := range cases {
		t.Run(c.string, func(t *testing.T) {
			expected := c.Path
			actual := DecodeSymbolicPath(c.string, md)
			assert.Equal(t, expected, actual)
			if actual != nil {
				ss, err := actual.SymbolicString(md)
				require.NoError(t, err)
				assert.Equal(t, c.string, ss)
			}
		})
	}
}

func TestPathHasPrefix(t *testing.T) {
	type tc struct {
		p      string
		prefix string
		bool
	}
	cases := []tc{
		{"1", "1", true},
		{"1", "2", false},
		{"1/2", "1/2", true},
		{"1/2", "1/2/3", false},
		{"1/2", "1/1", false},
		{"1/2", "1", true},
		{"1/2", "2", false},
		{"1[1]", "1[1]", true},
		{"1[1]", "1", true},
		{"1[1]", "2", false},
		{"1[1]", "1[2]", false},
		{"1[0x0a]", "1[0x0a]", true},
		{"1[0x0a]", "1[0x0b]", false},
		{"1[0x0a]", "1", true},
		{"", "", true},
		{"1", "", true},
	}
	for _, c := range cases {
		t.Run(c.p, func(t *testing.T) {
			p, prefix := DecodePath(c.p), DecodePath(c.prefix)
			actual := p.HasPrefix(prefix)
			assert.Equal(t, c.bool, actual, "(%v).HasPrefix(%v) ≠ %v", c.p, c.prefix, c.bool)
			assert.False(t, p.HasPrefix(nil), "(%v).HasPrefix(nil) should return false", c.p)
		})
	}
}
