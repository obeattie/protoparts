package main

import (
	"fmt"
	"io/ioutil"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"

	ppproto "github.com/obeattie/protoparts/test/proto"
)

func s(str string) *string {
	return &str
}

func msg(name, streetAddress, city *string, tags []string, boop [][]byte, kv map[string]string) *dynamicpb.Message {
	personMd := (&ppproto.Person{}).ProtoReflect().Descriptor()
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

func main() {
	msgs := []*dynamicpb.Message{
		msg(nil, nil, nil, nil, nil, nil),
		msg(s(""), nil, nil, nil, nil, nil),
		msg(s("Oliver Beattie"), nil, nil, nil, nil, nil),
		msg(s("Lindy Bishop"), nil, nil, nil, nil, nil),
		msg(
			s("Ryan Gosling"),
			s("3532 Hayden Ave"),
			s("Culver City"),
			[]string{"The Driver", "Sebastian Wilder", ""},
			[][]byte{[]byte("🕺"), []byte("🏍")},
			map[string]string{
				"k1": "v1",
				"k2": "v2",
				"k3": "v3",
				"k4": "v4",
			}),
		msg(
			s("Ryan Gosling"),
			s("3532 Hayden Ave"),
			s("Culver City"),
			[]string{"The Driver", "Sebastian Wilder"},
			[][]byte{[]byte("🕺"), []byte("🏍")},
			map[string]string{
				"k1": "v1",
				"k2": "v2",
				"k3": "v3",
				"k4": "v4",
			}),
	}
	for i, msg := range msgs {
		b, err := proto.MarshalOptions{Deterministic: true}.Marshal(msg)
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile(fmt.Sprintf("corpus/%d", i), b, 0644)
		if err != nil {
			panic(err)
		}
	}
}
