var whosonfirst = whosonfirst || {};
whosonfirst.spatial = whosonfirst.spatial || {};

whosonfirst.spatial.intersects = (function(){
    
    var self = {

	init: function(map){

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
	    
	},

	getIntersecting: function(f) {

	    var args = {
		geometry: f.geometry,
	    };
	    
	    whosonfirst.spatial.api.intersects(args).then((rsp) => {
		console.log("INTERSECTS", rsp);
	    }).catch((err) => {
		console.error("Failed to perform intersects query", err);
	    });

	}
	
    };

    return self;
    
})();
