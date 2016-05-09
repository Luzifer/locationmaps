assets/ol.js:
	curl -ssLo assets/ol.js http://www.openlayers.org/api/OpenLayers.js
	sed 's#theme/default/style.css#/assets/style.css#' assets/ol.js > assets/ol.tmp.js
	mv assets/ol.tmp.js assets/ol.js
	curl -ssLo assets/style.css http://www.openlayers.org/api/theme/default/style.css
