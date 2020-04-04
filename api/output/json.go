package output

import (
	"encoding/json"
	"net/http"
)

func AsJSON(rsp http.ResponseWriter, data interface{}) {

	js, err := json.Marshal(data)

	if err != nil {
		http.Error(rsp, err.Error(), http.StatusInternalServerError)
		return
	}

	rsp.Header().Set("Content-Type", "application/json")
	rsp.Write(js)
}
