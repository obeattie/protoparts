.SECONDEXPANSION:
protos: $$(targets)

targets += test/proto/proto.pb.go
test/proto/proto.pb.go: test/proto/proto.proto
	protoc --go_out=. --go_opt=paths=source_relative test/proto/proto.proto

targets += test/proto/news.pb.go
test/proto/news.pb.go: test/proto/news.proto
	protoc --go_out=. --go_opt=paths=source_relative test/proto/news.proto
