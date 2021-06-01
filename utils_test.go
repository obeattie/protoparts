package protoparts

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/reflect/protoreflect"

	ppproto "github.com/obeattie/protoparts/test/proto"
)

func TestValueFromMessage(t *testing.T) {
	m := testMsg(t, s("Oliver Beattie"), s("203 Chautauqua Blvd"), nil, nil, nil, nil)
	md := m.Descriptor()

	emptyMap := m.NewField(md.Fields().ByName("map_string_string")).Map()
	emptyList := m.NewField(md.Fields().ByName("tags")).List()

	type tc struct {
		path  string
		value protoreflect.Value
		ok    bool
	}
	cases := []tc{
		// Happy paths
		{`name`, protoreflect.ValueOfString("Oliver Beattie"), true},
		{`address/street_address`, protoreflect.ValueOfString("203 Chautauqua Blvd"), true},
		// Valid path but unset = empty value, not present
		{`address/city`, protoreflect.ValueOfString(""), false},
		{`map_string_string`, protoreflect.ValueOfMap(emptyMap), false},
		{`tags`, protoreflect.ValueOfList(emptyList), false},
		// This one is a bit weird – I naïvely expected this to be an empty string, as it would be if it was an unset
		// string value. But the Protobuf libraries seem to return a nil value for an absent map key.
		{`map_string_string[0x0a]`, protoreflect.Value{}, false},
		// Invalid paths = nil value, not present
		{`garbage`, protoreflect.Value{}, false},
	}

	for _, c := range cases {
		t.Run(c.path, func(t *testing.T) {
			p := DecodeSymbolicPath(c.path, md)
			vv, ok := valueFromMessage(m, p)
			assert.Equal(t, c.ok, ok)
			v := vv.Interface()
			switch v := v.(type) {
			case protoreflect.List:
				expectedList := c.value.List()
				if assert.Equal(t, expectedList.Len(), v.Len()) {
					for i := 0; i < expectedList.Len(); i++ {
						assert.Equal(t, expectedList.Get(i), v.Get(i))
					}
				}
			case protoreflect.Map:
				expectedMap := c.value.Map()
				if assert.Equal(t, expectedMap.Len(), v.Len()) {
					expectedMap.Range(func(mk protoreflect.MapKey, mv protoreflect.Value) bool {
						return assert.Equal(t, mv, v.Get(mk))
					})
				}
			default:
				assert.Equal(t, c.value.Interface(), v)
			}
		})
	}
}

func TestFieldDescriptorInMessage(t *testing.T) {
	md := (&ppproto.Person{}).ProtoReflect().Descriptor()
	type tc struct {
		p  string
		fd protoreflect.FieldDescriptor
	}
	cases := []tc{
		{`name`, md.Fields().ByName(`name`)},
		{`name[2]`, md.Fields().ByName(`name`)},
		{`address/city`, md.Fields().ByName(`address`).Message().Fields().ByName(`city`)},
		{`moar_addresses/city`, md.Fields().ByName(`moar_addresses`).Message().Fields().ByName(`city`)},
		{`map_string_string`, md.Fields().ByName(`map_string_string`)},
		{`map_string_string[0x0a]`, md.Fields().ByName(`map_string_string`).MapValue()},
		{`address.city`, nil}, // invalid path = nil
		{`boopboopboop`, nil},
		{`address[0x0a]`, nil},
	}

	for _, c := range cases {
		t.Run(c.p, func(t *testing.T) {
			expected := c.fd
			p := DecodeSymbolicPath(c.p, md)
			actual := fieldDescriptorInMessage(md, p)
			assert.Equal(t, expected, actual, "%s", c.p)
		})
	}

	// DecodeSymbolicPath returns nil for invalid paths, so try it with a populated but invalid Path too
	assert.Nil(t, fieldDescriptorInMessage(md, Path{{Tag: protowire.Number(-1)}}))
}
