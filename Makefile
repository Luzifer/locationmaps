VERSION = $(shell git describe --tags)

bindata: assets/ol.js
	go generate

build:
	docker run -v $(CURDIR):/src -e LDFLAGS="-X main.version $(VERSION)" centurylink/golang-builder:latest

container: build ca-certificates.pem
	docker build -t luzifer/locationmaps .

ca-certificates.pem:
	curl -ssLo certdata.txt https://hg.mozilla.org/mozilla-central/raw-file/tip/security/nss/lib/ckfw/builtins/certdata.txt
	./extract-nss-root-certs > ca-certificates.pem
	rm certdata.txt

assets/ol.js:
	curl -ssLo assets/ol.js http://www.openlayers.org/api/OpenLayers.js
	sed 's#theme/default/style.css#/assets/style.css#' assets/ol.js > assets/ol.tmp.js
	mv assets/ol.tmp.js assets/ol.js
	curl -ssLo assets/style.css http://www.openlayers.org/api/theme/default/style.css
