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
package rubrikcdm

import (
	"errors"
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
)

// EndManagedVolumeSnapshot corresponds to POST /internal/managed_volume/{id}/end_snapshot
type EndManagedVolumeSnapshot struct {
	ID                     string   `json:"id"`
	Date                   string   `json:"date"`
	ExpirationDate         string   `json:"expirationDate"`
	SourceObjectType       string   `json:"sourceObjectType"`
	IsOnDemandSnapshot     bool     `json:"isOnDemandSnapshot"`
	CloudState             int      `json:"cloudState"`
	ConsistencyLevel       string   `json:"consistencyLevel"`
	IndexState             int      `json:"indexState"`
	ReplicationLocationIds []string `json:"replicationLocationIds"`
	ArchivalLocationIds    []string `json:"archivalLocationIds"`
	SLAID                  string   `json:"slaId"`
	SLAName                string   `json:"slaName"`
	Links                  struct {
		Self struct {
			Href string `json:"href"`
			Rel  string `json:"rel"`
		} `json:"self"`
	} `json:"links"`
}

// Cluster corresponds to /v1/cluster/{id}
type Cluster struct {
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

// ObjectID will search the Rubrik cluster for the provided "objectName" and return its ID/
//
// Valid "objectType" choices are:
//
//	vmware, sla, vmwareHost, physicalHost, filesetTemplate, managedVolume, vcenter, and ec2.
// When the "objectType" is "ec2", the objectName should correspond to the AWS Instance ID.
func (c *Credentials) ObjectID(objectName, objectType string, timeout int, hostOS ...string) (string, error) {

	validObjectType := map[string]bool{
		"vmware":          true,
		"sla":             true,
		"vmwareHost":      true,
		"physicalHost":    true,
		"filesetTemplate": true,
		"managedVolume":   true,
		"vcenter":         true,
		"ec2":             true,
		"ahv":             true,
	}

	if validObjectType[objectType] == false {
		return "", fmt.Errorf("The 'objectType' must be 'vmware', 'sla', 'vmwareHost', 'physicalHost', 'filesetTemplate', 'managedVolume', or 'vcenter'")
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
				return "", errors.New("The hostOS must be either 'Linux' or 'Windows'")

			}
		} else if len(hostOS) == 0 {
			return "", errors.New("You must provide the Fileset Tempalte OS type")
		}
		objectSummaryAPIVersion = "v1"
		objectSummaryAPIEndpoint = fmt.Sprintf("/fileset_template?primary_cluster_id=local&operating_system_type=%s&name=%s", hostOperatingSystem, objectName)
	case "managedVolume":
		objectSummaryAPIVersion = "internal"
		objectSummaryAPIEndpoint = fmt.Sprintf("/managed_volume?is_relic=false&primary_cluster_id=local&name=%s", objectName)
	case "vcenter":
		objectSummaryAPIVersion = "v1"
		objectSummaryAPIEndpoint = "/vmware/vcenter"
	case "ec2":
		objectSummaryAPIVersion = "internal"
		objectSummaryAPIEndpoint = fmt.Sprintf("/aws/ec2_instance?name=%s&is_relic=false&sort_by=instanceId&sort_order=asc", objectName)
	case "ahv":
		objectSummaryAPIVersion = "internal"
		objectSummaryAPIEndpoint = fmt.Sprintf("/nutanix/vm?primary_cluster_id=local&is_relic=false&name=%s", objectName)
	}

	apiRequest, err := c.Get(objectSummaryAPIVersion, objectSummaryAPIEndpoint, timeout)
	if err != nil {
		return "", err
	}
	if apiRequest.(map[string]interface{})["total"] == 0 {
		return "", fmt.Errorf("The %s object '%s' was not found on the Rubrik cluster", objectType, objectName)
	} else if apiRequest.(map[string]interface{})["total"].(float64) > 0 {
		objectIDs := make([]string, 0)
		// # Define the "object name" to search for
		var nameValue string
		if objectType == "physicalHost" {
			nameValue = "hostname"
		} else if objectType == "ec2" {
			nameValue = "instanceId"
		} else {
			nameValue = "name"
		}

		for _, v := range apiRequest.(map[string]interface{})["data"].([]interface{}) {
			if v.(interface{}).(map[string]interface{})[nameValue].(string) == objectName {
				objectIDs = append(objectIDs, v.(interface{}).(map[string]interface{})["id"].(string))
			}
		}

		if len(objectIDs) > 1 {
			return "", fmt.Errorf("Multiple %s objects named '%s' were found on the Rubrik cluster. Unable to return a specific object id", objectType, objectName)
		} else if len(objectIDs) == 0 {
			return "", fmt.Errorf("The %s object '%s' was not found on the Rubrik cluster", objectType, objectName)
		} else {
			return objectIDs[0], nil
		}
	}

	return "", fmt.Errorf("The %s object '%s' was not found on the Rubrik cluster", objectType, objectName)

}

// AssignSLA adds the "objectName" to the "slaName". vmware and ahv are the only supported "objectType". To exclude the object from all SLA assignments
// use "do not protect" as the "slaName". To assign the selected object to the SLA of the next higher level object, use "clear" as the "slaName".
//
// The function will return one of the following:
//	No change required. The vSphere VM '{objectName}' is already assigned to the '{slaName}' SLA Domain.
//
//	The full API response for POST /internal/sla_domain/{slaID}/assign.
func (c *Credentials) AssignSLA(objectName, objectType, slaName string, timeout ...int) (*StatusCode, error) {

	httpTimeout := httpTimeout(timeout)

	validObjectType := map[string]bool{
		"vmware": true,
		"ahv":    true,
	}

	if validObjectType[objectType] == false {
		return nil, fmt.Errorf("The 'objectType' must be 'vmware' or 'ahv'.")
	}

	var slaID string
	var err error
	switch slaName {
	case "do not protect":
		slaID = "UNPROTECTED"
	case "clear":
		slaID = "INHERIT"
	default:
		slaID, err = c.ObjectID(slaName, "sla", httpTimeout)
		if err != nil {
			return nil, err
		}
	}

	config := map[string]interface{}{}
	switch objectType {
	case "vmware":
		vmID, err := c.ObjectID(objectName, "vmware", httpTimeout)
		if err != nil {
			return nil, err
		}

		vmSummary, err := c.Get("v1", fmt.Sprintf("/vmware/vm/%s", vmID), httpTimeout)
		if err != nil {
			return nil, err
		}

		var currentSLAID string
		switch slaID {
		case "INHERIT":
			currentSLAID = vmSummary.(map[string]interface{})["configuredSlaDomainId"].(string)
		default:
			currentSLAID = vmSummary.(map[string]interface{})["effectiveSlaDomainId"].(string)
		}

		if slaID == currentSLAID {
			return nil, fmt.Errorf("No change required. The vSphere VM '%s' is already assigned to the '%s' SLA Domain", objectName, slaName)
		}

		config["managedIds"] = []string{vmID}
	case "ahv":
		vmID, err := c.ObjectID(objectName, "ahv", httpTimeout)
		if err != nil {
			return nil, err
		}

		vmSummary, err := c.Get("internal", fmt.Sprintf("/nutanix/vm/%s", vmID), httpTimeout)
		if err != nil {
			return nil, err
		}

		var currentSLAID string
		switch slaID {
		case "INHERIT":
			currentSLAID = vmSummary.(map[string]interface{})["configuredSlaDomainId"].(string)
		default:
			currentSLAID = vmSummary.(map[string]interface{})["effectiveSlaDomainId"].(string)
		}

		if slaID == currentSLAID {
			return nil, fmt.Errorf("No change required. The AHV VM '%s' is already assigned to the '%s' SLA Domain", objectName, slaName)
		}

		config["managedIds"] = []string{vmID}
	}

	apiRequest, err := c.Post("internal", fmt.Sprintf("/sla_domain/%s/assign", slaID), config, httpTimeout)
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

// BeginManagedVolumeSnapshot opens a managed volume for writes. All writes to the managed volume until the snapshot is
// ended will be part of its snapshot.
//
// The function will return one of the following:
//	No change required. The Managed Volume '{name}' is already in a writeable state.
//
//	The full API response for POST /internal/managed_volume/{managedVolumeID}/begin_snapshot
func (c *Credentials) BeginManagedVolumeSnapshot(name string, timeout ...int) (*StatusCode, error) {

	httpTimeout := httpTimeout(timeout)

	managedVolumeID, err := c.ObjectID(name, "managedVolume", httpTimeout)
	if err != nil {
		return nil, err
	}

	managedVolumeSummary, err := c.Get("internal", fmt.Sprintf("/managed_volume/%s", managedVolumeID), httpTimeout)
	if err != nil {
		return nil, err
	}

	if managedVolumeSummary.(map[string]interface{})["isWritable"].(bool) {
		return nil, fmt.Errorf("No change required. The Managed Volume '%s' is already in a writeable state", name)
	}

	config := map[string]string{}

	apiRequest, err := c.Post("internal", fmt.Sprintf("/managed_volume/%s/begin_snapshot", managedVolumeID), config, httpTimeout)
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

// EndManagedVolumeSnapshot closes a managed volume for writes. A snapshot will be created containing all writes since the last begin snapshot call.
//
// The function will return one of the following:
//	No change required. The Managed Volume '{name}' is already in a read-only state.
//
//	The full API response for POST /internal/managed_volume/{managedVolumeID}/end_snapshot
func (c *Credentials) EndManagedVolumeSnapshot(name, slaName string, timeout ...int) (*EndManagedVolumeSnapshot, error) {

	httpTimeout := httpTimeout(timeout)

	managedVolumeID, err := c.ObjectID(name, "managedVolume", httpTimeout)
	if err != nil {
		return nil, err
	}

	managedVolumeSummary, err := c.Get("internal", fmt.Sprintf("/managed_volume/%s", managedVolumeID), httpTimeout)
	if err != nil {
		return nil, err
	}

	if managedVolumeSummary.(map[string]interface{})["isWritable"].(bool) == false {
		return nil, fmt.Errorf("No change required. The Managed Volume '%s' is already in a read-only state", name)
	}

	var slaID string
	config := map[string]interface{}{}
	switch slaName {
	case "current":
	default:
		slaID, err = c.ObjectID(slaName, "sla", httpTimeout)
		if err != nil {
			return nil, err
		}
		config["retentionConfig"] = map[string]interface{}{}
		config["retentionConfig"].(map[string]interface{})["slaId"] = slaID
	}

	apiRequest, err := c.Post("internal", fmt.Sprintf("/managed_volume/%s/end_snapshot", managedVolumeID), config, httpTimeout)
	if err != nil {
		return nil, err
	}

	// Convert the API Response (map[string]interface{}) to a struct
	var endSnapshot EndManagedVolumeSnapshot
	mapErr := mapstructure.Decode(apiRequest, &endSnapshot)
	if mapErr != nil {
		return nil, mapErr
	}

	return &endSnapshot, nil

}

// GetSLAObjects returns the name and ID of a specific object type.
func (c *Credentials) GetSLAObjects(slaName, objectType string, timeout ...int) (interface{}, error) {

	httpTimeout := httpTimeout(timeout)

	validObjectType := map[string]bool{
		"vmware": true,
	}

	if validObjectType[objectType] == false {
		return nil, fmt.Errorf("The 'objectType' must be 'vmware'")
	}

	switch objectType {
	case "vmware":
		slaID, err := c.ObjectID(slaName, "sla", httpTimeout)
		if err != nil {
			return nil, err
		}

		allVMinSLA, err := c.Get("v1", fmt.Sprintf("/vmware/vm?effective_sla_domain_id=%s&is_relic=false", slaID), httpTimeout)
		if err != nil {
			return nil, err
		}

		if allVMinSLA.(map[string]interface{})["total"].(float64) == 0 {
			return fmt.Sprintf("The SLA '%s' is currently not protecting any %s objects.", slaName, objectType), nil
		}

		vmNameID := map[interface{}]interface{}{}
		for _, v := range allVMinSLA.(map[string]interface{})["data"].([]interface{}) {
			vmNameID[v.(map[string]interface{})["name"]] = v.(map[string]interface{})["id"]
		}

		return vmNameID, nil

	}

	return "", nil
}

// PauseSnapshot suspends all snapshot activity for the provided object. The only "objectType" current supported is vmware.
//
// The function will return one of the following:
//	No change required. The '{objectName}' '{objectType}' is already paused.
//
//	The full API response for POST /internal/vmware/vm/{vmID}
func (c *Credentials) PauseSnapshot(objectName, objectType string, timeout ...int) (interface{}, error) {

	httpTimeout := httpTimeout(timeout)

	// Change the default to 180
	if httpTimeout == 15 {
		httpTimeout = 180
	}

	validObjectType := map[string]bool{
		"vmware": true,
	}

	if validObjectType[objectType] == false {
		return nil, fmt.Errorf("The 'objectType' must be 'vmware'")
	}

	switch objectType {
	case "vmware":
		vmID, err := c.ObjectID(objectName, "vmware", httpTimeout)
		if err != nil {
			return nil, err
		}

		vmSummary, err := c.Get("v1", fmt.Sprintf("/vmware/vm/%s", vmID), httpTimeout)
		if err != nil {
			return nil, err
		}

		if vmSummary.(map[string]interface{})["blackoutWindowStatus"].(map[string]interface{})["isSnappableBlackoutActive"].(bool) {
			return fmt.Sprintf("No change required. The '%s' '%s' is already paused.", objectName, objectType), nil
		}

		config := map[string]bool{}
		config["isVmPaused"] = true

		apiRequest, err := c.Patch("v1", fmt.Sprintf("/vmware/vm/%s", vmID), config, httpTimeout)
		if err != nil {
			return nil, err
		}

		return apiRequest, nil

	}

	return "", nil
}

// ResumeSnapshot resumes all snapshot activity for the provided object. The only "objectType" currently supported is vmware.
//
// The function will return one of the following:
//	No change required. The '{objectName}' '{objectType}' is currently not paused.
//
//	The full API response for POST /internal/vmware/vm/{vmID}
func (c *Credentials) ResumeSnapshot(objectName, objectType string, timeout ...int) (interface{}, error) {

	httpTimeout := httpTimeout(timeout)

	// Change the default to 180
	if httpTimeout == 15 {
		httpTimeout = 180
	}

	validObjectType := map[string]bool{
		"vmware": true,
	}

	if validObjectType[objectType] == false {
		return nil, fmt.Errorf("The 'objectType' must be 'vmware'")
	}

	switch objectType {
	case "vmware":
		vmID, err := c.ObjectID(objectName, "vmware", httpTimeout)
		if err != nil {
			return nil, err
		}

		vmSummary, err := c.Get("v1", fmt.Sprintf("/vmware/vm/%s", vmID), httpTimeout)
		if err != nil {
			return nil, err
		}

		if vmSummary.(map[string]interface{})["blackoutWindowStatus"].(map[string]interface{})["isSnappableBlackoutActive"].(bool) == false {
			return fmt.Sprintf("No change required. The '%s' '%s' is currently not paused.", objectName, objectType), nil
		}

		config := map[string]bool{}
		config["isVmPaused"] = false

		apiRequest, err := c.Patch("v1", fmt.Sprintf("/vmware/vm/%s", vmID), config, httpTimeout)
		if err != nil {
			return nil, err
		}

		return apiRequest, nil

	}

	return "", nil
}

// OnDemandSnapshotVM initiates an on-demand snapshot for the "objectName". The only "objectType" currently supported is vmware. To use the currently
// assigned SLA Domain for the snapshot use "current" for the slaName.
//
// The function will return:
//	The job status URL for the on-demand Snapshot
func (c *Credentials) OnDemandSnapshotVM(objectName, objectType, slaName string, timeout ...int) (string, error) {

	httpTimeout := httpTimeout(timeout)

	// Change the default to 180
	if httpTimeout == 15 {
		httpTimeout = 180
	}

	validObjectType := map[string]bool{
		"vmware": true,
	}

	if validObjectType[objectType] == false {
		return "", fmt.Errorf("The 'objectType' must be 'vmware'")
	}

	switch objectType {
	case "vmware":
		vmID, err := c.ObjectID(objectName, "vmware", httpTimeout)
		if err != nil {
			return "", err
		}

		var slaID interface{}
		switch slaName {
		case "current":
			slaID, err = c.Get("v1", fmt.Sprintf("/vmware/vm/%s", vmID))
			if err != nil {
				return "", err
			}
		default:
			slaID, err = c.ObjectID(slaName, "sla", httpTimeout)
			if err != nil {
				return "", err
			}
		}

		config := map[string]string{}
		config["slaId"] = slaID.(map[string]interface{})["effectiveSlaDomainId"].(string)

		apiRequest, err := c.Post("v1", fmt.Sprintf("/vmware/vm/%s/snapshot", vmID), config, httpTimeout)
		if err != nil {
			return "", err
		}

		return apiRequest.(map[string]interface{})["links"].([]interface{})[0].(map[string]interface{})["href"].(string), nil

	}

	return "", nil
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
func (c *Credentials) OnDemandSnapshotPhysical(hostName, slaName, fileset, hostOS string, timeout ...int) (string, error) {

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
		return "", fmt.Errorf("The 'hostOS' must be 'Linux' or 'Windows")
	}

	hostID, err := c.ObjectID(hostName, "physicalHost", httpTimeout)
	if err != nil {
		return "", err
	}

	filesetTemplateID, err := c.ObjectID(fileset, "filesetTemplate", httpTimeout, hostOS)
	if err != nil {
		return "", err
	}

	filesetSummary, err := c.Get("v1", fmt.Sprintf("/fileset?primary_cluster_id=local&host_id=%s&is_relic=false&template_id=%s", hostID, filesetTemplateID))
	if err != nil {
		return "", err
	}

	if filesetSummary.(map[string]interface{})["total"] == 0 {
		return "", fmt.Errorf("The Physical Host '%s' is not assigned to the '%s' Fileset", hostName, fileset)
	}

	filesetID := filesetSummary.(map[string]interface{})["data"].([]interface{})[0].(map[string]interface{})["id"].(string)

	var slaID string
	switch slaName {
	case "current":
		slaID = filesetSummary.(map[string]interface{})["data"].([]interface{})[0].(map[string]interface{})["effectiveSlaDomainId"].(string)
	default:
		slaID, err = c.ObjectID(slaName, "sla", httpTimeout)
		if err != nil {
			return "", err
		}

	}

	config := map[string]string{}
	config["slaId"] = slaID

	apiRequeset, err := c.Post("v1", fmt.Sprintf("/fileset/%s/snapshot", filesetID), config, httpTimeout)
	if err != nil {
		return "", err
	}

	return apiRequeset.(map[string]interface{})["links"].([]interface{})[0].(map[string]interface{})["href"].(string), nil
}

func (c *Credentials) DateTimeConversion(dateTime string, timeout ...int) (string, error) {

	httpTimeout := httpTimeout(timeout)

	apiRequest, err := c.Get("v1", "/cluster/me", httpTimeout)
	if err != nil {
		return "", err
	}

	// Convert the API Response (map[string]interface{}) to a struct
	var cluster Cluster
	mapErr := mapstructure.Decode(apiRequest, &cluster)
	if mapErr != nil {
		return "", mapErr
	}

	currentTimezone, _ := time.LoadLocation(cluster.Timezone.Timezone)
	// Month-Day-Year

	snapshotDateTime, err := time.ParseInLocation("01-02-2006 3:04 PM", dateTime, currentTimezone)
	if err != nil {
		return "", fmt.Errorf("The provided 'dateTime' does not match the required format (Month-Day-Year Hour:Minute AM/PM). Ex. 01-02-2006 3:04 PM")
	}

	return snapshotDateTime.UTC().Format(time.RFC3339), nil

}

// RecoverFileDownload initiates to create a file download job from a fileset backup.
//
// Valid "hostOS" choices are:
//
//	Linux and Windows
//
// The function will return:
//	The job status URL for file download job from a fileset backup
func (c *Credentials) RecoverFileDownload(hostName, fileset, hostOS, filePath, dateTime string, timeout ...int) (string, error) {
	httpTimeout := httpTimeout(timeout)

	validHostOs := map[string]bool{
		"Linux":   true,
		"Windows": true,
	}

	if validHostOs[hostOS] == false {
		return "", fmt.Errorf("The 'hostOS' must be 'Linux' or 'Windows")
	}

	hostID, err := c.ObjectID(hostName, "physicalHost", httpTimeout)
	if err != nil {
		return "", err
	}

	filesetTemplateID, err := c.ObjectID(fileset, "filesetTemplate", httpTimeout, hostOS)
	if err != nil {
		return "", err
	}

	filesetSummary, err := c.Get("v1", fmt.Sprintf("/fileset?primary_cluster_id=local&host_id=%s&is_relic=false&template_id=%s", hostID, filesetTemplateID))
	if err != nil {
		return "", err
	}

	if filesetSummary.(map[string]interface{})["total"] == 0 {
		return "", fmt.Errorf("The Physical Host '%s' is not assigned to the '%s' Fileset", hostName, fileset)
	}

	filesetID := filesetSummary.(map[string]interface{})["data"].([]interface{})[0].(map[string]interface{})["id"].(string)

	filesetDetail, err := c.Get("v1", fmt.Sprintf("/fileset/%s", filesetID))
	if err != nil {
		return "", err
	}
	snapshotSummary := filesetDetail.(map[string]interface{})["snapshots"].([]interface{})

	if len(snapshotSummary) == 0 {
		return "", fmt.Errorf("The Physical Host '%s' does not have any snapshot by '%s' Fileset", hostName, fileset)
	}

	snapshotDateTimeStr, err := c.DateTimeConversion(dateTime)
	if err != nil {
		return "", err
	}

	snapshotDateTime, _ := time.Parse(time.RFC3339, snapshotDateTimeStr)
	if err != nil {
		return "", err
	}

	var snapshotID string
	for _, v := range snapshotSummary {
		date, _ := time.Parse(time.RFC3339, v.(map[string]interface{})["date"].(string))
		if err != nil {
			return "", err
		}
		diff := date.Sub(snapshotDateTime)
		if 0 <= diff && diff < time.Duration(60)*time.Second {
			snapshotID = v.(map[string]interface{})["id"].(string)
			break
		}
	}
	if snapshotID == "" {
		return "", fmt.Errorf("The Physical Host '%s' does not have any snapshot at '%s'", hostName, dateTime)
	}
	config := map[string]string{
		"sourceDir": filePath,
	}
	apiRequest, err := c.Post("v1", fmt.Sprintf("/fileset/snapshot/%s/download_file", snapshotID), config)

	if err != nil {
		return "", err
	}

	return apiRequest.(map[string]interface{})["links"].([]interface{})[0].(map[string]interface{})["href"].(string), nil
}
