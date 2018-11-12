package rubrikcdm

import (
	"fmt"
	"log"
)

// ObjectID will search the Rubrik cluster for the provided "objectName" and return its ID/
//
// Valid "awsRegion" choices are:
//
//	vmware, sla, vmwareHost, physicalHost, filesetTemplate, managedVolume
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
	case "physicalHost":

		objectSummaryAPIVersion = "v1"
		objectSummaryAPIEndpoint = fmt.Sprintf("/host?primary_cluster_id=local&hostname=%s", objectName)
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

// AssignSLA adds the "objectName" to the "slaName". vmware is currently the only supported "objectType". To exclude the object from all SLA assignments
// use "do not protect" as the "slaName". To assign the selected object to the SLA of the next higher level object, use "clear" as the "slaName".
//
// The function will return one of the following:
//	No change required. The vSphere VM '{objectName}' is already assigned to the '{slaName}' SLA Domain.
//
//	The full API response for POST /internal/sla_domain/{slaID}/assign.
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

// BeginManagedVolumeSnapshot opens a managed volume for writes. All writes to the managed volume until the snapshot is
// ended will be part of its snapshot.
//
// The function will return one of the following:
//	No change required. The Managed Volume '{name}' is already in a writeable state.
//
//	The full API response for POST /internal/managed_volume/{managedVolumeID}/begin_snapshot
func (c *Credentials) BeginManagedVolumeSnapshot(name string, timeout ...int) interface{} {

	httpTimeout := httpTimeout(timeout)

	managedVolumeID := c.ObjectID(name, "managedVolume")

	managedVolumeSummary := c.Get("internal", fmt.Sprintf("/managed_volume/%s", managedVolumeID), httpTimeout)

	if managedVolumeSummary.(map[string]interface{})["isWritable"].(bool) {

		return fmt.Sprintf("No change required. The Managed Volume '%s' is already in a writeable state.", name)
	}

	config := map[string]string{}

	return c.Post("internal", fmt.Sprintf("/managed_volume/%s/begin_snapshot", managedVolumeID), config, httpTimeout)

}

// EndManagedVolumeSnapshot closes a managed volume for writes. A snapshot will be created containing all writes since the last begin snapshot call.
//
// The function will return one of the following:
//	No change required. The Managed Volume '{name}' is already in a read-only state.
//
//	The full API response for POST /internal/managed_volume/{managedVolumeID}/end_snapshot
func (c *Credentials) EndManagedVolumeSnapshot(name, slaName string, timeout ...int) interface{} {

	httpTimeout := httpTimeout(timeout)

	managedVolumeID := c.ObjectID(name, "managedVolume")

	managedVolumeSummary := c.Get("internal", fmt.Sprintf("/managed_volume/%s", managedVolumeID), httpTimeout)

	if managedVolumeSummary.(map[string]interface{})["isWritable"].(bool) == false {

		return fmt.Sprintf("No change required. The Managed Volume '%s' is already in a read-only state.", name)
	}

	var slaID string
	switch slaName {
	case "current":
		slaID = managedVolumeSummary.(map[string]interface{})["configuredSlaDomainId"].(string)
	default:
		slaID = c.ObjectID(slaName, "sla")
	}

	config := map[string]interface{}{}
	config["retentionConfig"] = map[string]interface{}{}
	config["retentionConfig"].(map[string]interface{})["slaId"] = slaID

	return c.Post("internal", fmt.Sprintf("/managed_volume/%s/end_snapshot", managedVolumeID), config, httpTimeout)

}

// GetSLAObjects returns the name and ID of a specific object type.
func (c *Credentials) GetSLAObjects(slaName, objectType string, timeout ...int) interface{} {

	httpTimeout := httpTimeout(timeout)

	validObjectType := map[string]bool{
		"vmware": true,
	}

	if validObjectType[objectType] == false {
		log.Fatalf("Error: The 'objectType' must be 'vmware'")
	}

	switch objectType {
	case "vmware":
		slaID := c.ObjectID(slaName, "sla")

		allVMinSLA := c.Get("v1", fmt.Sprintf("/vmware/vm?effective_sla_domain_id=%s&is_relic=false", slaID), httpTimeout).(map[string]interface{})

		if allVMinSLA["total"].(float64) == 0 {
			return fmt.Sprintf("The SLA '%s' is currently not protecting any %s objects.", slaName, objectType)
		}

		vmNameID := map[interface{}]interface{}{}
		for _, v := range allVMinSLA["data"].([]interface{}) {
			vmNameID[v.(map[string]interface{})["name"]] = v.(map[string]interface{})["id"]
		}

		return vmNameID

	}

	return ""
}

// PauseSnapshot suspends all snapshot activity for the provided object. The only "objectType" current supported is vmware.
//
// The function will return one of the following:
//	No change required. The '{objectName}' '{objectType}' is already paused.
//
//	The full API response for POST /internal/vmware/vm/{vmID}
func (c *Credentials) PauseSnapshot(objectName, objectType string, timeout ...int) interface{} {

	httpTimeout := httpTimeout(timeout)

	// Change the default to 180
	if httpTimeout == 15 {
		httpTimeout = 180
	}

	validObjectType := map[string]bool{
		"vmware": true,
	}

	if validObjectType[objectType] == false {
		log.Fatalf("Error: The 'objectType' must be 'vmware'")
	}

	switch objectType {
	case "vmware":
		vmID := c.ObjectID(objectName, "vmware")

		vmSummary := c.Get("v1", fmt.Sprintf("/vmware/vm/%s", vmID), httpTimeout).(map[string]interface{})

		if vmSummary["blackoutWindowStatus"].(map[string]interface{})["isSnappableBlackoutActive"].(bool) {
			return fmt.Sprintf("No change required. The '%s' '%s' is already paused.", objectName, objectType)
		}

		config := map[string]bool{}
		config["isVmPaused"] = true

		return c.Patch("v1", fmt.Sprintf("/vmware/vm/%s", vmID), config, httpTimeout)

	}

	return ""
}

// ResumeSnapshot resumes all snapshot activity for the provided object. The only "objectType" currently supported is vmware.
//
// The function will return one of the following:
//	No change required. The '{objectName}' '{objectType}' is currently not paused.
//
//	The full API response for POST /internal/vmware/vm/{vmID}
func (c *Credentials) ResumeSnapshot(objectName, objectType string, timeout ...int) interface{} {

	httpTimeout := httpTimeout(timeout)

	// Change the default to 180
	if httpTimeout == 15 {
		httpTimeout = 180
	}

	validObjectType := map[string]bool{
		"vmware": true,
	}

	if validObjectType[objectType] == false {
		log.Fatalf("Error: The 'objectType' must be 'vmware'")
	}

	switch objectType {
	case "vmware":
		vmID := c.ObjectID(objectName, "vmware")

		vmSummary := c.Get("v1", fmt.Sprintf("/vmware/vm/%s", vmID), httpTimeout).(map[string]interface{})

		if vmSummary["blackoutWindowStatus"].(map[string]interface{})["isSnappableBlackoutActive"].(bool) == false {
			return fmt.Sprintf("No change required. The '%s' '%s' is currently not paused.", objectName, objectType)
		}

		config := map[string]bool{}
		config["isVmPaused"] = false

		return c.Patch("v1", fmt.Sprintf("/vmware/vm/%s", vmID), config, httpTimeout)

	}

	return ""
}

// OnDemandSnapshotVM initiates an on-demand snapshot for the "objectName". The only "objectType" currently supported is vmware. To use the currently
// assigned SLA Domain for the snapshot use "current" for the slaName.
//
// The function will return:
//	The job status URL for the on-demand Snapshot
func (c *Credentials) OnDemandSnapshotVM(objectName, objectType, slaName string, timeout ...int) string {

	httpTimeout := httpTimeout(timeout)

	// Change the default to 180
	if httpTimeout == 15 {
		httpTimeout = 180
	}

	validObjectType := map[string]bool{
		"vmware": true,
	}

	if validObjectType[objectType] == false {
		log.Fatalf("Error: The 'objectType' must be 'vmware'")
	}

	switch objectType {
	case "vmware":
		vmID := c.ObjectID(objectName, "vmware")

		var slaID string
		switch slaName {
		case "current":
			slaID = c.Get("v1", fmt.Sprintf("/vmware/vm/%s", vmID)).(map[string]interface{})["effectiveSlaDomainId"].(string)
		default:
			slaID = c.ObjectID(slaName, "sla")
		}

		config := map[string]string{}
		config["slaId"] = slaID

		return c.Post("v1", fmt.Sprintf("/vmware/vm/%s/snapshot", vmID), config, httpTimeout).(map[string]interface{})["links"].([]interface{})[0].(map[string]interface{})["href"].(string)

	}

	return ""
}

// OnDemandSnapshotPhysical initiates an on-demand snapshot for a physical host ("hostname"). To use the currently  assigned SLA Domain for the
// snapshot use "current" for the slaName.
//
// Valid "hostOS" choices are:
//
//	Linux and Windows
//
// The function will return:
//	The job status URL for the on-demand Snapshot
func (c *Credentials) OnDemandSnapshotPhysical(hostName, slaName, fileset, hostOS string, timeout ...int) string {

	httpTimeout := httpTimeout(timeout)

	// Change the default to 180
	if httpTimeout == 15 {
		httpTimeout = 180
	}

	validHostOs := map[string]bool{
		"Linux":   true,
		"Windows": true,
	}

	if validHostOs[hostOS] == false {
		log.Fatalf("Error: The 'hostOS' must be 'Linux' or 'Windows.")
	}

	hostID := c.ObjectID(hostName, "physicalHost")

	filesetTemplateID := c.ObjectID(fileset, "filesetTemplate", hostOS)

	filesetSummary := c.Get("v1", fmt.Sprintf("/fileset?primary_cluster_id=local&host_id=%s&is_relic=false&template_id=%s", hostID, filesetTemplateID)).(map[string]interface{})

	if filesetSummary["total"] == 0 {
		log.Fatalf(fmt.Sprintf("Error: The Physical Host '%s' is not assigned to the '%s' Fileset.", hostName, fileset))
	}

	filesetID := filesetSummary["data"].([]interface{})[0].(map[string]interface{})["id"].(string)

	var slaID string
	switch slaName {
	case "current":
		slaID = filesetSummary["data"].([]interface{})[0].(map[string]interface{})["effectiveSlaDomainId"].(string)
	default:
		slaID = c.ObjectID(slaName, "sla")

	}

	config := map[string]string{}
	config["slaId"] = slaID

	return c.Post("v1", fmt.Sprintf("/fileset/%s/snapshot", filesetID), config, httpTimeout).(map[string]interface{})["links"].([]interface{})[0].(map[string]interface{})["href"].(string)
}
