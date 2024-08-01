gen_cert:
	cd cert && ./gen_cert.sh
test:
	go test ./...go test
port_grpc_protobuf_gen:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./ports/gRPC/server.proto
clean:
	cd cert && rm -f *.cert *.key *.req 
