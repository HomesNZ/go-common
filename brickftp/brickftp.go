package brickftp

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"
)

// Client for the BrickFTP REST API
type Client struct {
	Client  *http.Client
	BaseURL url.URL
	APIKey  string
}

type FileInfo struct {
	Path         string    `json:"path"`
	DisplayName  string    `json:"display_name"`
	Type         string    `json:"type"`
	Size         int64     `json:"size"`
	ModifiedTime time.Time `json:"mtime"`
	CRC32        string    `json:"crc32"`
	MD5          string    `json:"md5"`
	DownloadURI  string    `json:"download_uri"`
}

// Download requests a file to be downloaded
func (c Client) Download(fileName string) (*FileInfo, error) {
	req, err := c.newRequest(http.MethodGet, c.filesURL(fileName), nil)
	if err != nil {
		return nil, err
	}
	resp := &FileInfo{}
	res, err := c.doRequest(req, resp)
	if err != nil {
		return nil, err
	}
	const successStatus = http.StatusOK
	if res.StatusCode != successStatus {
		return nil, c.err("Download", "expected status %d, got: %d filename: %s", successStatus, res.StatusCode, fileName)
	}
	return resp, nil
}

type filesMoveRequest struct {
	MoveDestination string `json:"move-destination"`
}

// Move a file
func (c Client) Move(oldName string, newName string) error {
	r := filesMoveRequest{newName}
	req, err := c.newRequest(http.MethodPost, c.filesURL(oldName), r)
	if err != nil {
		return err
	}
	res, err := c.doRequest(req, nil)
	if err != nil {
		return err
	}
	if (res.StatusCode < 200) && (res.StatusCode > 300) {
		return c.err("Move", "expected status, got: %d filename: %s", res.StatusCode, oldName)
	}
	return nil
}

type actionRequest struct {
	Action string `json:"action"`
	Ref    string `json:"ref,omitempty"`
	Part   int    `json:"part,omitempty"`
}

type Upload struct {
	Ref                string `json:"ref"`
	Path               string `json:"path"`
	Action             string `json:"action"`
	AskAboutOverwrites bool   `json:"ask_about_overwrites"`
	HTTPMethod         string `json:"http_method"`
	UploadURI          string `json:"upload_uri"`
	PartSize           int64  `json:"partsize"`
	PartNumber         int    `json:"part_number"`
	AvailableParts     int    `json:"available_parts"`
}

// StartUpload requests an upload to be started
func (c Client) StartUpload(name string) (*Upload, error) {
	r := actionRequest{Action: "put"}
	req, err := c.newRequest(http.MethodPost, c.filesURL(name), r)
	if err != nil {
		return nil, err
	}
	resp := &Upload{}
	res, err := c.doRequest(req, resp)
	if err != nil {
		return resp, err
	}
	const successStatus = http.StatusOK
	if res.StatusCode != successStatus {
		return resp, c.err("StartUpload", "expected status %d, got: %d filename: %s", successStatus, res.StatusCode, name)
	}
	return resp, nil
}

// UploadPart uploads a file to the UploadURI given by Upload.
// Max file size is 5GB per part
func (c Client) UploadPart(u *Upload, r io.Reader, contentLength int64) error {
	req, err := http.NewRequest(http.MethodPut, u.UploadURI, r)
	if err != nil {
		return err
	}
	req.ContentLength = contentLength
	res, err := c.doRequest(req, nil)
	if err != nil {
		return err
	}
	const successStatus = http.StatusOK
	if res.StatusCode != successStatus {
		return c.err("UploadPart", "expected status %d, got: %d filename: %s", successStatus, res.StatusCode, u.Path)
	}
	return nil
}

// RequestUploadPart requests a new upload URI for additional parts
func (c Client) RequestUploadPart(u *Upload, part int) (*Upload, error) {
	r := actionRequest{Action: "put", Ref: u.Ref, Part: part}
	req, err := c.newRequest(http.MethodPost, c.filesURL(u.Path), r)
	if err != nil {
		return nil, err
	}
	resp := &Upload{}
	res, err := c.doRequest(req, resp)
	if err != nil {
		return resp, err
	}
	const successStatus = http.StatusOK
	if res.StatusCode != successStatus {
		return resp, c.err("RequestUploadPart", "expected status %d, got: %d filename: %s", successStatus, res.StatusCode, u.Path)
	}
	return resp, nil
}

// CompleteUpload notifies that all parts have been uploaded
func (c Client) CompleteUpload(u *Upload) (*FileInfo, error) {
	r := actionRequest{Action: "end", Ref: u.Ref}
	req, err := c.newRequest(http.MethodPost, c.filesURL(u.Path), r)
	if err != nil {
		return nil, err
	}
	resp := &FileInfo{}
	res, err := c.doRequest(req, resp)
	if err != nil {
		return resp, err
	}
	const successStatus = http.StatusCreated
	if res.StatusCode != successStatus {
		return resp, c.err("CompleteUpload", "expected status %d, got: %d filename: %s", successStatus, res.StatusCode, u.Path)
	}
	return resp, nil
}

// newRequest wraps http.NewRequest and adds essential headers for any BrickFTP
// API request. These include an authentication cookie and content types.
func (c Client) newRequest(method, url string, request interface{}) (*http.Request, error) {
	var body io.Reader
	if request != nil {
		encoded, err := json.Marshal(request)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(encoded)
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	auth := basicAuth(c.APIKey, "x")
	req.Header.Set("Authorization", "Basic "+auth)
	return req, nil
}

func (c Client) doRequest(req *http.Request, response interface{}) (*http.Response, error) {
	resp, err := c.Client.Do(req)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, err
	}
	if response == nil {
		// if response object is nil, we don't need to unmarshal
		return resp, nil
	}
	err = json.Unmarshal(body, response)
	return resp, err
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func (c Client) err(endpoint string, message string, args ...interface{}) error {
	return fmt.Errorf("BrickFTP ["+endpoint+"]: "+message, args...)
}

func (c Client) filesURL(file string) string {
	u := c.BaseURL
	u.Path = path.Join("/api/rest/v1/files", file)
	return u.String()
}
