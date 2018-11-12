package rubrikcdm

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"
)

// ClusterVersion returns the CDM version of the Rubrik cluster.
func (c *Credentials) ClusterVersion() string {
	apiRequest := c.Get("v1", "/cluster/me")
	return apiRequest.(map[string]interface{})["version"].(string)
}

// ClusterVersionCheck is used to determine if the Rubrik cluster is using running an earlier release than the provided CDM "clusterVersion".
// If the CDM version is an earlier release than the "clusterVersion", the following message error message is thrown:
//	Error: The Rubrik cluster must be running CDM version {clusterVersion} or later.
func (c *Credentials) ClusterVersionCheck(clusterVersion float64) {
	currentClusterVersion, _ := strconv.ParseFloat(c.ClusterVersion()[:3], 2)

	if currentClusterVersion < clusterVersion {
		log.Fatalf(fmt.Sprintf("Error: The Rubrik cluster must be running CDM version %.1f or later.", clusterVersion))
	}
}

// ClusterNodeIP returns all Node IPs in the Rubrik cluster.
func (c *Credentials) ClusterNodeIP() []string {
	apiRequest := c.Get("internal", "/cluster/me/node").(map[string]interface{})

	var nodeList []string

	for _, v := range apiRequest["data"].([]interface{}) {
		nodeList = append(nodeList, v.(interface{}).(map[string]interface{})["ipAddress"].(string))
	}

	return nodeList
}

// ClusterNodeName returns the name of all nodes in the Rubrik cluster.
func (c *Credentials) ClusterNodeName() []string {
	apiRequest := c.Get("internal", "/cluster/me/node").(map[string]interface{})

	var nodeName []string

	for _, v := range apiRequest["data"].([]interface{}) {
		nodeName = append(nodeName, v.(interface{}).(map[string]interface{})["id"].(string))
	}

	return nodeName
}

// EndUserAuthorization assigns an End User account privileges for a VMware virtual machine. vmware is currently the only
// supported "objectType"
//
// The function will return one of the following:
//	No change required. The End User '{endUser}' is already authorized to interact with the '{objectName}' VM.
//
//	The full API response for POST /internal/authorization/role/end_user
func (c *Credentials) EndUserAuthorization(objectName, endUser, objectType string, timeout ...int) interface{} {

	httpTimeout := httpTimeout(timeout)

	validObjectType := map[string]bool{
		"vmware": true,
	}

	if validObjectType[objectType] == false {
		log.Fatalf("Error: The 'objectType' must be 'vmware'.")
	}

	vmID := c.ObjectID(objectName, objectType)

	userLookup := c.Get("internal", fmt.Sprintf("/user?username=%s", endUser)).([]interface{})

	if len(userLookup) == 0 {
		log.Fatalf(fmt.Sprintf("Error: The Rubrik cluster does not contain a End User account named '%s'.", endUser))
	} else {
		userID := userLookup[0].(map[string]interface{})["id"]

		userAuthorization := c.Get("internal", fmt.Sprintf("/authorization/role/end_user?principals=%s", userID)).(map[string]interface{})

		authorizedObjects := userAuthorization["data"].([]interface{})[0].(map[string]interface{})["privileges"].(map[string]interface{})["restore"]

		for _, vm := range authorizedObjects.([]interface{}) {
			if vm == vmID {
				return fmt.Sprintf("No change required. The End User '%s' is already authorized to interact with the '%s' VM.", endUser, objectName)
			}
		}

		config := map[string]interface{}{}
		config["principals"] = []string{userID.(string)}
		config["privileges"] = map[string]interface{}{}
		config["privileges"].(map[string]interface{})["restore"] = []string{vmID}

		return c.Post("internal", "/authorization/role/end_user", config, httpTimeout)

	}

	return ""

}

// ConfigureTimezone provides the ability to set the time zone that is used by the Rubrik cluster which uses the specified
// time zone for time values in the web UI, all reports, SLA Domain settings, and all other time related operations.
//
// Valid timezone choices are:
//
// 	America/Anchorage, America/Araguaina, America/Barbados, America/Chicago, America/Denver, America/Los_Angeles America/Mexico_City, America/New_York,
//	America/Noronha, America/Phoenix, America/Toronto, America/Vancouver, Asia/Bangkok, Asia/Dhaka, Asia/Dubai, Asia/Hong_Kong, Asia/Karachi, Asia/Kathmandu,
//	Asia/Kolkata, Asia/Magadan, Asia/Singapore, Asia/Tokyo, Atlantic/Cape_Verde, Australia/Perth, Australia/Sydney, Europe/Amsterdam, Europe/Athens,
//	Europe/London, Europe/Moscow, Pacific/Auckland, Pacific/Honolulu, Pacific/Midway, or UTC.
//
// The function will return one of the following:
//	No change required. The Rubrik cluster is already configured with '{timezone}' as it's timezone.
//
//	The full API response for POST /v1/cluster/me
func (c *Credentials) ConfigureTimezone(timezone string, timeout ...int) interface{} {

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
		log.Fatalf("Error: The 'timezone' must be 'America/Anchorage', 'America/Araguaina', 'America/Barbados', 'America/Chicago', 'America/Denver', 'America/Los_Angeles' 'America/Mexico_City', 'America/New_York', 'America/Noronha', 'America/Phoenix', 'America/Toronto', 'America/Vancouver', 'Asia/Bangkok', 'Asia/Dhaka', 'Asia/Dubai', 'Asia/Hong_Kong', 'Asia/Karachi', 'Asia/Kathmandu', 'Asia/Kolkata', 'Asia/Magadan', 'Asia/Singapore', 'Asia/Tokyo', 'Atlantic/Cape_Verde', 'Australia/Perth', 'Australia/Sydney', 'Europe/Amsterdam', 'Europe/Athens', 'Europe/London', 'Europe/Moscow', 'Pacific/Auckland', 'Pacific/Honolulu', 'Pacific/Midway', or 'UTC'.")
	}

	clusterSummary := c.Get("v1", "/cluster/me")

	if clusterSummary.(map[string]interface{})["timezone"].(map[string]interface{})["timezone"] == timezone {
		return fmt.Sprintf("No change required. The Rubrik cluster is already configured with '%s' as it's timezone.", timezone)
	}

	config := map[string]interface{}{}
	config["timezone"] = map[string]string{}
	config["timezone"].(map[string]string)["timezone"] = timezone

	return c.Patch("v1", "/cluster/me", config, httpTimeout)

}

// ConfigureNTP provides the connection information for the NTP servers used for time synchronization.
//
// The function will return one of the following:
//	No change required. The NTP server(s) {ntpServers} has already been added to the Rubrik cluster.
//
//	The full API response for POST /internal/cluster/me/ntp_server
func (c *Credentials) ConfigureNTP(ntpServers []string, timeout ...int) interface{} {

	httpTimeout := httpTimeout(timeout)

	clusterNTP := c.Get("internal", "/cluster/me/ntp_server").(map[string]interface{})["data"]

	// Convert the clusterNTP Slice to []string
	var currentNTPServers = []string{}
	for _, server := range clusterNTP.([]interface{}) {
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
		return c.Post("internal", "/cluster/me/ntp_server", ntpServers, httpTimeout)
	}

	return fmt.Sprintf("No change required. The NTP server(s) %s has already been added to the Rubrik cluster.", ntpServers)

}

// ConfigureSyslog enables the Rubrik cluster to send syslog server messages that are based on the events that also appear in
// the Activity Log to the provided Syslog Server.
//
// Valid protocol choices are:
//
//	UDP, TCP
//
// The function will return one of the following:
//	No change required. The Rubrik cluster is already configured to use the syslog server '{syslogIP}' on port '{port}' using the '{protocol}' protocol.
//
//	The full API response for POST /internal/syslog
func (c *Credentials) ConfigureSyslog(syslogIP, protocol string, port float64, timeout ...int) interface{} {

	httpTimeout := httpTimeout(timeout)

	validProtocols := map[string]bool{
		"UDP": true,
		"TCP": true,
	}

	if validProtocols[protocol] == false {
		log.Fatalf("Error: The 'protocol' must be 'UDP' or 'TCP'.")
	}

	config := map[string]interface{}{}
	config["hostname"] = syslogIP
	config["protocol"] = protocol
	config["port"] = port

	clusterSyslog := c.Get("internal", "/syslog").(map[string]interface{})["data"]

	if len(clusterSyslog.([]interface{})) > 0 {

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
			c.Delete("internal", "/syslog/1")
		} else {
			return fmt.Sprintf("No change required. The Rubrik cluster is already configured to use the syslog server '%s' on port '%d' using the '%s' protocol.", syslogIP, int(port), protocol)
		}

	}
	return c.Post("internal", "/syslog", config, httpTimeout)

}

// ConfigureDNSServers provides the connection information for the DNS Servers used by the Rubrik cluster.
//
// The function will return one of the following:
//	No change required. The Rubrik cluster is already configured with the provided DNS servers.
//
//	The full API response for POST /internal/cluster/me/dns_nameserver
func (c *Credentials) ConfigureDNSServers(serverIP []string, timeout ...int) interface{} {

	httpTimeout := httpTimeout(timeout)

	currentDNSServers := c.Get("internal", "/cluster/me/dns_nameserver", httpTimeout).(map[string]interface{})["data"].([]interface{})

	if stringEq(serverIP, currentDNSServers) {
		return "No change required. The Rubrik cluster is already configured with the provided DNS servers."
	}

	return c.Post("internal", "/cluster/me/dns_nameserver", serverIP, httpTimeout)

}

// ConfigureSearchDomain provides the connection information for the DNS search domains used by the Rubrik cluster.
//
// The function will return one of the following:
//	No change required. The Rubrik cluster is already configured with the provided DNS search domains.
//
//	The full API response for POST /internal/cluster/me/dns_search_domain
func (c *Credentials) ConfigureSearchDomain(searchDomain []string, timeout ...int) interface{} {

	httpTimeout := httpTimeout(timeout)

	currentSearchDomains := c.Get("internal", "/cluster/me/dns_search_domain", httpTimeout).(map[string]interface{})["data"].([]interface{})

	if stringEq(searchDomain, currentSearchDomains) {
		return "No change required. The Rubrik cluster is already configured with the provided DNS search domains."
	}

	return c.Post("internal", "/cluster/me/dns_search_domain", searchDomain, httpTimeout)

}

// ConfigureSMTPSettings provides the connection information to send notification email messages for delivery to
// the administrator accounts.
//
// Valid encryption choices are:
//
//	NONE, SSL, and STARTTLS
//
// The function will return one of the following:
//	No change required. The Rubrik cluster is already configured with the provided SMTP settings.
//
//	The full API response for POST /internal/smtp_instance
//
// The full API response for PATCH /smtp_instance/{smtpID}
func (c *Credentials) ConfigureSMTPSettings(hostname, fromEmail, smtpUsername, smtpPassword, encryption string, port int, timeout ...int) interface{} {

	httpTimeout := httpTimeout(timeout)

	validEncryption := map[string]bool{
		"NONE":     true,
		"SSL":      true,
		"STARTTLS": true,
	}

	if validEncryption[encryption] == false {
		log.Fatalf("Error: The 'encryption' must be 'NONE', 'SSL' or 'STARTTLS'.")
	}

	config := map[string]interface{}{}
	config["smtpSecurity"] = encryption
	config["smtpHostname"] = hostname
	config["smtpPort"] = port
	config["smtpUsername"] = smtpUsername
	config["fromEmailId"] = fromEmail

	getSMTPSettings := c.Get("internal", "/smtp_instance", httpTimeout)

	if getSMTPSettings.(map[string]interface{})["total"] == float64(0) {
		config["smtpPassword"] = smtpPassword
		return c.Post("internal", "/smtp_instance", config, httpTimeout)
	}

	currentSMTPSettings := getSMTPSettings.(map[string]interface{})["data"].([]interface{})[0]

	smtpID := currentSMTPSettings.(map[string]interface{})["id"]
	delete(currentSMTPSettings.(map[string]interface{}), "id")
	// Convert the smtpPort to int for comparison
	currentSMTPSettings.(map[string]interface{})["smtpPort"] = int(currentSMTPSettings.(map[string]interface{})["smtpPort"].(float64))

	checkConfig := reflect.DeepEqual(config, currentSMTPSettings)
	if checkConfig {
		return fmt.Sprintf("No change required. The Rubrik cluster is already configured with the provided SMTP settings.")
	}

	return c.Patch("internal", fmt.Sprintf("/smtp_instance/%s", smtpID), config, httpTimeout)
}

// ConfigureVLAN provides the VLAN VLAN tagging information which is an optional feature that allows a Rubrik cluster to
// efficiently switch network traffic using Virtual Local Area Networks. The ips map should be in a {nodeName:IP} format.
//
// The function will return one of the following:
//	No change required. The Rubrik cluster is already configured with the provided VLAN information.
//
//	The full API response for POST /internal/cluster/me/vlan
func (c *Credentials) ConfigureVLAN(netmask string, vlan int, ips map[string]string, timeout ...int) interface{} {

	httpTimeout := httpTimeout(timeout)

	config := map[string]interface{}{}
	config["vlan"] = vlan
	config["netmask"] = netmask
	nodeIPInterfaces := []interface{}{}
	for k, v := range ips {
		nodeIPInterfaces = append(nodeIPInterfaces, map[string]interface{}{"node": k, "ip": v})
	}
	config["interfaces"] = nodeIPInterfaces

	getCurrentVLANs := c.Get("internal", "/cluster/me/vlan", httpTimeout)

	if getCurrentVLANs.(map[string]interface{})["total"] != float64(0) {
		currentVLANs := getCurrentVLANs.(map[string]interface{})["data"].([]interface{})[0]
		// Convert int from float64 to int
		currentVLANs.((map[string]interface{}))["vlan"] = int(currentVLANs.((map[string]interface{}))["vlan"].(float64))

		checkConfig := reflect.DeepEqual(config, currentVLANs)
		if checkConfig {
			return "No change required. The Rubrik cluster is already configured with the provided VLAN information."
		}

	}
	return c.Post("internal", "/cluster/me/vlan", config, httpTimeout)

}

// AddvCenter connects to the Rubrik cluster to a new vCenter instance.
//
// The function will return one of the following:
//	No change required. The vCenter '{vcenterIP}' has already been added to the Rubrik cluster.
//
//	The full API response for POST /v1/vmware/vcenter
func (c *Credentials) AddvCenter(vCenterIP, vCenterUsername, vCenterPassword string, vmLinking bool, timeout ...int) string {

	httpTimeout := httpTimeout(timeout)

	currentVCenter := c.Get("v1", "/vmware/vcenter?primary_cluster_id=local", httpTimeout).(map[string]interface{})

	for _, v := range currentVCenter["data"].([]interface{}) {

		if v.(interface{}).(map[string]interface{})["hostname"].(string) == vCenterIP {
			return fmt.Sprintf("No change required. The vCenter '%s' has already been added to the Rubrik cluster.", vCenterIP)
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

	return c.Post("v1", "/vmware/vcenter", config, httpTimeout).(map[string]interface{})["links"].([]interface{})[0].(map[string]interface{})["href"].(string)

}

// AddvCenterWithCert connects to the Rubrik cluster to a new vCenter instance using a CA certificate.
//
// The function will return one of the following:
//	No change required. The vCenter '{vcenterIP}' has already been added to the Rubrik cluster.
//
//	The full API response for POST /v1/vmware/vcenter
func (c *Credentials) AddvCenterWithCert(vCenterIP, vCenterUsername, vCenterPassword, caCertificate string, vmLinking bool, timeout ...int) string {

	httpTimeout := httpTimeout(timeout)

	currentVCenter := c.Get("v1", "/vmware/vcenter?primary_cluster_id=local", httpTimeout).(map[string]interface{})

	for _, v := range currentVCenter["data"].([]interface{}) {

		if v.(interface{}).(map[string]interface{})["hostname"].(string) == vCenterIP {
			return fmt.Sprintf("No change required. The vCenter '%s' has already been added to the Rubrik cluster.", vCenterIP)
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

	return c.Post("v1", "/vmware/vcenter", config, httpTimeout).(map[string]interface{})["links"].([]interface{})[0].(map[string]interface{})["href"].(string)

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
func (c *Credentials) Bootstrap(clusterName, adminEmail, adminPassword, managementGateway, managementSubnetMask string, dnsSearchDomains, dnsNameServers, ntpServers []string, nodeConfig map[string]string, enableEncryption, waitForCompletion bool, timeout ...int) interface{} {

	httpTimeout := httpTimeout(timeout)

	// Validate that the Credentials struck only has a node ip configured.
	if len(c.Username) != 0 {
		log.Fatalf("Error: When bootstrapping a cluster the 'username' variable must be a blank string.")
	}

	if len(c.Password) != 0 {
		log.Fatalf("Error: When bootstrapping a cluster the 'password' variable must be a blank string.")
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

	if c.Get("internal", "/node_management/is_bootstrapped", httpTimeout).(map[string]interface{})["value"].(bool) == true {
		return "The provided Rubrik node is already bootstrapped."
	}

	bootstrap := c.Post("internal", "/cluster/me/bootstrap", config, httpTimeout)
	bootstrapRequestID := bootstrap.(map[string]interface{})["id"].(float64)

	if waitForCompletion {

		for {

			bootstrapStatus := c.Get("internal", fmt.Sprintf("/cluster/me/bootstrap?request_id=%v", int(bootstrapRequestID)), httpTimeout)

			switch bootstrapStatus.(map[string]interface{})["status"] {
			case "IN_PROGRESS":
				fmt.Println("In Progress")
				time.Sleep(30 * time.Second)
			case "FAILURE":
				fmt.Println("Failure")
				log.Fatalf("Error: %s", bootstrapStatus.(map[string]interface{})["message"])
			case "FAILED":
				fmt.Println("Failed")
				log.Fatalf("Error: %s", bootstrapStatus.(map[string]interface{})["message"])
			default:
				fmt.Println("Default")
				return bootstrapStatus

			}

		}
	}

	return bootstrap
}
