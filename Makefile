VERSION = $(shell git rev-parse --short HEAD)

build:
	go generate
	docker run -v $(CURDIR):/src -e LDFLAGS="-X main.version $(VERSION)" centurylink/golang-builder:latest

container: build ca-certificates.pem
	docker build -t luzifer/locationmaps .

ca-certificates.pem:
		curl -s https://pki.google.com/roots.pem | grep -v "^#" | grep -v "^$$" > $@
		shasum $@
