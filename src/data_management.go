package rubrikcdm

import (
	"fmt"
	"log"
)

// ObjectID will search the Rubrik cluster for the provided objectName and return its Id
func (c *Credentials) ObjectID(objectType, objectName string) string {

	validObjectType := map[string]bool{
		"vmware":     true,
		"sla":        true,
		"vmwareHost": true,
	}

	if validObjectType[objectType] == false {
		log.Fatalf("Error: The `objectType` must be `vmware`, `sla`, or `vmwareHost`.")
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
	}

	apiRequest := c.Get(objectSummaryAPIVersion, objectSummaryAPIEndpoint)

	if apiRequest["total"] == 0 {
		log.Fatalf(fmt.Sprintf("Error: The %s object '%s' was not found on the Rubrik cluster.", objectType, objectName))
	} else {

		for _, v := range apiRequest["data"].([]interface{}) {
			if v.(interface{}).(map[string]interface{})["name"].(string) == objectName {
				return v.(interface{}).(map[string]interface{})["id"].(string)
			}
		}
	}

	log.Fatalf(fmt.Sprintf("Error: The %s object '%s' was not found on the Rubrik cluster.", objectType, objectName))
	return ""

}
