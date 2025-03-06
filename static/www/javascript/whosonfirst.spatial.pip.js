var whosonfirst = whosonfirst || {};
whosonfirst.spatial = whosonfirst.spatial || {};

whosonfirst.spatial.pip = (function(){

    var styles = {
	"match": {
	    "color": "#000",
	    "weight": 1,
	    "opacity": 1,
	    "fillColor": "#00308F",
	    "fillOpacity": 0.05
	}
    };
    
    var self = {

	'named_style': function(name){
	    return styles[name];
	},
	
	'default_properties': function(){

	    var props_table = {
		"wof:id":"",
		"wof:name":"",
		"wof:placetype":"",
		"edtf:inception":"",
		"edtf:cessation":"",		
	    };

	    return props_table;
	},
	
	'render_properties_table': function(features, props_table){

	    if (! props_table){
		props_table = self.default_properties();
	    }
	    
	    var count = features.length;
	    
	    var table = document.createElement("table");
	    table.setAttribute("class", "table table-striped");	   
	    
	    for (var i=0; i < count; i++){

		var props = features[i];
		
		// draw table header

		if (i % 10 == 0){

		    var tr = document.createElement("tr");
	    
		    for (var k in props_table){
			
			if (self.is_wildcard(k)){
			    
			    for (prop_k in props){
				
				if (! prop_k.startsWith(k)){
				    continue;
				}
				
				var v = prop_k;
				
				var th = document.createElement("th");
				th.appendChild(document.createTextNode(v));
				tr.appendChild(th);				
			    }
			    
			} else {
			    
			    var v = k;	// props_table[k]
			    var th = document.createElement("th");
			    th.appendChild(document.createTextNode(v));
			    tr.appendChild(th);			    
			}		
		    }
		    
		    var thead = document.createElement("thead");
		    thead.setAttribute("class", "thead-dark");
		    thead.appendChild(tr);
		    table.appendChild(thead);		    
		}
		
		var wof_id = props["wof:id"];
		
		var tr = document.createElement("tr");
		tr.setAttribute("id", "tr-" + wof_id);
		
		for (var k in props_table){

		    if (self.is_wildcard(k)){

			for (prop_k in props){

			    if (! prop_k.startsWith(k)){
				continue;
			    }

			    var v = props[prop_k];
			    var node = self.render_value(v);
			    
			    var td = document.createElement("td");
			
			    td.appendChild(node);
			    tr.appendChild(td);
			}
			
		    } else {
			
			var v = props[k];
			var node = self.render_value(v);

			var td = document.createElement("td");
			
			td.appendChild(node);
			tr.appendChild(td);
		    }
		    
		    table.appendChild(tr);
		}
		
	    }

	    var wrapper = document.createElement("div");
	    wrapper.setAttribute("class", "table-responsive");

	    wrapper.appendChild(table);
	    return wrapper;
	},

	'is_wildcard': function(str) {

	    if (str.endsWith(":")){
		return true;
	    }
	    
	    if (str.endsWith("*")){
		return true;
	    }

	    return false;
	},

	'render_value': function(v) {

	    if (typeof(v) == "object"){

		var enc_v = JSON.stringify(v, null, 2);
		var pre = document.createElement("pre");
		pre.appendChild(document.createTextNode(enc_v));

		var summary = document.createElement("summary");
		summary.appendChild(document.createTextNode("details"));
		    
		var details = document.createElement("details");
		details.appendChild(summary);
		details.appendChild(pre);
		
		return details;
	    }

	    else {
		return document.createTextNode(v);
	    }
	},

	init: function(map) {

	    var layers = L.layerGroup();
	    layers.addTo(map);
	    
	    var spinner = new L.Control.Spinner();
	    
	    var update_map = function(e){
		
		var pos = map.getCenter();	

		console.debug("Map center", pos);
		
		var args = {
		    'geometry': {
			'type': 'Point',
			'coordinates': [ pos['lng'], pos['lat'] ],
		    }
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
		
		var existential_filters = document.getElementsByClassName("point-in-polygon-filter-existential");
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
		
		var placetype_filters = document.getElementsByClassName("point-in-polygon-filter-placetype");	
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
		
		var edtf_filters = document.getElementsByClassName("point-in-polygon-filter-edtf");
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
				return whosonfirst.spatial.pip.named_style("match");
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
		    
		    var matches = document.getElementById("pip-matches");
		    matches.innerHTML = "";
		    
		    if (! count){
			return;
		    }
		    
		    for (var i=0; i < count; i++){
			var pl = places[i];
			show_feature(pl["wof:id"]);
		    }
		    
		    var table_props = whosonfirst.spatial.pip.default_properties();
		    
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
		    
		    var table = whosonfirst.spatial.pip.render_properties_table(places, table_props);
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
		
		whosonfirst.spatial.api.point_in_polygon(args).then((rsp) => {
		    on_success(rsp);
		}).catch((err) => {
		    on_error(err);
		});
		
		map.addControl(spinner);	
		layers.clearLayers();	
	    };
	    
	    map.on("moveend", update_map);
	    
	    var filters = document.getElementsByClassName("point-in-polygon-filter");
	    var count_filters = filters.length;
	    
	    for (var i=0; i < count_filters; i++){	    
		var el = filters[i];
		el.onchange = update_map;
	    }
	    
	    var extras = document.getElementsByClassName("point-in-polygon-extra");
	    var count_extras = extras.length;
	    
	    for (var i=0; i < count_extras; i++){	    
		var el = extras[i];
		el.onchange = update_map;
	    }
	    	    
	    slippymap.crosshairs.init(map);

	    whosonfirst.spatial.api.placetypes({}).then((rsp) => {

		var mk_checkbox = function(id, name){

		    var div = document.createElement("div");
		    div.setAttribute("class", "form-check form-check-inline");

		    var input = document.createElement("input");
		    input.setAttribute("type", "checkbox");
		    input.setAttribute("class", "form-check-input point-in-polygon-filter point-in-polygon-filter-placetype");
		    input.setAttribute("id", "placetype-" + id);
		    input.setAttribute("value", name);

		    var label = document.createElement("label");
		    label.setAttribute("class", "form-check-label");
		    label.setAttribute("for", "placetype-" + name);
		    label.appendChild(document.createTextNode(name));

		    div.appendChild(input);
		    div.appendChild(label);

		    return div;
		};

		var placetypes_el = document.getElementById("placetypes");
		
		var count = rsp.length;

		for (var i=0; i < count; i++){
		    var pt = rsp[i];
		    var cb = mk_checkbox(pt.id, pt.name);
		    placetypes_el.appendChild(cb);
		}
		
	    }).catch((err) => {
		console.error("SAD PLACETYPES", err);
	    });
	    
	}
	
    };

    return self;
    
})();
