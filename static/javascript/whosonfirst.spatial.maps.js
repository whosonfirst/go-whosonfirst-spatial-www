var whosonfirst = whosonfirst || {};
whosonfirst.spatial = whosonfirst.spatial || {};

whosonfirst.spatial.maps = (function(){

    var attribution = '<a href="https://github.com/tangrams" target="_blank">Tangram</a> | <a href="http://www.openstreetmap.org/copyright" target="_blank">&copy; OpenStreetMap contributors</a> | <a href="https://www.nextzen.org/" target="_blank">Nextzen</a>';
   
    var maps = {};

    var self = {

	'parseHash': function(hash_str){

	    if (hash_str.indexOf('#') === 0) {
		hash_str = hash_str.substr(1);
	    }

	    var lat;
	    var lon;
	    var zoom;
	    
	    var args = hash_str.split("/");
	    
	    if (args.length != 3){
		console.log("Unrecognized hash string");
		return null;
	    }
	    
	    zoom = args[0];
	    lat = args[1];
	    lon = args[2];			
	    
	    zoom = parseInt(zoom, 10);
	    lat = parseFloat(lat);
	    lon = parseFloat(lon);		
	    
	    if (isNaN(zoom) || isNaN(lat) || isNaN(lon)) {
		console.log("Invalid zoom/lat/lon", zoom, lat, lon);
		return null;
	    }

	    var parsed = {
		'latitude': lat,
		'longitude': lon,
		'zoom': zoom,
	    };

	    return parsed;
	},
	
	'getMap': function(map_el, args){

	    if (! args){
		args = {};
	    }
	    
	    var map_id = map_el.getAttribute("id");

	    if (! map_id){
		console.log("SAD");
		return;
	    }
	    
	    if (maps[map_id]){
		return maps[map_id];
	    }

	    var map = L.map("map");
	    
	    var tile_url = map_el.getAttribute("data-leaflet-tile-url");
	    tile_url = decodeURIComponent(tile_url);
	    
	    if (tile_url != "") {
		
		var layer = L.tileLayer(tile_url, {});
		layer.addTo(map);
		
	    } else {
		
		var api_key = args["api_key"];
		
		var tangram_opts = self.getTangramOptions(args);	   
		var tangramLayer = Tangram.leafletLayer(tangram_opts);
		
		tangramLayer.addTo(map);

		var attribution = self.getAttribution();
		map.attributionControl.addAttribution(attribution);	    		
	    }
	    
	    return map;
	},

	'getTangramOptions': function(args){

	    if (! args){
		args = {};
	    }

	    if (! args["api_key"]){
		return null;
	    }

	    /*
	    var sceneText = await fetch(new Request('https://somwehere.com/scene.zip', { headers: { 'Accept': 'application/zip' } })).then(r => r.text());
	    var sceneURL = URL.createObjectURL(new Blob([sceneText]));
	    scene.load(sceneURL, { base_path: 'https://somwehere.com/' });
	    */
	    
	    var api_key = args["api_key"];
	    var style_url = args["style_url"];
	    var tile_url = args["tile_url"];	    
	    
	    var tangram_opts = {
		scene: {
		    import: [
			style_url,
		    ],
		    sources: {
			mapzen: {
			    url: tile_url,
			    url_subdomains: ['a', 'b', 'c', 'd'],
			    url_params: {api_key: api_key},
			    tile_size: 512,
			    max_zoom: 18
			}
		    }
		}
	    };

	    return tangram_opts;
	},

	'getAttribution': function(){
	    return attribution;
	},
    };

    return self;
    
})();
