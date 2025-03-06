var whosonfirst = whosonfirst || {};
whosonfirst.spatial = whosonfirst.spatial || {};

whosonfirst.spatial.intersects = (function(){

    var layers;
    var spinner;
    
    var self = {

	init: function(map){
	    
	    layers = L.layerGroup();
	    layers.addTo(map);
	    
	    spinner = new L.Control.Spinner();
	    
	    map.pm.addControls({  
		position: 'topright',
		drawMarker: false,
		drawCircle: false,
		drawCircleMarker: false,
		drawPolyline: false,
		drawText: false,
		editMode: false,
		rotateMode: false,
		cutPolygon: false,
		dragMode: false,
	    });

	    map.on("pm:drawstart", (e) => {
		// console.log("draw start", e);
	    });

	    map.on("pm:drawend", (shp) => {
		console.log("draw start", shp);

		var feature_group = map.pm.getGeomanLayers(true);
		var feature_collection = feature_group.toGeoJSON();

		var features = feature_collection.features;
		var count = features.length;

		for (var i=0; i < count; i++){
		    self.getIntersecting(features[i]);
		}
	    });

	    whosonfirst.spatial.placetypes.init().catch((err) => {
		console.error("Failed to initialize placetypes", err);
	    });
	    
	},

	getIntersecting: function(f) {

	    var args = {
		geometry: f.geometry,
	    };

	    var properties = [];
	    
	    var extra_properties = document.getElementById("extras");
	    
	    if (extra_properties){
		
		var extras = extra_properties.value;
		
		if (extras){
		    properties = extras.split(",");
		    args['properties'] = properties;
		}
	    }
	    
	    var existential_filters = document.getElementsByClassName("spatial-filter-existential");
	    var count_existential = existential_filters.length;
	    
	    for (var i=0; i < count_existential; i++){
		
		var el = existential_filters[i];
		
		if (! el.checked){
		    continue;
		}
		
		var fl = el.value;
		args[fl] = [ 1 ];
	    }
	    
	    var placetypes = [];
	    
	    var placetype_filters = document.getElementsByClassName("spatial-filter-placetype");	
	    var count_placetypes = placetype_filters.length;
	    
	    for (var i=0; i < count_placetypes; i++){
		
		var el = placetype_filters[i];
		
		if (! el.checked){
		    continue;
		}
		
		var pt = el.value;
		placetypes.push(pt);
	    }
	    
	    if (placetypes.length > 0){
		args['placetypes'] = placetypes;
	    }
	    
	    var edtf_filters = document.getElementsByClassName("spatial-filter-edtf");
	    var count_edtf = edtf_filters.length;
		
	    for (var i=0; i < count_edtf; i++){
		
		var el = edtf_filters[i];
		
		var id = el.getAttribute("id");
		
		if (! id.match("^(inception|cessation)$")){
		    continue
		}
		
		var value = el.value;
		
		if (value == ""){
		    continue;
		}
		
		// TO DO: VALIDATE EDTF HERE WITH WASM
		// https://millsfield.sfomuseum.org/blog/2021/01/14/edtf/
		
		var key = id + "_date";
		args[key] = value;
	    };
	    
	    var show_feature = function(id){
		
		var url = "/data/" + id;
		
		var on_success = function(data){
		    
		    var l = L.geoJSON(data, {
			style: function(feature){
			    return whosonfirst.spatial.results.named_style("match");
			},
		    });
		    
		    layers.addLayer(l);
		    l.bringToFront();
		};
		
		var on_fail= function(err){
		    console.log("SAD", id, err);
		}
		
		whosonfirst.net.fetch(url, on_success, on_fail);
	    };
	    
	    var on_success = function(rsp){
		
		map.removeControl(spinner);
		
		var places = rsp["places"];
		var count = places.length;
		
		var matches = document.getElementById("intersects-matches");
		matches.innerHTML = "";
		
		if (! count){
		    return;
		}
		
		for (var i=0; i < count; i++){
		    var pl = places[i];
			show_feature(pl["wof:id"]);
		}
		
		var table_props = whosonfirst.spatial.results.default_properties();
		
		// START OF something something something
		
		var extras_el = document.getElementById("extras");
		
		if (extras_el){
		    
		    var str_extras = extras_el.value;
		    var extras = null;
		    
		    if (str_extras){
			extras = str_extras.split(",");  		    
		    }
		    
		    if (extras){
			
			var first = places[0];
			
			var count_extras = extras.length;		    
			var extra_props = [];
			
			for (var i=0; i < count_extras; i++){
			    
			    var ex = extras[i];
			    
			    if ((ex.endsWith(":")) || (ex.endsWith(":*"))){
				
				var prefix = ex.replace("*", "");
				
				for (k in first){
				    if (k.startsWith(prefix)){
					extra_props.push(k);
				    }
				}
				
			    } else {
				
				if (first[ex]) {
				    extra_props.push(ex);
				}
			    }
			}
			
			for (idx in extra_props){
			    var ex = extra_props[idx];
			    table_props[ex] = "";
			}
		    }
		    
		}
		
		// END OF something something something
		
		var table = whosonfirst.spatial.results.render_properties_table(places, table_props);
		matches.appendChild(table);
		
	    };
	    
	    var on_error = function(err){
		
		var matches = document.getElementById("pip-matches");
		matches.innerHTML = "";
		
		map.removeControl(spinner);	    
		console.error("Point in polygon request failed", err);
	    }
	    
	    args["sort"] = [
		"placetype://",
		"name://",
		"inception://",
	    ];
	    
	    whosonfirst.spatial.api.intersects(args).then((rsp) => {
		console.log("INTERSECTS", rsp);
		on_success(rsp);
	    }).catch((err) => {
		on_error(err);
	    });

	    map.addControl(spinner);
	    layers.clearLayers();
	}
	
    };

    return self;
    
})();
