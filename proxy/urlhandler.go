package proxy

import (
	//"errors"
	"net/http"
)

// APIHandler handle the http request
type APIHandler struct {
	runner *AppRunner
}

// NewAPIHandler get the instance of APIHandler
func NewAPIHandler(runner *AppRunner) *APIHandler {
	return &APIHandler{runner: runner}
}

//DispatchRequest process the http request
func (o *APIHandler) DispatchRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.EscapedPath() != "/callapi.do" {
		http.Error(w, "URL not handler:"+r.URL.EscapedPath(), 500)
		return
	}
	o.handleCallAPI(w, r)
}

func (o *APIHandler) copyArg(r *http.Request, names []string) ([]string, error) {
	ret := make([]string, len(names))
	for i, v := range names {
		val := r.FormValue(v)
		//if len(val) == 0 {
		//	return []string{}, errors.New("argument of:" + v + " is missing")
		//}
		ret[i] = val
	}
	return ret, nil
}

func (o *APIHandler) handleCallAPI(w http.ResponseWriter, r *http.Request) {
	names := []string{"req", "sig"}
	args, err := o.copyArg(r, names)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	buf := o.runner.SendRequest("callapi", args)
	w.Write(buf)
}
