package protoparts

import (
	"bytes"
	"encoding/base64"
	"os/exec"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/dynamicpb"

	"github.com/obeattie/protoparts/internal/testproto"
	"github.com/stretchr/testify/assert"
)

// FuzzSplit works off the basic premise that if a message is considered valid by the Protobuf library, then Split()
// should be able to handle it. If that's ever not the case, then the test fails.
func FuzzSplit(f *testing.F) {
	msgs := []*dynamicpb.Message{
		testMsg(f, nil, nil, nil, nil, nil, nil),
		testMsg(f, s(""), nil, nil, nil, nil, nil),
		testMsg(f, s("Oliver Beattie"), nil, nil, nil, nil, nil),
		testMsg(f, s("Lindy Bishop"), nil, nil, nil, nil, nil),
		testMsg(f,
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
		testMsg(f,
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
	for _, msg := range msgs {
		f.Add(marshalProto(f, msg))
	}

	f.Fuzz(func(t *testing.T, b []byte) {
		pb := &testproto.Person{}
		if err := proto.Unmarshal(b, pb); err != nil {
			t.Skip()
		}
		md := pb.ProtoReflect().Descriptor()
		_, err := Split(b, md)
		if !assert.NoError(t, err) {
			t.Logf("base64 protobuf: 		%s", base64.StdEncoding.EncodeToString(b))
			t.Logf("parsed protobuf: 		%s", pb.String())

			protocCmd := exec.Command("protoc", "--decode_raw")
			protocCmd.Stdin = bytes.NewReader(b)
			protocOutput, err := protocCmd.Output()
			if err == nil {
				t.Logf("protoc --decode_raw: 	%s", string(protocOutput))
			}
		}
	})
}
