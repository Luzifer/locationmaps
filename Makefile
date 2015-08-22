VERSION = $(shell git describe --tags)

bindata: assets/ol.js
	go generate

build:
	docker run -v $(CURDIR):/src -e LDFLAGS="-X main.version $(VERSION)" centurylink/golang-builder:latest

container: build
	docker build -t luzifer/locationmaps .

assets/ol.js:
	curl -ssLo assets/ol.js http://www.openlayers.org/api/OpenLayers.js
	sed 's#theme/default/style.css#/assets/style.css#' assets/ol.js > assets/ol.tmp.js
	mv assets/ol.tmp.js assets/ol.js
	curl -ssLo assets/style.css http://www.openlayers.org/api/theme/default/style.css
