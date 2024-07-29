gen_cert:
	cd cert && ./gen_cert.sh
clean:
	cd cert && rm -f *.cert *.key *.req 
