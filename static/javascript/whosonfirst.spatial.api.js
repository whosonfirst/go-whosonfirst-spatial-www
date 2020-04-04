var whosonfirst = whosonfirst || {};
whosonfirst.spatial = whosonfirst.spatial || {};

whosonfirst.spatial.api = (function(){

    var self = {

	'point_in_polygon': function(args, on_success, on_error) {

	    var qs = self.query_string(args);

	    var rel_url = "/point-in-polygon?" + qs;
	    return self.get(rel_url, on_success, on_error);
	},

	'point_in_polygon_candidates': function(args, on_success, on_error) {

	    var lat = args['latitude'];
	    var lon = args['longitude'];

	    var rel_url = "/point-in-polygon/candidates?latitude=" + lat + "&longitude=" + lon;
	    return self.get(rel_url, on_success, on_error);
	},
	
	'get': function(rel_url, on_success, on_error) {

	    var abs_url = self.abs_url(rel_url);
	    
	    var req = new XMLHttpRequest();
	    
	    req.onload = function(){
		
		var rsp;
		
		try {
		    rsp = JSON.parse(this.responseText);
            	}
		
		catch (e){
		    console.log("ERR", abs_url, e);
		    on_error(e);
		    return false;
		}
		
		on_success(rsp);
       	    };
	    
	    req.open("get", abs_url, true);
	    req.send();	    
	},

	'abs_url': function(rel_url) {

	    return location.protocol + "//" + location.host + '/api' + rel_url;	// READ ME FROM A DATA ATTRIBUTE...
	},

	'query_string': function(args){

	    var pairs = [];

	    for (var k in args){

		var v = args[k];

		var enc_k = encodeURIComponent(k);
		var enc_v = encodeURIComponent(v);
		
		var pair = enc_k + "=" + enc_v;
		pairs.push(pair);
	    }

	    return pairs.join("&");
	},
    };

    return self;
    
})();
