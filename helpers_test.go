package protoparts

import (
	"fmt"
	"testing"

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
	personMd := (&Person{}).ProtoReflect().Descriptor()
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
	personMd := (&Person{}).ProtoReflect().Descriptor()
	addressMd := personMd.Fields().ByName("address").Message()
	latLngMd := addressMd.Fields().ByName("lat_lng").Message()

	person := dynamicpb.NewMessage(personMd)
	person.Set(personMd.Fields().ByName("name"), protoreflect.ValueOf("Ryan Gosling"))
	address := dynamicpb.NewMessage(addressMd)
	address.Set(addressMd.Fields().ByName("street_address"), protoreflect.ValueOf("3532 Hayden Ave"))
	address.Set(addressMd.Fields().ByName("city"), protoreflect.ValueOf("Culver City"))
	latLng := dynamicpb.NewMessage(latLngMd)
	latLng.Set(latLngMd.Fields().ByName("latitude"), protoreflect.ValueOfFloat64(34.0257178))
	latLng.Set(latLngMd.Fields().ByName("longitude"), protoreflect.ValueOfFloat64(-118.3802275))
	address.Set(addressMd.Fields().ByName("lat_lng"), protoreflect.ValueOf(latLng))
	person.Set(personMd.Fields().ByName("address"), protoreflect.ValueOf(address))
	moarAddresses := person.NewField(personMd.Fields().ByName("moar_addresses")).List()
	moarAddresses.Append(protoreflect.ValueOf(address))
	moarAddresses.Append(protoreflect.ValueOf(address))
	person.Set(personMd.Fields().ByName("moar_addresses"), protoreflect.ValueOf(moarAddresses))
	tags := person.NewField(personMd.Fields().ByName("tags")).List()
	tags.Append(protoreflect.ValueOf("The Driver"))
	tags.Append(protoreflect.ValueOf("Sebastian Wilder"))
	tags.Append(protoreflect.ValueOf("")) // deliberately diabolical
	person.Set(personMd.Fields().ByName("tags"), protoreflect.ValueOf(tags))
	boop := person.NewField(personMd.Fields().ByName("boop")).List()
	boop.Append(protoreflect.ValueOf([]byte("🕺")))
	boop.Append(protoreflect.ValueOf([]byte("🏍")))
	person.Set(personMd.Fields().ByName("boop"), protoreflect.ValueOf(boop))
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

	mapStringString := person.NewField(personMd.Fields().ByName("map_string_string")).Map()
	mapStringString.Set(protoreflect.ValueOf("k1").MapKey(), protoreflect.ValueOf("v1"))
	mapStringString.Set(protoreflect.ValueOf("k2").MapKey(), protoreflect.ValueOf("v2"))
	mapStringString.Set(protoreflect.ValueOf("k3").MapKey(), protoreflect.ValueOf("v3"))
	mapStringString.Set(protoreflect.ValueOf("").MapKey(), protoreflect.ValueOf("")) // deliberately diabolical
	person.Set(personMd.Fields().ByName("map_string_string"), protoreflect.ValueOf(mapStringString))

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

// split is a shortcut which splits the dynamic message into parts
func split(t testing.TB, msg protoreflect.Message) Parts {
	pb, err := proto.Marshal(msg.Interface())
	require.NoError(t, err)
	parts, err := Split(pb, msg.Descriptor())
	require.NoError(t, err)
	return parts
}
