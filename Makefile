.PHONY: proto
proto:
	@go install github.com/golang/protobuf/protoc-gen-go@latest
	@protoc -I . proto/*.proto --go_out=plugins=grpc:gen --go_opt=paths=source_relative