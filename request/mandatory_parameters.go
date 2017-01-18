package request

import "net/http"

//ErrorMissingRequiredParameters is returned when a request does not contain required parameters
var ErrorMissingRequiredParameters error

//ValidParameters determines whether the given request contains a value for every parameter in params
func ValidParameters(req *http.Request, params []string) bool {
	p := req.URL.Query()

	for _, parameter := range params {
		if p.Get(parameter) == "" {
			return false
		}
	}
	return true
}
