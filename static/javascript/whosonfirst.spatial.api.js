var whosonfirst = whosonfirst || {};
whosonfirst.spatial = whosonfirst.spatial || {};

whosonfirst.spatial.api = (function(){

    var self = {

	'point_in_polygon': function(args, on_success, on_error) {

	    var rel_url = "/point-in-polygon";
	    return self.post(rel_url, args, on_success, on_error);
	},

	'point_in_polygon_candidates': function(args, on_success, on_error) {

	    return self.post(rel_url, args, on_success, on_error);
	},

	'post': function(rel_url, args, on_success, on_error) {

	    var abs_url = self.abs_url(rel_url);
	    
	    var req = new XMLHttpRequest();

	    /*
	    var form_data = args;

	    if (! form_data.append){
		
		form_data = new FormData();
		
		for (key in args){
		    form_data.append(key, args[key]);
		}
	    }

	    if (method.verb() == "GET"){

		if (form_data.keys()){}
	    }

	    */
					    
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
	    
	    req.open("POST", abs_url, true);

	    req.setRequestHeader("Content-type", "application/json");
	    req.setRequestHeader("Accept", "application/geo+json");
	    
	    var enc_args = JSON.stringify(args);
	    req.send(enc_args);	    
	},
	
	'get': function(rel_url, args, on_success, on_error) {

	    var qs = self.query_string(args);
	    rel_url = rel_url + "?" + qs;

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
