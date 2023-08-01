package protoparts

import (
	"fmt"
	"testing"

	"github.com/obeattie/protoparts/testproto"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

func s(str string) *string {
	return &str
}

func marshalProto(t testing.TB, msg proto.Message) []byte {
	b, err := proto.MarshalOptions{Deterministic: true}.Marshal(msg)
	require.NoError(t, err)
	return b
}

func testMsg(t testing.TB, name, streetAddress, city *string, tags []string, boop [][]byte, kv map[string]string) *dynamicpb.Message {
	personMd := (&testproto.Person{}).ProtoReflect().Descriptor()
	person := dynamicpb.NewMessage(personMd)
	if name != nil {
		person.Set(personMd.Fields().ByName("name"), protoreflect.ValueOfString(*name))
	}
	if streetAddress != nil || city != nil {
		addressMd := personMd.Fields().ByName("address").Message()
		address := dynamicpb.NewMessage(addressMd)
		if streetAddress != nil {
			address.Set(addressMd.Fields().ByName("street_address"), protoreflect.ValueOfString(*streetAddress))
		}
		if city != nil {
			address.Set(addressMd.Fields().ByName("city"), protoreflect.ValueOfString(*city))
		}
		person.Set(personMd.Fields().ByName("address"), protoreflect.ValueOfMessage(address))
	}
	if tags != nil {
		tagsL := person.NewField(personMd.Fields().ByName("tags")).List()
		for _, t := range tags {
			tagsL.Append(protoreflect.ValueOfString(t))
		}
		person.Set(personMd.Fields().ByName("tags"), protoreflect.ValueOfList(tagsL))
	}
	if boop != nil {
		boopL := person.NewField(personMd.Fields().ByName("boop")).List()
		for _, b := range boop {
			boopL.Append(protoreflect.ValueOfBytes(b))
		}
		person.Set(personMd.Fields().ByName("boop"), protoreflect.ValueOfList(boopL))
	}
	if kv != nil {
		kvM := person.NewField(personMd.Fields().ByName("map_string_string")).Map()
		for k, v := range kv {
			kvM.Set(protoreflect.ValueOfString(k).MapKey(), protoreflect.ValueOfString(v))
		}
		person.Set(personMd.Fields().ByName("map_string_string"), protoreflect.ValueOfMap(kvM))
	}
	return person
}

func quickTestMsg(t testing.TB) *dynamicpb.Message {
	person := testMsg(t,
		s("Ryan Gosling"),
		s("3532 Hayden Ave"),
		s("Culver City"),
		[]string{
			"The Driver",
			"Sebastian Wilder",
			"", // deliberately diabolical
		},
		[][]byte{
			[]byte("🕺"),
			[]byte("🏍"),
		},
		map[string]string{
			"k1": "v1",
			"k2": "v2",
			"k3": "v3",
			"k4": "v4",
			"":   "", // deliberately diabolical
		})
	personMd := person.Type().Descriptor()

	address := person.Get(personMd.Fields().ByName("address")).Message()
	addressMd := address.Type().Descriptor()
	latLngMd := addressMd.Fields().ByName("lat_lng").Message()
	latLng := dynamicpb.NewMessage(latLngMd)
	latLng.Set(latLngMd.Fields().ByName("latitude"), protoreflect.ValueOfFloat64(34.0257178))
	latLng.Set(latLngMd.Fields().ByName("longitude"), protoreflect.ValueOfFloat64(-118.3802275))
	address.Set(addressMd.Fields().ByName("lat_lng"), protoreflect.ValueOf(latLng))
	person.Set(personMd.Fields().ByName("address"), protoreflect.ValueOf(address))

	moarAddresses := person.NewField(personMd.Fields().ByName("moar_addresses")).List()
	moarAddresses.Append(protoreflect.ValueOf(address))
	moarAddresses.Append(protoreflect.ValueOf(address))
	person.Set(personMd.Fields().ByName("moar_addresses"), protoreflect.ValueOf(moarAddresses))

	single := personMd.Fields().ByName("marital_status").Enum().Values().ByName("SINGLE")
	person.Set(personMd.Fields().ByName("marital_status"), protoreflect.ValueOf(single.Number()))
	person.Set(personMd.Fields().ByName("maybe_latlng"), protoreflect.ValueOf(latLng))

	mapStringLatLng := person.NewField(personMd.Fields().ByName("map_string_latlng")).Map()
	for i := 0; i < 100; i++ { // 100 map elements
		mapStringLatLng.Set(protoreflect.ValueOf(fmt.Sprintf("%da", i)).MapKey(), protoreflect.ValueOf(latLng))
		mapStringLatLng.Set(protoreflect.ValueOf(fmt.Sprintf("%db", i)).MapKey(), protoreflect.ValueOf(latLng))
		mapStringLatLng.Set(protoreflect.ValueOf(fmt.Sprintf("%dc", i)).MapKey(), protoreflect.ValueOf(latLng))
		mapStringLatLng.Set(protoreflect.ValueOf(fmt.Sprintf("%dd", i)).MapKey(), protoreflect.ValueOf(latLng))
	}
	// Also set an empty-valued key – this is an important edge case, serialising as a map element with a 0-length value
	emptyLatLng := dynamicpb.NewMessage(latLngMd)
	mapStringLatLng.Set(protoreflect.ValueOf("empty").MapKey(), protoreflect.ValueOf(emptyLatLng))
	person.Set(personMd.Fields().ByName("map_string_latlng"), protoreflect.ValueOf(mapStringLatLng))

	return person
}

// benchmarkMsg returns a message for use in all benchmarks (so benchmarks involving serialisation/deserialisation code
// are using a common baseline)
func benchmarkMsg(t testing.TB) *dynamicpb.Message {
	return testMsg(t,
		s("Sirius Black"),
		s("12 Grimmauld Place"),
		s("London"),
		[]string{"confringo"},
		nil,
		map[string]string{
			"wizard": "yes"})
}

// split is a shortcut which splits the dynamic message into Parts
func split(t testing.TB, msg protoreflect.Message) Parts {
	pb := marshalProto(t, msg.Interface())
	parts, err := Split(pb, msg.Descriptor())
	require.NoError(t, err)
	return parts
}
