package rubrik_cdm

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func Get(version, endpoint string) map[string]interface{} {
	// Get

	username := ""
	passwd := ""

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// url := "https://172.21.8.53/api/v1/cluster/me"
	url := fmt.Sprintf("https://172.21.8.53/api/%s%s", version, endpoint)

	request, err := http.NewRequest("GET", url, nil)
	request.SetBasicAuth(username, passwd)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	api_response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(api_response.Body)

	json_body := []byte(body)

	map_body := map[string]interface{}{}

	if err := json.Unmarshal(json_body, &map_body); err != nil {

		panic(err)
	}

	return map_body

}
