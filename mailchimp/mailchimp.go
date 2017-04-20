package mailchimp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const baseURI = "https://%s.api.mailchimp.com/3.0"

// MailChimp represents a configuration state for the API client.
type MailChimp struct {
	config     *Config
	httpClient *http.Client
}

// NewClient returns a pointer to a new instance of MailChimp.
func NewClient(config *Config) (*MailChimp, error) {
	err := config.Validate()
	if err != nil {
		return nil, err
	}
	return &MailChimp{
		config:     config,
		httpClient: &http.Client{},
	}, nil
}

func (m MailChimp) baseURI() string {
	return fmt.Sprintf(baseURI, m.config.DataCenter)
}

func (m MailChimp) Get(url string, params interface{}, responseData interface{}) error {
	return m.doRequest("GET", url, params, nil, responseData)
}

func (m MailChimp) Post(url string, params interface{}, requestData interface{}, responseData interface{}) error {
	return m.doRequest("POST", url, params, requestData, responseData)
}

func (m MailChimp) Patch(url string, params interface{}, requestData interface{}, responseData interface{}) error {
	return m.doRequest("PATCH", url, params, requestData, responseData)
}

func (m MailChimp) Put(url string, params interface{}, requestData interface{}, responseData interface{}) error {
	return m.doRequest("PUT", url, params, requestData, responseData)
}

func (m MailChimp) Delete(url string, params interface{}, requestData interface{}) error {
	return m.doRequest("DELETE", url, params, requestData, nil)
}

func (m MailChimp) doRequest(method string, url string, params interface{}, requestData interface{}, responseData interface{}) error {
	urlStr, err := formatUrl(url, params)
	if err != nil {
		return err
	}
	body := bytes.Buffer{}
	if requestData != nil {
		var b []byte
		// format request data as json
		b, err = json.Marshal(requestData)
		if err != nil {
			return err
		}
		body.Write(b)
	}
	req, err := http.NewRequest(method, m.baseURI()+urlStr, &body)
	if err != nil {
		return err
	}
	res, err := m.httpClient.Do(req)
	if err != nil {
		return err
	}
	if responseData == nil {
		// Don't need to parse the response if there is no object to parse it into.
		return nil
	}
	defer res.Body.Close()
	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(responseBody, responseData)
}
