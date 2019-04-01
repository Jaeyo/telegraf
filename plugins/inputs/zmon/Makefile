

build-proto:
	@rm -rf protodata/* && \
	protoc -I ./proto/ --go_out=plugins=grpc:protodata ./proto/*.proto
