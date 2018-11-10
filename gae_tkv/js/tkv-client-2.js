var france = L.map('france').setView([47, 0], 5);

L.tileLayer('https://api.tiles.mapbox.com/v4/{id}/{z}/{x}/{y}.png?access_token=pk.eyJ1IjoibWFwYm94IiwiYSI6ImNpejY4NXVycTA2emYycXBndHRqcmZ3N3gifQ.rJcFIG214AriISLbB6B5aw', {
	maxZoom: 18,
	attribution: 'Map data &copy; <a href="https://www.openstreetmap.org/">OpenStreetMap</a> contributors, ' +
		'<a href="https://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>, ' +
		'Imagery © <a href="https://www.mapbox.com/">Mapbox</a>',
	id: 'mapbox.streets'
}).addTo(france);

var hostname
var protocol
var port
var targetService

var oReq 

var littleIcon = L.icon({
iconUrl: '9pixels.png',

	iconSize:     [3, 3], // size of the icon
	iconAnchor:   [22, 94], // point of the icon which will correspond to marker's location
});



function onFranceMapClick(e) {

	hostname = window.location.hostname
	protocol = window.location.protocol
	port = window.location.port
	targetService = protocol + "//"+ hostname + ":" + port + "/"


	var jsonLatLng = JSON.stringify( e.latlng);
	console.log( jsonLatLng);

	oReq = new XMLHttpRequest();
	// oReq.responseType = 'json';
	oReq.addEventListener("load", reqListener);
	oReq.open("POST", targetService +'translateLatLngInSourceCountryToLatLngInTargetCountry');
	oReq.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
	oReq.send( jsonLatLng);				
};

function reqListener( evt) {
	
	var jsonResponse = JSON.parse( this.response)
	
	console.log('village translateLatLngInSourceCountryToLatLngInTargetCountry answer', 
		jsonResponse.X, jsonResponse.Y);

	lat = parseFloat(jsonResponse.LatClosest);
	lng = parseFloat(jsonResponse.LngClosest);

	latTarget = parseFloat(jsonResponse.LatTarget);
	lngTarget = parseFloat(jsonResponse.LngTarget);

	xSpead = parseFloat(jsonResponse.Xspread);
	ySpead = parseFloat(jsonResponse.Yspread);

	message = "Territory X="+ 
		Math.floor(100*jsonResponse.Xspread)+" Y="+
		Math.floor(100*jsonResponse.Yspread);
			
	L.marker([lat, lng]).addTo(france)
		.bindPopup( message).openPopup();
		
	L.marker([latTarget, lngTarget]).addTo(haiti)
		.bindPopup( message).openPopup();


	for (var i = 0; i < jsonResponse.SourceBorderPoints[0].length; i++) {

		lng = parseFloat(jsonResponse.SourceBorderPoints[0][i][0]);
		lat = parseFloat(jsonResponse.SourceBorderPoints[0][i][1]);

		marker = new L.marker([lat,lng], {icon: littleIcon})
			.addTo(france);
	}
};


france.on('click', onFranceMapClick);

var haiti = L.map('haiti').setView([18, -72], 5);

L.tileLayer('https://api.tiles.mapbox.com/v4/{id}/{z}/{x}/{y}.png?access_token=pk.eyJ1IjoibWFwYm94IiwiYSI6ImNpejY4NXVycTA2emYycXBndHRqcmZ3N3gifQ.rJcFIG214AriISLbB6B5aw', {
	maxZoom: 18,
	attribution: 'Map data &copy; <a href="https://www.openstreetmap.org/">OpenStreetMap</a> contributors, ' +
		'<a href="https://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>, ' +
		'Imagery © <a href="https://www.mapbox.com/">Mapbox</a>',
	id: 'mapbox.streets'
}).addTo(haiti);
