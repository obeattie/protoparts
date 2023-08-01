.SECONDEXPANSION:
protos: $$(targets)

targets += internal/testproto/proto.pb.go
internal/testproto/proto.pb.go: internal/testproto/proto.proto
	protoc --go_out=. --go_opt=paths=source_relative internal/testproto/proto.proto

targets += internal/testproto/news.pb.go
internal/testproto/news.pb.go: internal/testproto/news.proto
	protoc --go_out=. --go_opt=paths=source_relative internal/testproto/news.proto
