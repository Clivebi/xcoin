package proxy

import (
	"errors"
	"net/http"
)

type APIHandler struct {
	runner *AppRunner
}

func NewAPIHandler(runner *AppRunner) *APIHandler {
	return &APIHandler{runner: runner}
}

func (o *APIHandler) DispatchRequest(w http.ResponseWriter, r *http.Request) {
	switch r.URL.EscapedPath() {
	case "createbank":
		o.handle_createbank(w, r)
	case "getbank":
		o.handle_getbank(w, r)
	case "changebanklimit":
		o.handle_changebanklimit(w, r)
	case "adduser":
		o.handle_adduser(w, r)
	case "getuser":
		o.handle_getuser(w, r)
	case "issue":
		o.handle_issue(w, r)
	case "chippay":
		o.handle_chippay(w, r)
	case "cashin":
		o.handle_cashin(w, r)
	case "cashout":
		o.handle_cashout(w, r)
	case "transfer":
		o.handle_transfer(w, r)
	default:
		http.Error(w, "invalid URL", 500)
	}
}

func (o *APIHandler) copyArg(r *http.Request, names []string) ([]string, error) {
	ret := make([]string, len(names))
	for i, v := range names {
		val := r.FormValue(v)
		if len(val) == 0 {
			return []string{}, errors.New("argument of:" + v + " is missing")
		}
		ret[i] = val
	}
	return ret, nil
}

func (o *APIHandler) handle_createbank(w http.ResponseWriter, r *http.Request) {
	names := []string{"bankname", "currency", "chip", "exchanger"}
	args, err := o.copyArg(r, names)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	buf := o.runner.SendRequest("addbank", args)
	w.Write(buf)
}

func (o *APIHandler) handle_getbank(w http.ResponseWriter, r *http.Request) {
	names := []string{"bankname"}
	args, err := o.copyArg(r, names)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	buf := o.runner.SendRequest("getbank", args)
	w.Write(buf)
}

func (o *APIHandler) handle_changebanklimit(w http.ResponseWriter, r *http.Request) {
	names := []string{"bankname", "add_threshold"}
	args, err := o.copyArg(r, names)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	buf := o.runner.SendRequest("addbanklimit", args)
	w.Write(buf)
}

func (o *APIHandler) handle_adduser(w http.ResponseWriter, r *http.Request) {

}

func (o *APIHandler) handle_getuser(w http.ResponseWriter, r *http.Request) {

}

func (o *APIHandler) handle_issue(w http.ResponseWriter, r *http.Request) {

}

func (o *APIHandler) handle_chippay(w http.ResponseWriter, r *http.Request) {

}

func (o *APIHandler) handle_cashin(w http.ResponseWriter, r *http.Request) {

}

func (o *APIHandler) handle_cashout(w http.ResponseWriter, r *http.Request) {

}

func (o *APIHandler) handle_transfer(w http.ResponseWriter, r *http.Request) {

}
