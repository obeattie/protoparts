.SECONDEXPANSION:
protos: $$(targets)

targets += testproto/proto.pb.go
testproto/proto.pb.go: testproto/proto.proto
	protoc --go_out=. --go_opt=paths=source_relative testproto/proto.proto

targets += testproto/news.pb.go
testproto/news.pb.go: testproto/news.proto
	protoc --go_out=. --go_opt=paths=source_relative testproto/news.proto
