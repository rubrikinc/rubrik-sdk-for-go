// Copyright 2018 Rubrik, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License prop
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package rubrikcdm transforms the Rubrik API functionality into easy to consume functions. This eliminates the need to understand
// how to consume raw Rubrik APIs with Go and extends upon one of Rubrik’s main design centers - simplicity. Rubrik’s API first architecture enables
// organizations to embrace and integrate Rubrik functionality into their existing automation processes.
package rubrikcdm

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"reflect"
	"sort"
	"time"
)

// Type and Constants are used for escaping Get requests
type encoding int

const (
	encodePath encoding = 1 + iota
	encodePathSegment
	encodeQueryComponent
)

// Credentials contains the parameters used to authenticate against the Rubrik cluster and can be consumed
// through ConnectEnv() or Connect().
type Credentials struct {
	NodeIP   string
	Username string
	Password string
}

// Connect initializes a new API client based on manually provided Rubrik cluster credentials. When possible,
// the Rubrik credentials should not be stored as plain text in your .go file. ConnectEnv() can be used
// as a safer alternative.
func Connect(nodeIP, username, password string) *Credentials {
	client := &Credentials{
		NodeIP:   nodeIP,
		Username: username,
		Password: password,
	}

	return client
}

// ConnectEnv is the preferred method to initialize a new API client by attempting to read the
// following environment variables:
//
//  rubrik_cdm_node_ip
//
//  rubrik_cdm_username
//
//  rubrik_cdm_password
func ConnectEnv() (*Credentials, error) {

	nodeIP, ok := os.LookupEnv("rubrik_cdm_node_ip")
	if ok != true {
		return nil, errors.New("The `rubrik_cdm_node_ip` environment variable is not present")
	}
	username, ok := os.LookupEnv("rubrik_cdm_username")
	if ok != true {
		return nil, errors.New("The `rubrik_cdm_username` environment variable is not present")
	}
	password, ok := os.LookupEnv("rubrik_cdm_password")
	if ok != true {
		return nil, errors.New("The `rubrik_cdm_password` environment variable is not present")
	}

	client := &Credentials{
		NodeIP:   nodeIP,
		Username: username,
		Password: password,
	}

	return client, nil
}

// Consolidate the base API functions.
func (c *Credentials) commonAPI(callType, apiVersion, apiEndpoint string, config interface{}, timeout int) (interface{}, error) {

	if apiVersionValidation(apiVersion) == false {
		return nil, errors.New("Enter a valid API version")
	}

	if endpointValidation(apiEndpoint) == "errorStart" {
		return nil, errors.New("The API Endpoint should begin with '/' (ex: /cluster/me)")
	} else if endpointValidation(apiEndpoint) == "errorEnd" {
		return nil, errors.New("The API Endpoint should not end with '/' (ex. /cluster/me)")
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * time.Duration(timeout),
	}

	requestURL := fmt.Sprintf("https://%s/api/%s%s", c.NodeIP, apiVersion, apiEndpoint)

	var request *http.Request
	switch callType {
	case "GET":
		request, _ = http.NewRequest(callType, getEscape(requestURL), nil)
	case "POST":
		convertedConfig, _ := json.Marshal(config)
		request, _ = http.NewRequest(callType, requestURL, bytes.NewBuffer(convertedConfig))
	case "PATCH":
		convertedConfig, _ := json.Marshal(config)
		request, _ = http.NewRequest(callType, requestURL, bytes.NewBuffer(convertedConfig))
	case "DELETE":
		request, _ = http.NewRequest(callType, requestURL, nil)
	case "JOB_STATUS":
		// Overwrite the default requstURL with the job status url and convert to string
		requestURL = config.(string)
		request, _ = http.NewRequest("GET", requestURL, nil)
	}
	if len(c.Username) != 0 {
		request.SetBasicAuth(c.Username, c.Password)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	apiRequest, err := client.Do(request)
	if err, ok := err.(net.Error); ok && err.Timeout() {
		return nil, errors.New("Unable to establish a connection to the Rubrik cluster")
	} else if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(apiRequest.Body)

	apiResponse := []byte(body)

	var convertedAPIResponse interface{}
	if err := json.Unmarshal(apiResponse, &convertedAPIResponse); err != nil {

		// DELETE request will return a 204 No Content status
		if apiRequest.StatusCode == 204 {
			convertedAPIResponse = map[string]interface{}{}
			convertedAPIResponse.(map[string]interface{})["statusCode"] = apiRequest.StatusCode
		} else if apiRequest.StatusCode != 200 {
			return nil, fmt.Errorf("%s", apiRequest.Status)
		}

	}

	//
	if reflect.TypeOf(convertedAPIResponse).Kind() == reflect.Slice {
		return convertedAPIResponse, nil
	}

	if _, ok := convertedAPIResponse.(map[string]interface{})["errorType"]; ok {
		return nil, fmt.Errorf("%s", convertedAPIResponse.(map[string]interface{})["message"])
	}

	if _, ok := convertedAPIResponse.(map[string]interface{})["message"]; ok {
		// Add exception for bootstrap
		if _, ok := convertedAPIResponse.(map[string]interface{})["setupEncryptionAtRest"]; ok {
			return convertedAPIResponse, nil

		}

		return nil, fmt.Errorf("%s", convertedAPIResponse.(map[string]interface{})["message"])
	}

	return convertedAPIResponse, nil

}

// apiVersionValidation validates the API Version provided in the Base API functions. Valid versions are v1, v2 and internal.
func apiVersionValidation(apiVersion string) bool {
	validAPIVersions := []string{"v1", "v2", "internal"}

	for _, version := range validAPIVersions {
		if version == apiVersion {
			return true
		}
	}
	return false
}

// endpointValidation validates that the endpoint provided in the Base API functions starts with a / but does not end with one.
func endpointValidation(apiEndpoint string) string {

	if string(apiEndpoint[0]) != "/" {
		return "errorStart"
	} else if string(apiEndpoint[len(apiEndpoint)-1]) == "/" {
		return "errorEnd"
	}
	return "success"
}

// httpTimeout returns a default timeout value of 15 or use the value provided by the end user
func httpTimeout(timeout []int) int {
	if len(timeout) == 0 {
		return int(15) // if not timeout value is provided, set the default to 15
	}
	return int(timeout[0]) // set the timeout value to the first value in the timeout slice

}

// getEscape is a custom implementation of url.PathEscape.
func getEscape(s string) string {
	return escape(s, encodePathSegment)
}

func escape(s string, mode encoding) string {
	spaceCount, hexCount := 0, 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if shouldEscape(c, mode) {
			if c == ' ' && mode == encodeQueryComponent {
				spaceCount++
			} else {
				hexCount++
			}
		}
	}

	if spaceCount == 0 && hexCount == 0 {
		return s
	}

	t := make([]byte, len(s)+2*hexCount)
	j := 0
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case c == ' ' && mode == encodeQueryComponent:
			t[j] = '+'
			j++
		case shouldEscape(c, mode):
			t[j] = '%'
			t[j+1] = "0123456789ABCDEF"[c>>4]
			t[j+2] = "0123456789ABCDEF"[c&15]
			j += 3
		default:
			t[j] = s[i]
			j++
		}
	}
	return string(t)
}

func shouldEscape(c byte, mode encoding) bool {
	// §2.3 Unreserved characters (alphanum)
	if 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z' || '0' <= c && c <= '9' {
		return false
	}

	switch c {
	case '-', '_', '.', '~': // §2.3 Unreserved characters (mark)
		return false

	case '$', '&', '+', ',', '/', ':', ';', '=', '?', '@': // §2.2 Reserved characters (reserved)
		// Different sections of the URL allow a few of
		// the reserved characters to appear unescaped.

		return c == ';' || c == ','

	}

	// Everything else must be escaped.
	return true
}

// Get sends a GET request to the provided Rubrik API endpoint and returns the full API response. Supported "apiVersions" are v1, v2, and internal.
// The optional timeout value corresponds to the number of seconds to wait to establish a connection to the Rubrik cluster before returning a
// timeout error. If no value is provided, a default of 15 seconds will be used.
func (c *Credentials) Get(apiVersion, apiEndpoint string, timeout ...int) (interface{}, error) {

	httpTimeout := httpTimeout(timeout)

	apiRequest, err := c.commonAPI("GET", apiVersion, apiEndpoint, nil, httpTimeout)
	if err != nil {
		return nil, err
	}

	return apiRequest, nil

}

// JobStatus performs a GET operation to monitor the status of a specific Rubrik job want waits for it's completion.
func (c *Credentials) JobStatus(jobStatusURL string, timeout ...int) (interface{}, error) {

	httpTimeout := httpTimeout(timeout)

	// Dummy place holder values to pass validation
	apiVersion := "v1"
	apiEndpoint := "/placeholder"

	for {
		apiRequest, err := c.commonAPI("JOB_STATUS", apiVersion, apiEndpoint, jobStatusURL, httpTimeout)
		if err != nil {
			return nil, err
		}

		jobStatus := apiRequest.(map[string]interface{})["status"].(string)

		switch jobStatus {
		case "SUCCEEDED":
			return apiRequest, nil
		case "QUEUED":
			time.Sleep(10 * time.Second)
		case "RUNNING":
			time.Sleep(10 * time.Second)
		case "FINISHING":
			time.Sleep(10 * time.Second)
		default:
			return apiRequest, errors.New("Job failed")
		}
	}

}

// Post sends a POST request to the provided Rubrik API endpoint and returns the full API response. Supported "apiVersions" are v1, v2, and internal.
// The optional timeout value corresponds to the number of seconds to wait to establish a connection to the Rubrik cluster before returning a
// timeout error. If no value is provided, a default of 15 seconds will be used.
func (c *Credentials) Post(apiVersion, apiEndpoint string, config interface{}, timeout ...int) (interface{}, error) {

	httpTimeout := httpTimeout(timeout)

	apiRequest, err := c.commonAPI("POST", apiVersion, apiEndpoint, config, httpTimeout)
	if err != nil {
		return nil, err
	}

	return apiRequest, nil

}

// Patch sends a PATCH request to the provided Rubrik API endpoint and returns the full API response. Supported "apiVersions" are v1, v2, and internal.
// The optional timeout value corresponds to the number of seconds to wait to establish a connection to the Rubrik cluster before returning a
// timeout error. If no value is provided, a default of 15 seconds will be used.
func (c *Credentials) Patch(apiVersion, apiEndpoint string, config interface{}, timeout ...int) (interface{}, error) {

	httpTimeout := httpTimeout(timeout)

	apiRequest, err := c.commonAPI("PATCH", apiVersion, apiEndpoint, config, httpTimeout)
	if err != nil {
		return nil, err
	}

	return apiRequest, nil

}

// Delete sends a DELETE request to the provided Rubrik API endpoint and returns the full API response. Supported "apiVersions" are v1, v2, and internal.
// The optional timeout value corresponds to the number of seconds to wait to establish a connection to the Rubrik cluster before returning a
// timeout error. If no value is provided, a default of 15 seconds will be used.
func (c *Credentials) Delete(apiVersion, apiEndpoint string, timeout ...int) (interface{}, error) {

	httpTimeout := httpTimeout(timeout)

	apiRequest, err := c.commonAPI("DELETE", apiVersion, apiEndpoint, nil, httpTimeout)
	if err != nil {
		return nil, err
	}

	return apiRequest, nil
}

// stringEq converts b to []string, sorts the two []string, and checks for equality
func stringEq(a []string, b []interface{}) bool {

	// Convert []interface {} to []string
	c := make([]string, len(b))
	for i, v := range b {
		c[i] = fmt.Sprint(v)
	}

	sort.Strings(a)
	sort.Strings(c)

	// If one is nil, the other must also be nil.
	if (a == nil) != (c == nil) {
		return false
	}

	if len(a) != len(c) {
		return false
	}

	for i := range a {
		if a[i] != c[i] {
			return false
		}
	}

	return true
}
