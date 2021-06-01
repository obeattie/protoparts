proto_pb_test.go: proto_test.proto
	protoc --go_out=. --go_opt=paths=source_relative proto_test.proto
	mv proto_test.pb.go proto_pb_test.go
