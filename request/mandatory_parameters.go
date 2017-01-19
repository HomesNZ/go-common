package request

import (
	"fmt"
	"net/http"
)

//ValidParameters determines whether the given request contains a value for every parameter in params
func ValidParameters(req *http.Request, params []string) error {
	p := req.URL.Query()
	for _, parameter := range params {
		if p.Get(parameter) == "" {
			return fmt.Errorf("Invalid Parameters")
		}
	}
	return nil
}
