VERSION = $(shell git rev-parse --short HEAD)

build: assets/ol.js
	go generate
	docker run -v $(CURDIR):/src -e LDFLAGS="-X main.version $(VERSION)" centurylink/golang-builder:latest

container: build ca-certificates.pem
	docker build -t luzifer/locationmaps .

ca-certificates.pem:
		curl -s https://pki.google.com/roots.pem | grep -v "^#" | grep -v "^$$" > $@
		shasum $@

assets/ol.js:
	curl -ssLo assets/ol.js http://www.openlayers.org/api/OpenLayers.js
	sed 's#theme/default/style.css#/assets/style.css#' assets/ol.js > assets/ol.tmp.js
	mv assets/ol.tmp.js assets/ol.js
	curl -ssLo assets/style.css http://www.openlayers.org/api/theme/default/style.css
