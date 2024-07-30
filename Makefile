gen_cert:
	cd cert && ./gen_cert.sh
test:
	go test ./...go test
clean:
	cd cert && rm -f *.cert *.key *.req 
