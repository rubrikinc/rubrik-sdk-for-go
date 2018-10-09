package rubrikcdm

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

// Type and Constants are used for escaping Get requests
type encoding int

const (
	encodePath encoding = 1 + iota
	encodePathSegment
	encodeQueryComponent
)

// Credentials - Contains parameters used to authenticate against the Rubrik cluster
type Credentials struct {
	NodeIP   string
	Username string
	Password string
}

// Connect initializes a new API client based on manually provided credentials.
func Connect(nodeIP, username, password string) *Credentials {
	client := &Credentials{
		NodeIP:   nodeIP,
		Username: username,
		Password: password,
	}

	return client
}

// ConnectEnv initializes a new API client based on environment variables.
func ConnectEnv() *Credentials {

	nodeIP, ok := os.LookupEnv("rubrik_cdm_node_ip")
	if ok != true {
		log.Fatalf("Error: The `rubrik_cdm_node_ip` environment variable is not present.")
	}
	username, ok := os.LookupEnv("rubrik_cdm_username")
	if ok != true {
		log.Fatalf("Error: The `rubrik_cdm_username` environment variable is not present.")
	}
	password, ok := os.LookupEnv("rubrik_cdm_password")
	if ok != true {
		log.Fatalf("Error: The `rubrik_cdm_password` environment variable is not present.")
	}

	client := &Credentials{
		NodeIP:   nodeIP,
		Username: username,
		Password: password,
	}

	return client
}

// Consolidate the base API functions.
func (c *Credentials) commonAPI(callType, apiVersion, apiEndpoint string, config map[string]string, timeout int) map[string]interface{} {

	if apiVersionValidation(apiVersion) == false {
		log.Fatalf("Error: Enter a valid API version.")
	}

	if endpointValidation(apiEndpoint) == "errorStart" {
		log.Fatalf("Error: The API Endpoint should begin with '/' (ex: /cluster/me).")
	} else if endpointValidation(apiEndpoint) == "errorEnd" {
		log.Fatal("Error: The API Endpoint should not end with '/' (ex. /cluster/me).")
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
	}
	request.SetBasicAuth(c.Username, c.Password)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	apiRequest, err := client.Do(request)
	if err, ok := err.(net.Error); ok && err.Timeout() {
		log.Fatalf("Error: Unable to establish a connection to the Rubrik cluster.")
	} else if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(apiRequest.Body)

	apiResponse := []byte(body)

	mapAPIResponse := map[string]interface{}{}

	if err := json.Unmarshal(apiResponse, &mapAPIResponse); err != nil {

		// DELETE request will return a 204 No Content status
		if apiRequest.StatusCode == 204 {
			mapAPIResponse["statusCode"] = apiRequest.StatusCode
		} else {
			log.Fatalf("Error: %s", apiRequest.Status)
		}

	}

	for k := range mapAPIResponse {
		if k == "errorType" || k == "message" {
			log.Fatalf("Error: %s", mapAPIResponse["message"])
		}
	}
	return mapAPIResponse

}

// Validate the API Version provided in the Base API functions
func apiVersionValidation(apiVersion string) bool {
	validAPIVersions := []string{"v1", "internal"}

	for _, version := range validAPIVersions {
		if version == apiVersion {
			return true
		}
	}
	return false
}

// Validate the endpoint provided in the Base API functions
func endpointValidation(apiEndpoint string) string {

	if string(apiEndpoint[0]) != "/" {
		return "errorStart"
	} else if string(apiEndpoint[len(apiEndpoint)-1]) == "/" {
		return "errorEnd"
	}
	return "succes"
}

// Return a default timeout value of 15 or use the value provided by the end user
func httpTimeout(timeout []int) int {
	if len(timeout) == 0 {
		return int(15) // if not timeout value is provided, set the default to 15
	}
	return int(timeout[0]) // set the timeout value to the first value in the timeout slice

}

// Custom implementation of url.PathEscape
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

// Get - Send a GET request to the provided Rubrik API endpoint.
func (c *Credentials) Get(apiVersion, apiEndpoint string, timeout ...int) map[string]interface{} {

	httpTimeout := httpTimeout(timeout)

	return c.commonAPI("GET", apiVersion, apiEndpoint, nil, httpTimeout)

}

// Post - Send a POST request to the provided Rubrik API endpoint.
func (c *Credentials) Post(apiVersion, apiEndpoint string, config map[string]string, timeout ...int) map[string]interface{} {

	httpTimeout := httpTimeout(timeout)

	return c.commonAPI("POST", apiVersion, apiEndpoint, config, httpTimeout)
}

// Patch - Send a PATCH request to the provided Rubrik API endpoint.
func (c *Credentials) Patch(apiVersion, apiEndpoint string, config map[string]string, timeout ...int) map[string]interface{} {

	httpTimeout := httpTimeout(timeout)

	return c.commonAPI("PATCH", apiVersion, apiEndpoint, config, httpTimeout)
}

// Delete - Send a DELETE request to the provided Rubrik API endpoint.
func (c *Credentials) Delete(apiVersion, apiEndpoint string, timeout ...int) map[string]interface{} {

	httpTimeout := httpTimeout(timeout)

	return c.commonAPI("DELETE", apiVersion, apiEndpoint, nil, httpTimeout)
}
