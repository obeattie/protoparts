test/proto/proto.pb.go: test/proto/proto.proto
	protoc --go_out=. --go_opt=paths=source_relative test/proto/proto.proto
