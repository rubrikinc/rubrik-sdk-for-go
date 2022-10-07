// Copyright 2018 Rubrik, Inc.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License prop
//  http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rubrikcdm

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
)

// ClusterVersion corresponds to /v1/cluster/me/version
type ClusterVersion struct {
	Version string `json:"version"`
}

// EndUserAuthorization corresponds to POST /internal/authorization/role/end_user
type EndUserAuthorization struct {
	HasMore bool `json:"hasMore"`
	Data    []struct {
		Principal  string `json:"principal"`
		Privileges struct {
			DestructiveRestore []string `json:"destructiveRestore"`
			Restore            []string `json:"restore"`
			ProvisionOnInfra   []string `json:"provisionOnInfra"`
		} `json:"privileges"`
		OrganizationID string `json:"organizationId"`
	} `json:"data"`
	Total int `json:"total"`
}

// ClusterProperties corresponds to PATCH /v1/cluster/{id}
type ClusterProperties struct {
	ID         string `json:"id"`
	Version    string `json:"version"`
	APIVersion string `json:"apiVersion"`
	Name       string `json:"name"`
	Timezone   struct {
		Timezone string `json:"timezone"`
	} `json:"timezone"`
	Geolocation struct {
		Address string `json:"address"`
	} `json:"geolocation"`
	AcceptedEulaVersion string `json:"acceptedEulaVersion"`
	LatestEulaVersion   string `json:"latestEulaVersion"`
}

// StatusCode is used when the only API response is a status code
type StatusCode struct {
	StatusCode int `json:"statusCode"`
}

// Syslog corresponds to POST /internal/syslog
type Syslog struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	ID       string `json:"id"`
}

// SMTP corresponds to PATCH /internal/smtp_instance/{id}
type SMTP struct {
	ID           string `json:"id"`
	SMTPHostname string `json:"smtpHostname"`
	SMTPPort     int    `json:"smtpPort"`
	SMTPSecurity string `json:"smtpSecurity"`
	SMTPUsername string `json:"smtpUsername"`
	FromEmailID  string `json:"fromEmailId"`
}

// ClusterVersion returns the CDM version of the Rubrik cluster.
func (c *Credentials) ClusterVersion(timeout ...int) (string, error) {

	httpTimeout := httpTimeout(timeout)

	apiRequest, err := c.Get("v1", "/cluster/me/version", httpTimeout)
	if err != nil {
		return "", err
	}

	// Convert the API Response (map[string]interface{}) to a struct
	var currentClusterVersion ClusterVersion
	mapErr := mapstructure.Decode(apiRequest, &currentClusterVersion)
	if mapErr != nil {
		return "", mapErr
	}

	return currentClusterVersion.Version, nil
}

// ClusterVersionCheck is used to determine if the Rubrik cluster is using running an earlier release than the provided CDM "clusterVersion".
// If the CDM version is an earlier release than the "clusterVersion", the following message error message is thrown:
// Error: The Rubrik cluster must be running CDM version {clusterVersion} or later.
func (c *Credentials) ClusterVersionCheck(clusterVersion float64, timeout ...int) error {

	httpTimeout := httpTimeout(timeout)

	currentClusterVersion, err := c.ClusterVersion(httpTimeout)
	if err != nil {
		return err
	}

	convertedClusterVersion, _ := strconv.ParseFloat(currentClusterVersion[:3], 2)

	if convertedClusterVersion < clusterVersion {
		return fmt.Errorf("The Rubrik cluster must be running CDM version %.1f or later", clusterVersion)
	}

	return nil
}

// ClusterNodeIP returns all Node IPs in the Rubrik cluster.
func (c *Credentials) ClusterNodeIP(timeout ...int) ([]string, error) {

	httpTimeout := httpTimeout(timeout)

	apiRequest, err := c.Get("internal", "/cluster/me/node", httpTimeout)
	if err != nil {
		return nil, err
	}

	var nodeList []string

	for _, v := range apiRequest.(map[string]interface{})["data"].([]interface{}) {
		nodeList = append(nodeList, v.(interface{}).(map[string]interface{})["ipAddress"].(string))
	}

	return nodeList, nil
}

// ClusterNodeName returns the name of all nodes in the Rubrik cluster.
func (c *Credentials) ClusterNodeName(timeout ...int) ([]string, error) {

	httpTimeout := httpTimeout(timeout)

	apiRequest, err := c.Get("internal", "/cluster/me/node", httpTimeout)
	if err != nil {
		return nil, err
	}

	var nodeName []string

	for _, v := range apiRequest.(map[string]interface{})["data"].([]interface{}) {
		nodeName = append(nodeName, v.(interface{}).(map[string]interface{})["id"].(string))
	}

	return nodeName, nil
}

// ClusterBootstrapStatus checks whether the cluster has been bootstrapped.
func (c *Credentials) ClusterBootstrapStatus(timeout ...int) (bool, error) {

	httpTimeout := httpTimeout(timeout)

	numberOfAttempts := 0
	for {
		numberOfAttempts++
		apiRequest, err := c.Get("internal", "/node_management/is_bootstrapped", httpTimeout)
		if err != nil {

			// Give the cluster 4 minutes to start responding to API calls before returning an error
			if strings.Contains(err.Error(), "tcp") {
				if numberOfAttempts == 24 {
					return false, err
				}
			} else if strings.Contains(err.Error(), "Unable to establish a connection") {

				if numberOfAttempts == 6 {
					return false, err
				}

			} else {
				return false, err

			}
		}

		if err == nil {
			return apiRequest.(map[string]interface{})["value"].(bool), nil
		}
		time.Sleep(10 * time.Second)

	}

}

// EndUserAuthorization assigns an End User account privileges for a VMware virtual machine. VMware is currently the only
// supported "objectType"
//
// The function will return one of the following:
//
//	No change required. The End User '{endUser}' is already authorized to interact with the '{objectName}' VM.
//
//	The full API response for POST /internal/authorization/role/end_user
func (c *Credentials) EndUserAuthorization(objectName, endUser, objectType string, timeout ...int) (*EndUserAuthorization, error) {

	httpTimeout := httpTimeout(timeout)

	validObjectType := map[string]bool{
		"VMware": true,
	}

	if validObjectType[objectType] == false {
		return nil, errors.New("The 'objectType' must be 'VMware'")
	}

	vmID, err := c.ObjectID(objectName, objectType, httpTimeout)
	if err != nil {
		return nil, err
	}

	userLookup, err := c.Get("internal", fmt.Sprintf("/user?username=%s", endUser))
	if err != nil {
		return nil, err
	}

	if len(userLookup.([]interface{})) == 0 {
		return nil, fmt.Errorf("The Rubrik cluster does not contain a End User account named '%s'", endUser)
	}
	userID := userLookup.([]interface{})[0].(map[string]interface{})["id"]

	userAuthorization, err := c.Get("internal", fmt.Sprintf("/authorization/role/end_user?principals=%s", userID))
	if err != nil {
		return nil, err
	}

	authorizedObjects := userAuthorization.(map[string]interface{})["data"].([]interface{})[0].(map[string]interface{})["privileges"].(map[string]interface{})["restore"]

	for _, vm := range authorizedObjects.([]interface{}) {
		if vm == vmID {
			return nil, fmt.Errorf("No change required. The End User '%s' is already authorized to interact with the '%s' VM", endUser, objectName)
		}
	}

	config := map[string]interface{}{}
	config["principals"] = []string{userID.(string)}
	config["privileges"] = map[string]interface{}{}
	config["privileges"].(map[string]interface{})["restore"] = []string{vmID}

	apiRequest, err := c.Post("internal", "/authorization/role/end_user", config, httpTimeout)
	if err != nil {
		return nil, err
	}

	// Convert the API Response (map[string]interface{}) to a struct
	var apiResponse EndUserAuthorization
	mapErr := mapstructure.Decode(apiRequest, &apiResponse)
	if mapErr != nil {
		return nil, mapErr
	}

	return &apiResponse, nil

}

// ConfigureTimezone provides the ability to set the time zone that is used by the Rubrik cluster which uses the specified
// time zone for time values in the web UI, all reports, SLA Domain settings, and all other time related operations.
//
// Valid timezone choices are:
//
//	America/Anchorage, America/Araguaia, America/Barbados, America/Chicago, America/Denver, America/Los_Angeles America/Mexico_City, America/New_York,
//	America/Noronha, America/Phoenix, America/Toronto, America/Vancouver, Asia/Bangkok, Asia/Dhaka, Asia/Dubai, Asia/Hong_Kong, Asia/Karachi, Asia/Kathmandu,
//	Asia/Kolkata, Asia/Magadan, Asia/Singapore, Asia/Tokyo, Atlantic/Cape_Verde, Australia/Perth, Australia/Sydney, Europe/Amsterdam, Europe/Athens,
//	Europe/London, Europe/Moscow, Pacific/Auckland, Pacific/Honolulu, Pacific/Midway, or UTC.
//
// The function will return one of the following:
//
//	No change required. The Rubrik cluster is already configured with '{timezone}' as it's timezone.
//
//	The full API response for POST /v1/cluster/me
func (c *Credentials) ConfigureTimezone(timezone string, timeout ...int) (*ClusterProperties, error) {

	httpTimeout := httpTimeout(timeout)

	validObjectType := map[string]bool{
		"America/Anchorage":   true,
		"America/Araguaina":   true,
		"America/Barbados":    true,
		"America/Chicago":     true,
		"America/Denver":      true,
		"America/Los_Angeles": true,
		"America/Mexico_City": true,
		"America/New_York":    true,
		"America/Noronha":     true,
		"America/Phoenix":     true,
		"America/Toronto":     true,
		"America/Vancouver":   true,
		"Asia/Bangkok":        true,
		"Asia/Dhaka":          true,
		"Asia/Dubai":          true,
		"Asia/Hong_Kong":      true,
		"Asia/Karachi":        true,
		"Asia/Kathmandu":      true,
		"Asia/Kolkata":        true,
		"Asia/Magadan":        true,
		"Asia/Singapore":      true,
		"Asia/Tokyo":          true,
		"Atlantic/Cape_Verde": true,
		"Australia/Perth":     true,
		"Australia/Sydney":    true,
		"Europe/Amsterdam":    true,
		"Europe/Athens":       true,
		"Europe/London":       true,
		"Europe/Moscow":       true,
		"Pacific/Auckland":    true,
		"Pacific/Honolulu":    true,
		"Pacific/Midway":      true,
		"UTC":                 true,
	}

	if validObjectType[timezone] == false {
		return nil, fmt.Errorf("The 'timezone' must be 'America/Anchorage', 'America/Araguaina', 'America/Barbados', 'America/Chicago', 'America/Denver', 'America/Los_Angeles' 'America/Mexico_City', 'America/New_York', 'America/Noronha', 'America/Phoenix', 'America/Toronto', 'America/Vancouver', 'Asia/Bangkok', 'Asia/Dhaka', 'Asia/Dubai', 'Asia/Hong_Kong', 'Asia/Karachi', 'Asia/Kathmandu', 'Asia/Kolkata', 'Asia/Magadan', 'Asia/Singapore', 'Asia/Tokyo', 'Atlantic/Cape_Verde', 'Australia/Perth', 'Australia/Sydney', 'Europe/Amsterdam', 'Europe/Athens', 'Europe/London', 'Europe/Moscow', 'Pacific/Auckland', 'Pacific/Honolulu', 'Pacific/Midway', or 'UTC'")
	}

	clusterSummary, err := c.Get("v1", "/cluster/me", httpTimeout)
	if err != nil {
		return nil, err
	}

	if clusterSummary.(map[string]interface{})["timezone"].(map[string]interface{})["timezone"] == timezone {
		return nil, fmt.Errorf("No change required. The Rubrik cluster is already configured with '%s' as it's timezone", timezone)
	}

	config := map[string]interface{}{}
	config["timezone"] = map[string]string{}
	config["timezone"].(map[string]string)["timezone"] = timezone

	apiRequest, err := c.Patch("v1", "/cluster/me", config, httpTimeout)
	if err != nil {
		return nil, err
	}

	// Convert the API Response (map[string]interface{}) to a struct
	var apiResponse ClusterProperties
	mapErr := mapstructure.Decode(apiRequest, &apiResponse)
	if mapErr != nil {
		return nil, mapErr
	}

	return &apiResponse, nil

}

// ConfigureNTP provides the connection information for the NTP servers used for time synchronization.
//
// The function will return one of the following:
//
//	No change required. The NTP server(s) {ntpServers} has already been added to the Rubrik cluster.
//
//	The full API response for POST /internal/cluster/me/ntp_server
func (c *Credentials) ConfigureNTP(ntpServers []string, timeout ...int) (*StatusCode, error) {

	httpTimeout := httpTimeout(timeout)

	clusterNTP, err := c.Get("internal", "/cluster/me/ntp_server")
	if err != nil {
		return nil, err
	}

	// Convert the clusterNTP Slice to []string
	var currentNTPServers = []string{}
	for _, server := range clusterNTP.(map[string]interface{})["data"].([]interface{}) {
		currentNTPServers = append(currentNTPServers, server.(string))
	}

	updateNTP := false
	if len(ntpServers) == len(currentNTPServers) {
		for i := range ntpServers {
			if ntpServers[i] != currentNTPServers[i] {
				updateNTP = true
			}
		}
	} else {
		updateNTP = true
	}

	if updateNTP {
		apiRequest, err := c.Post("internal", "/cluster/me/ntp_server", ntpServers, httpTimeout)
		if err != nil {
			return nil, err
		}

		// Convert the API Response (map[string]interface{}) to a struct
		var apiResponse StatusCode
		mapErr := mapstructure.Decode(apiRequest, &apiResponse)
		if mapErr != nil {
			return nil, mapErr
		}

		return &apiResponse, nil
	}

	return nil, fmt.Errorf("No change required. The NTP server(s) %s has already been added to the Rubrik cluster", ntpServers)

}

// ConfigureSyslog enables the Rubrik cluster to send syslog server messages that are based on the events that also appear in
// the Activity Log to the provided Syslog Server.
//
// Valid protocol choices are:
//
//	UDP, TCP
//
// The function will return one of the following:
//
//	No change required. The Rubrik cluster is already configured to use the syslog server '{syslogIP}' on port '{port}' using the '{protocol}' protocol.
//
//	The full API response for POST /internal/syslog
func (c *Credentials) ConfigureSyslog(syslogIP, protocol string, port float64, timeout ...int) (*Syslog, error) {

	httpTimeout := httpTimeout(timeout)

	validProtocols := map[string]bool{
		"UDP": true,
		"TCP": true,
	}

	if validProtocols[protocol] == false {
		return nil, errors.New("The 'protocol' must be 'UDP' or 'TCP'")
	}

	config := map[string]interface{}{}
	config["hostname"] = syslogIP
	config["protocol"] = protocol
	config["port"] = port

	clusterSyslog, err := c.Get("internal", "/syslog")
	if err != nil {
		return nil, err
	}

	if len(clusterSyslog.(map[string]interface{})["data"].([]interface{})) > 0 {

		activeSyslog := clusterSyslog.([]interface{})[0]
		// Remove id for comparison
		delete(activeSyslog.(map[string]interface{}), "id")

		deleteSyslog := false
		if len(config) == len(activeSyslog.(map[string]interface{})) {
			for i := range config {
				if config[i] != activeSyslog.(map[string]interface{})[i] {
					deleteSyslog = true

				}
			}

		}
		if deleteSyslog {
			_, err := c.Delete("internal", "/syslog/1")
			if err != nil {
				return nil, err
			}

		} else {
			return nil, fmt.Errorf("No change required. The Rubrik cluster is already configured to use the syslog server '%s' on port '%d' using the '%s' protocol", syslogIP, int(port), protocol)
		}

	}

	apiRequest, err := c.Post("internal", "/syslog", config, httpTimeout)
	if err != nil {
		return nil, err
	}

	// Convert the API Response (map[string]interface{}) to a struct
	var apiResponse Syslog
	mapErr := mapstructure.Decode(apiRequest, &apiResponse)
	if mapErr != nil {
		return nil, mapErr
	}

	return &apiResponse, nil
}

// ConfigureDNSServers provides the connection information for the DNS Servers used by the Rubrik cluster.
//
// The function will return one of the following:
//
//	No change required. The Rubrik cluster is already configured with the provided DNS servers.
//
//	The full API response for POST /internal/cluster/me/dns_nameserver
func (c *Credentials) ConfigureDNSServers(serverIP []string, timeout ...int) (*StatusCode, error) {

	httpTimeout := httpTimeout(timeout)

	currentDNSServers, err := c.Get("internal", "/cluster/me/dns_nameserver", httpTimeout)
	if err != nil {
		return nil, err
	}

	if stringEq(serverIP, currentDNSServers.(map[string]interface{})["data"].([]interface{})) {
		return nil, errors.New("No change required. The Rubrik cluster is already configured with the provided DNS servers")
	}

	apiRequest, err := c.Post("internal", "/cluster/me/dns_nameserver", serverIP, httpTimeout)
	if err != nil {
		return nil, err
	}

	// Convert the API Response (map[string]interface{}) to a struct
	var apiResponse StatusCode
	mapErr := mapstructure.Decode(apiRequest, &apiResponse)
	if mapErr != nil {
		return nil, mapErr
	}

	return &apiResponse, nil

}

// ConfigureSearchDomain provides the connection information for the DNS search domains used by the Rubrik cluster.
//
// The function will return one of the following:
//
//	No change required. The Rubrik cluster is already configured with the provided DNS search domains.
//
//	The full API response for POST /internal/cluster/me/dns_search_domain
func (c *Credentials) ConfigureSearchDomain(searchDomain []string, timeout ...int) (*StatusCode, error) {

	httpTimeout := httpTimeout(timeout)

	currentSearchDomains, err := c.Get("internal", "/cluster/me/dns_search_domain", httpTimeout)
	if err != nil {
		return nil, err
	}

	if stringEq(searchDomain, currentSearchDomains.(map[string]interface{})["data"].([]interface{})) {
		return nil, errors.New("No change required. The Rubrik cluster is already configured with the provided DNS search domains")
	}

	apiRequest, err := c.Post("internal", "/cluster/me/dns_search_domain", searchDomain, httpTimeout)
	if err != nil {
		return nil, err
	}

	// Convert the API Response (map[string]interface{}) to a struct
	var apiResponse StatusCode
	mapErr := mapstructure.Decode(apiRequest, &apiResponse)
	if mapErr != nil {
		return nil, mapErr
	}

	return &apiResponse, nil
}

// ConfigureSMTPSettings provides the connection information to send notification email messages for delivery to
// the administrator accounts.
//
// Valid encryption choices are:
//
//	NONE, SSL, and STARTTLS
//
// The function will return one of the following:
//
//	No change required. The Rubrik cluster is already configured with the provided SMTP settings.
//
//	The full API response for POST /internal/smtp_instance
//
// The full API response for PATCH /smtp_instance/{smtpID}
func (c *Credentials) ConfigureSMTPSettings(hostname, fromEmail, smtpUsername, smtpPassword, encryption string, port int, timeout ...int) (*SMTP, error) {

	httpTimeout := httpTimeout(timeout)

	validEncryption := map[string]bool{
		"NONE":     true,
		"SSL":      true,
		"STARTTLS": true,
	}

	if validEncryption[encryption] == false {
		return nil, errors.New("The 'encryption' must be 'NONE', 'SSL' or 'STARTTLS'")
	}

	config := map[string]interface{}{}
	config["smtpSecurity"] = encryption
	config["smtpHostname"] = hostname
	config["smtpPort"] = port
	config["smtpUsername"] = smtpUsername
	config["fromEmailId"] = fromEmail

	getSMTPSettings, err := c.Get("internal", "/smtp_instance", httpTimeout)
	if err != nil {
		return nil, err
	}

	if getSMTPSettings.(map[string]interface{})["total"] == float64(0) {
		config["smtpPassword"] = smtpPassword
		apiRequest, err := c.Post("internal", "/smtp_instance", config, httpTimeout)
		if err != nil {
			return nil, err
		}

		// Convert the API Response (map[string]interface{}) to a struct
		var apiResponse SMTP
		mapErr := mapstructure.Decode(apiRequest, &apiResponse)
		if mapErr != nil {
			return nil, mapErr
		}

		return &apiResponse, nil

	}

	currentSMTPSettings := getSMTPSettings.(map[string]interface{})["data"].([]interface{})[0]

	smtpID := currentSMTPSettings.(map[string]interface{})["id"]
	delete(currentSMTPSettings.(map[string]interface{}), "id")
	// Convert the smtpPort to int for comparison
	currentSMTPSettings.(map[string]interface{})["smtpPort"] = int(currentSMTPSettings.(map[string]interface{})["smtpPort"].(float64))

	checkConfig := reflect.DeepEqual(config, currentSMTPSettings)
	if checkConfig {
		return nil, fmt.Errorf("No change required. The Rubrik cluster is already configured with the provided SMTP settings")
	}

	apiRequest, err := c.Patch("internal", fmt.Sprintf("/smtp_instance/%s", smtpID), config, httpTimeout)
	if err != nil {
		return nil, err
	}

	// Convert the API Response (map[string]interface{}) to a struct
	var apiResponse SMTP
	mapErr := mapstructure.Decode(apiRequest, &apiResponse)
	if mapErr != nil {
		return nil, mapErr
	}

	return &apiResponse, nil

}

// ConfigureVLAN provides the VLAN VLAN tagging information which is an optional feature that allows a Rubrik cluster to
// efficiently switch network traffic using Virtual Local Area Networks. The ips map should be in a {nodeName:IP} format.
//
// The function will return one of the following:
//
//	No change required. The Rubrik cluster is already configured with the provided VLAN information.
//
//	The full API response for POST /internal/cluster/me/vlan
func (c *Credentials) ConfigureVLAN(netmask string, vlan int, ips map[string]string, timeout ...int) (*StatusCode, error) {

	httpTimeout := httpTimeout(timeout)

	config := map[string]interface{}{}
	config["vlan"] = vlan
	config["netmask"] = netmask
	nodeIPInterfaces := []interface{}{}
	for k, v := range ips {
		nodeIPInterfaces = append(nodeIPInterfaces, map[string]interface{}{"node": k, "ip": v})
	}
	config["interfaces"] = nodeIPInterfaces

	getCurrentVLANs, err := c.Get("internal", "/cluster/me/vlan", httpTimeout)
	if err != nil {
		return nil, err
	}

	if getCurrentVLANs.(map[string]interface{})["total"] != float64(0) {
		currentVLANs := getCurrentVLANs.(map[string]interface{})["data"].([]interface{})[0]
		// Convert int from float64 to int
		currentVLANs.((map[string]interface{}))["vlan"] = int(currentVLANs.((map[string]interface{}))["vlan"].(float64))

		checkConfig := reflect.DeepEqual(config, currentVLANs)
		if checkConfig {
			return nil, errors.New("No change required. The Rubrik cluster is already configured with the provided VLAN information")
		}

	}
	apiRequest, err := c.Post("internal", "/cluster/me/vlan", config, httpTimeout)
	if err != nil {
		return nil, err
	}

	// Convert the API Response (map[string]interface{}) to a struct
	var apiResponse StatusCode
	mapErr := mapstructure.Decode(apiRequest, &apiResponse)
	if mapErr != nil {
		return nil, mapErr
	}

	return &apiResponse, nil

}

// AddvCenter connects to the Rubrik cluster to a new vCenter instance.
//
// The function will return one of the following:
//
//	No change required. The vCenter '{vcenterIP}' has already been added to the Rubrik cluster.
//
//	The full API response for POST /v1/VMware/vcenter
func (c *Credentials) AddvCenter(vCenterIP, vCenterUsername, vCenterPassword string, vmLinking bool, timeout ...int) (interface{}, error) {

	httpTimeout := httpTimeout(timeout)

	currentVCenter, err := c.Get("v1", "/VMware/vcenter?primary_cluster_id=local", httpTimeout)
	if err != nil {
		return "", err
	}

	for _, v := range currentVCenter.(map[string]interface{})["data"].([]interface{}) {

		if v.(interface{}).(map[string]interface{})["hostname"].(string) == vCenterIP {
			return fmt.Sprintf("No change required. The vCenter '%s' has already been added to the Rubrik cluster.", vCenterIP), nil
		}
	}

	config := map[string]string{}
	config["hostname"] = vCenterIP
	config["username"] = vCenterUsername
	config["password"] = vCenterPassword
	if vmLinking {
		config["conflictResolutionAuthz"] = "AllowAutoConflictResolution"
	} else if vmLinking == false {
		config["conflictResolutionAuthz"] = "NoConflictResolution"
	}

	apiRequest, err := c.Post("v1", "/VMware/vcenter", config, httpTimeout)
	if err != nil {
		return "", err
	}

	// Convert the API Response (map[string]interface{}) to a struct
	var addVcenter JobStatus
	mapErr := mapstructure.Decode(apiRequest, &addVcenter)
	if mapErr != nil {
		return nil, mapErr
	}

	status, err := c.JobStatus(addVcenter.Links[0].Href)
	if err != nil {
		return nil, err
	}

	return status, nil

}

// AddvCenterWithCert connects to the Rubrik cluster to a new vCenter instance using a CA certificate.
//
// The function will return one of the following:
//
//	No change required. The vCenter '{vcenterIP}' has already been added to the Rubrik cluster.
//
//	The full API response for POST /v1/VMware/vcenter
func (c *Credentials) AddvCenterWithCert(vCenterIP, vCenterUsername, vCenterPassword, caCertificate string, vmLinking bool, timeout ...int) (interface{}, error) {

	httpTimeout := httpTimeout(timeout)

	currentVCenter, err := c.Get("v1", "/VMware/vcenter?primary_cluster_id=local", httpTimeout)
	if err != nil {
		return "", err
	}

	for _, v := range currentVCenter.(map[string]interface{})["data"].([]interface{}) {

		if v.(interface{}).(map[string]interface{})["hostname"].(string) == vCenterIP {
			return fmt.Sprintf("No change required. The vCenter '%s' has already been added to the Rubrik cluster.", vCenterIP), nil
		}
	}

	config := map[string]string{}
	config["hostname"] = vCenterIP
	config["username"] = vCenterUsername
	config["password"] = vCenterPassword
	if vmLinking {
		config["conflictResolutionAuthz"] = "AllowAutoConflictResolution"
	} else if vmLinking == false {
		config["conflictResolutionAuthz"] = "NoConflictResolution"
	}
	config["caCerts"] = caCertificate

	apiRequest, err := c.Post("v1", "/VMware/vcenter", config, httpTimeout)
	if err != nil {
		return "", err
	}

	// Convert the API Response (map[string]interface{}) to a struct
	var addVcenter JobStatus
	mapErr := mapstructure.Decode(apiRequest, &addVcenter)
	if mapErr != nil {
		return nil, mapErr
	}

	status, err := c.JobStatus(addVcenter.Links[0].Href)
	if err != nil {
		return nil, err
	}

	return status, nil

}

// Bootstrap will complete the bootstrap process for a Rubrik cluster and requires a single node to have it's management interface
// configured. You will also need to use Connect() with the "username" and "password" set to blank strings. The "nodeConfig" should be in a
// {nodeName: nodeManagementIP} format. To monitor the bootstrap process and wait for the process to complete, set "waitForCompletion" to true.
//
// The function will return one of the following:
//
//	The full API response for POST /internal/cluster/me/bootstrap?request_id={requestID} (waitForCompletion is set to true)
//
//	The full API response for POST /internal/cluster/me/bootstrap (waitForCompletion is set to false)
func (c *Credentials) Bootstrap(clusterName, adminEmail, adminPassword, managementGateway, managementSubnetMask string, dnsSearchDomains, dnsNameServers, ntpServers []string, nodeConfig map[string]string, enableEncryption, waitForCompletion bool, timeout ...int) (interface{}, error) {

	httpTimeout := httpTimeout(timeout)

	// Validate that the Credentials struck only has a node ip configured.
	if len(c.Username) != 0 {
		return nil, errors.New("When bootstrapping a cluster the 'username' variable must be a blank string")
	}

	if len(c.Password) != 0 {
		return nil, errors.New("When bootstrapping a cluster the 'password' variable must be a blank string")
	}

	config := map[string]interface{}{}
	config["enableSoftwareEncryptionAtRest"] = enableEncryption
	config["name"] = clusterName
	config["dnsNameservers"] = dnsNameServers
	config["dnsSearchDomains"] = dnsSearchDomains
	config["ntpServers"] = ntpServers

	config["adminUserInfo"] = map[string]string{}
	config["adminUserInfo"].(map[string]string)["password"] = adminPassword
	config["adminUserInfo"].(map[string]string)["emailAddress"] = adminEmail
	config["adminUserInfo"].(map[string]string)["id"] = "admin"

	config["nodeConfigs"] = map[string]interface{}{}
	for nodeName, nodeIP := range nodeConfig {
		config["nodeConfigs"].(map[string]interface{})[nodeName] = map[string]interface{}{}
		config["nodeConfigs"].(map[string]interface{})[nodeName].(map[string]interface{})["managementIpConfig"] = map[string]string{}
		config["nodeConfigs"].(map[string]interface{})[nodeName].(map[string]interface{})["managementIpConfig"].(map[string]string)["netmask"] = managementSubnetMask
		config["nodeConfigs"].(map[string]interface{})[nodeName].(map[string]interface{})["managementIpConfig"].(map[string]string)["gateway"] = managementGateway
		config["nodeConfigs"].(map[string]interface{})[nodeName].(map[string]interface{})["managementIpConfig"].(map[string]string)["address"] = nodeIP
	}

	currentBootstrapStatus, err := c.ClusterBootstrapStatus(httpTimeout)
	if err != nil {
		return nil, err
	}

	if currentBootstrapStatus == true {
		return "The provided Rubrik node is already bootstrapped.", nil
	}
	bootstrap, err := c.Post("internal", "/cluster/me/bootstrap", config, httpTimeout)
	if err != nil {

		return nil, err
	}
	bootstrapRequestID := bootstrap.(map[string]interface{})["id"].(float64)

	if waitForCompletion {

		for {

			bootstrapStatus, err := c.Get("internal", fmt.Sprintf("/cluster/me/bootstrap?request_id=%v", int(bootstrapRequestID)), httpTimeout)
			if err != nil {
				return nil, err
			}

			switch bootstrapStatus.(map[string]interface{})["status"] {
			case "IN_PROGRESS":
				time.Sleep(30 * time.Second)
			case "FAILURE":
				return nil, fmt.Errorf("%s", bootstrapStatus.(map[string]interface{})["message"])
			case "FAILED":
				return nil, fmt.Errorf("%s", bootstrapStatus.(map[string]interface{})["message"])

			default:
				return bootstrapStatus, nil

			}

		}
	}

	return bootstrap, nil
}

// RegisterCluster submits the registration details for the specified Rubrik cluster. The username and password should
// correspond to your Rubrik Support Portal account. The default timeout value is 160 seconds.
func (c *Credentials) RegisterCluster(username, password string, timeout ...int) (interface{}, error) {

	httpTimeout := httpTimeout(timeout)

	// Change the default to 160
	if httpTimeout == 15 {
		httpTimeout = 160
	}

	isRegistered, err := c.Get("internal", "/cluster/me/is_registered")
	if err != nil {
		return nil, err
	}

	if isRegistered.(map[string]interface{})["value"] == true {
		return "No change required. The cluster is already registered.", nil
	}

	config := map[string]interface{}{}
	config["username"] = username
	config["password"] = password

	register, err := c.Post("internal", "/cluster/me/register", config, httpTimeout)
	if err != nil {
		return nil, err
	}

	return register, nil

}

// RefreshvCenter updates the the metadata for the specified vCenter Server and waits for the job to complete before returning the JobStatus API response.
func (c *Credentials) RefreshvCenter(vCenterIP string, timeout ...int) (interface{}, error) {

	httpTimeout := httpTimeout(timeout)

	vcenterID, err := c.ObjectID(vCenterIP, "vcenter", httpTimeout)
	if err != nil {
		return nil, err
	}

	refresh, err := c.Post("v1", fmt.Sprintf("/VMware/vcenter/%s/refresh", vcenterID), httpTimeout)
	if err != nil {
		return nil, err
	}

	// Convert the API Response (map[string]interface{}) to a struct
	var refreshJob JobStatus
	mapErr := mapstructure.Decode(refresh, &refreshJob)
	if mapErr != nil {
		return nil, mapErr
	}

	status, err := c.JobStatus(refreshJob.Links[0].Href)
	if err != nil {
		return nil, err
	}

	return status, nil

}
