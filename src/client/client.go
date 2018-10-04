package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

// Contains parameters used to authenticate against the Rubrik cluster
type Connect struct {
	NodeIP   string
	Username string
	Password string
}

// Consolidate the base API functions.
func (c *Connect) commonAPI(callType, apiVersion, apiEndpoint string, timeout int) map[string]interface{} {

	if apiVersionValidation(apiVersion) == false {
		log.SetFlags(0)
		log.Fatalf("Error: Enter a valid API version.")
	}

	if endpointValidation(apiEndpoint) == "errorStart" {
		log.SetFlags(0)
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

	request, err := http.NewRequest(callType, requestURL, nil)
	request.SetBasicAuth(c.Username, c.Password)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	apiRequest, err := client.Do(request)
	if err, ok := err.(net.Error); ok && err.Timeout() {
		log.SetFlags(0)
		log.Fatalf("Error: Unable to establish a connection to the Rubrik cluster.")
	} else if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(apiRequest.Body)

	apiResponse := []byte(body)

	mapAPIResponse := map[string]interface{}{}

	if err := json.Unmarshal(apiResponse, &mapAPIResponse); err != nil {

		panic(err)
	}

	for k := range mapAPIResponse {
		if k == "errorType" || k == "message" {
			log.SetFlags(0)
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

// Get - Send a GET request to the provided Rubrik API endpoint.
func (c *Connect) Get(apiVersion, apiEndpoint string, timeout ...int) map[string]interface{} {

	httpTimeout := httpTimeout(timeout)

	return c.commonAPI("GET", apiVersion, apiEndpoint, httpTimeout)
}
