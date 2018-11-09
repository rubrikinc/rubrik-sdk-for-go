package rubrikcdm

import (
	"fmt"
	"log"
)

// ObjectID will search the Rubrik cluster for the provided objectName and return its Id
func (c *Credentials) ObjectID(objectName, objectType string, hostOS ...string) string {

	validObjectType := map[string]bool{
		"vmware":          true,
		"sla":             true,
		"vmwareHost":      true,
		"physicalHost":    true,
		"filesetTemplate": true,
		"managedVolume":   true,
	}

	if validObjectType[objectType] == false {
		log.Fatalf("Error: The 'objectType' must be 'vmware', 'sla', 'vmwareHost', 'physicalHost', 'filesetTemplate', or 'managedVolume'.")
	}

	var objectSummaryAPIVersion string
	var objectSummaryAPIEndpoint string
	switch objectType {
	case "vmware":
		objectSummaryAPIVersion = "v1"
		objectSummaryAPIEndpoint = fmt.Sprintf("/vmware/vm?primary_cluster_id=local&is_relic=false&name=%s", objectName)
	case "sla":
		objectSummaryAPIVersion = "v1"
		objectSummaryAPIEndpoint = fmt.Sprintf("/sla_domain?primary_cluster_id=local&name=%s", objectName)
	case "vmwareHost":
		objectSummaryAPIVersion = "v1"
		objectSummaryAPIEndpoint = "/vmware/host?primary_cluster_id=local"
	case "filesetTemplate":
		var hostOperatingSystem string
		if len(hostOS) > 0 {
			hostOperatingSystem = hostOS[0]
			switch hostOperatingSystem {
			case "Linux":
			case "Windows":
			default:
				log.Fatalf("Error: The hostOS must be either 'Linux' or 'Windows'.")

			}
		} else if len(hostOS) == 0 {
			log.Fatalf("Error: You must provide the Fileset Tempalte OS type. ")

		}

		objectSummaryAPIVersion = "v1"
		objectSummaryAPIEndpoint = fmt.Sprintf("/fileset_template?primary_cluster_id=local&operating_system_type=%s&name=%s", hostOperatingSystem, objectName)
	case "managedVolume":
		objectSummaryAPIVersion = "internal"
		objectSummaryAPIEndpoint = fmt.Sprintf("/managed_volume?is_relic=false&primary_cluster_id=local&name=%s", objectName)
	}

	apiRequest := c.Get(objectSummaryAPIVersion, objectSummaryAPIEndpoint).(map[string]interface{})
	if apiRequest["total"] == 0 {
		log.Fatalf(fmt.Sprintf("Error: The %s object '%s' was not found on the Rubrik cluster.", objectType, objectName))
	} else if apiRequest["total"].(float64) > 0 {
		objectIDs := make([]string, 0)
		// # Define the "object name" to search for
		var nameValue string
		if objectType == "physicalHost" {
			nameValue = "hostname"
		} else {
			nameValue = "name"
		}

		for _, v := range apiRequest["data"].([]interface{}) {
			if v.(interface{}).(map[string]interface{})[nameValue].(string) == objectName {
				objectIDs = append(objectIDs, v.(interface{}).(map[string]interface{})["id"].(string))
			}
		}

		if len(objectIDs) > 1 {
			log.Fatalf(fmt.Sprintf("Error: Multiple %s objects named '%s' were found on the Rubrik cluster. Unable to return a specific object id.", objectType, objectName))
		} else if len(objectIDs) == 0 {
			log.Fatalf(fmt.Sprintf("Error: The %s object '%s' was not found on the Rubrik cluster.", objectType, objectName))
		} else {
			return objectIDs[0]
		}
	}

	log.Fatalf(fmt.Sprintf("Error: The %s object '%s' was not found on the Rubrik cluster.", objectType, objectName))
	return ""

}

// AssignSLA
func (c *Credentials) AssignSLA(objectName, objectType, slaName string, timeout ...int) interface{} {

	httpTimeout := httpTimeout(timeout)

	validObjectType := map[string]bool{
		"vmware": true,
	}

	if validObjectType[objectType] == false {
		log.Fatalf("Error: The 'objectType' must be 'vmware'.")
	}

	var slaID string
	switch slaName {
	case "do not protect":
		slaID = "UNPROTECTED"
	case "clear":
		slaID = "INHERIT"
	default:
		slaID = c.ObjectID(slaName, "sla")
	}

	config := map[string]interface{}{}
	switch objectType {
	case "vmware":
		vmID := c.ObjectID(objectName, "vmware")

		vmSummary := c.Get("v1", fmt.Sprintf("/vmware/vm/%s", vmID), httpTimeout)

		var currentSLAID string
		switch slaID {
		case "INHERIT":
			currentSLAID = vmSummary.(map[string]interface{})["configuredSlaDomainId"].(string)
		default:
			currentSLAID = vmSummary.(map[string]interface{})["effectiveSlaDomainId"].(string)
		}

		if slaID == currentSLAID {
			return fmt.Sprintf("No change required. The vSphere VM '%s' is already assigned to the '%s' SLA Domain.", objectName, slaName)
		}

		config["managedIds"] = []string{vmID}
	}

	return c.Post("internal", fmt.Sprintf("/sla_domain/%s/assign", slaID), config, httpTimeout)
}
