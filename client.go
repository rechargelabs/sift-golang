//
//

package sift

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// Response -
type Response struct {
	// Used for debugging purposes I guess....
	HTTPStatus       string      `json:"-"`
	HTTPStatusCode   int         `json:"-"`
	HTTPStatusHeader http.Header `json:"-"`
	HTTPResponseBody string      `json:"-"`
	// ------------------------------------------------------------
	Status       int           `json:"status"`
	ErrorMessage string        `json:"error_message"`
	Time         time.Duration `json:"time"`
	Request      string        `json:"request"`
}

// IsOK - Check status of response. Is it error'ed or succeed?
func (r *Response) IsOK() bool {
	if _, ok := NoContentStatusCodes[r.HTTPStatusCode]; ok {
		return 204 == r.HTTPStatusCode
	}

	return r.Status == 0
}

// Config - Configuration struct used once per Sift Environment
type Config struct {
	ApiUrl     string        `json:"api_url"`
	ApiVersion int           `json:"api_version"`
	ApiKey     string        `json:"api_key"`
	Timeout    time.Duration `json:"timeout"`
}

// Client - Designed to be as connection point between code and Sift Science API
type Client struct {
	Config `json:"config"`
}

// SetApiUrl - Set Sift API Url. Should not be modified unless you know
// what you're doing. Default API Url can be seen in constants.go
func (c *Client) SetApiUrl(url string) {
	c.ApiUrl = url
}

// SetApiKey - Set Sift API key. You can find your keys at https://siftscience.com/console/developer/api-keys
// Pay closer attention to Production/Sandbox Mode as keys are different.
func (c *Client) SetApiKey(key string) {
	c.ApiKey = key
}

// SetApiVersion - Set Sift API version. Default API Url can be seen in constants.go
func (c *Client) SetApiVersion(version int) {
	c.ApiVersion = version
}

// SetTimeout - Set new API request timeout. Default API Timeout can be
// seen in constants.go
func (c *Client) SetTimeout(timeout time.Duration) {
	c.Timeout = timeout
}

// UserAgent - Returns User Agent that will be used with request towards Sift Science
func (c *Client) UserAgent() string {
	return fmt.Sprintf("SiftScience/%d sift-golang/%s", c.ApiVersion, VERSION)
}

// GetEventsUrl -
func (c *Client) GetEventsUrl() string {
	return c.BuildApiUrl("events")
}

// GetScoreUrl -
func (c *Client) GetScoreUrl(userId string) string {
	return c.BuildApiUrl(fmt.Sprintf("score/%s", userId))
}

// GetScoreUrl -
func (c *Client) GetLabelUrl(userId string) string {
	return c.BuildApiUrl(fmt.Sprintf("users/%s/labels", userId))
}

// BuildApiUrl -
func (c *Client) BuildApiUrl(uri string) string {
	// Make sure correct API version is set. Fail-safe and to setup defaults
	if c.ApiVersion == 0 {
		c.SetApiVersion(API_VERSION)
	}

	// Make sure correct API URL is set. Fail-safe and setup defaults.
	if c.ApiUrl == "" {
		c.SetApiUrl(API_URL)
	}

	return fmt.Sprintf("%s/v%d/%s", c.ApiUrl, c.ApiVersion, uri)
}

// GetRequest -
func (c *Client) GetRequest(method string, url string, params map[string]interface{}) (*Response, error) {

	// Set this here so it acts as global configuration
	params["$api_key"] = c.ApiKey

	if _, ok := AvailableMethods[strings.ToUpper(method)]; !ok {
		return nil, fmt.Errorf("Passed request (method: %s) is not supported by Sift Science yet!", method)
	}

	b, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(b))
	req.Header.Set("User-Agent", c.UserAgent())
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)

	r := Response{
		HTTPStatus:       resp.Status,
		HTTPStatusCode:   resp.StatusCode,
		HTTPStatusHeader: resp.Header,
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &r, err
	}

	log.Println("response Body:", string(body))

	if err := json.Unmarshal([]byte(body), &r); err != nil {
		return &r, err
	}

	if r.IsOK() == false {
		return &r, errors.New(r.ErrorMessage)
	}

	return &r, nil
}
