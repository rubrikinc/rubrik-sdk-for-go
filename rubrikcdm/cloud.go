// // Copyright 2018 Rubrik, Inc.
// //
// // Licensed under the Apache License, Version 2.0 (the "License");
// // you may not use this file except in compliance with the License.
// // You may obtain a copy of the License prop
// //  http://www.apache.org/licenses/LICENSE-2.0
// //
// // Unless required by applicable law or agreed to in writing, software
// // distributed under the License is distributed on an "AS IS" BASIS,
// // WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// // See the License for the specific language governing permissions and
// // limitations under the License.

package rubrikcdm

import (
	"fmt"
	"reflect"
	"strings"
)

// AddAWSNativeAccount enables the management and protection of Amazon Elastic Compute Cloud (Amazon EC2) instances. The "regionalBoltNetworkConfigs"
// should be a list of dictionaries in the following format:
//
// 	usEast1 := map[string]string{}
//	usEast1["region"] = "us-east-1"
//	usEast1["region"] = "us-east-1"
//	usEast1["region"] = "us-east-1"
//	usEast1["vNetId"] = "vpc-11a44968"
//	usEast1["subnetId"] = "subnet-3ac58e06"
//	usEast1["securityGroupId"] = "sg-9ba90ee5"
//
//	regionalBoltNetworkConfigs := []interface{}{usEast1}
//
// Valid "awsRegion" choices are:
//
//	ap-south-1,ap-northeast-3, ap-northeast-2, ap-southeast-1, ap-southeast-2, ap-northeast-1, ca-central-1, cn-north-1, cn-northwest-1, eu-central-1, eu-west-1,
//	eu-west-2, eu-west-3, us-west-1, us-east-1, us-east-2, and us-west-2.
//
// The function will return one of the following:
//	No change required. Cloud native source with access key '{awsAccessKey}' is already configured on the Rubrik cluster.
// //
//	The full API response for POST /internal/aws/account.
func (c *Credentials) AddAWSNativeAccount(awsAccountName, awsAccessKey, awsSecretKey string, awsRegions []string, regionalBoltNetworkConfigs interface{}, timeout ...int) (interface{}, error) {

	minimumClusterVersion := c.ClusterVersionCheck(4.2)
	if minimumClusterVersion != nil {
		return nil, minimumClusterVersion
	}

	httpTimeout := httpTimeout(timeout)

	validAWSRegions := map[string]bool{
		"ap-south-1":     true,
		"ap-northeast-3": true,
		"ap-northeast-2": true,
		"ap-southeast-1": true,
		"ap-southeast-2": true,
		"ap-northeast-1": true,
		"ca-central-1":   true,
		"cn-north-1":     true,
		"cn-northwest-1": true,
		"eu-central-1":   true,
		"eu-west-1":      true,
		"eu-west-2":      true,
		"eu-west-3":      true,
		"us-west-1":      true,
		"us-east-1":      true,
		"us-east-2":      true,
		"us-west-2":      true,
	}

	for _, region := range awsRegions {
		if validAWSRegions[region] == false {
			return nil, fmt.Errorf("'%s' is not a valid AWS Region", region)
		}

	}

	cloudNativeOnCluster, err := c.Get("internal", "/aws/account", httpTimeout)
	if err != nil {
		return nil, err
	}

	if len(cloudNativeOnCluster.(map[string]interface{})["data"].([]interface{})) > 0 {
		var currentAWSAccountName string
		var currentAWSConfigID string
		for _, v := range cloudNativeOnCluster.(map[string]interface{})["data"].([]interface{}) {
			currentAWSAccountName = (v.(interface{}).(map[string]interface{})["name"].(string))
			currentAWSConfigID = (v.(interface{}).(map[string]interface{})["id"].(string))

		}

		// TODO - Add additional logic checks for the region and bolt configs
		currentAccessKey, err := c.Get("internal", fmt.Sprintf("/aws/account/%s", currentAWSConfigID), httpTimeout)
		if err != nil {
			return nil, err
		}

		if currentAccessKey.(interface{}).(map[string]interface{})["accessKey"].(string) == awsAccessKey {
			return fmt.Sprintf("No change required. Cloud native source with access key '%s' is already configured on the Rubrik cluster.", awsAccessKey), nil
		}

		if currentAWSAccountName == awsAccountName {
			return nil, fmt.Errorf("A Cloud native source with name '%s' already exists. Please enter a unique 'awsAccountName'", awsAccountName)
		}

	}

	config := map[string]interface{}{}
	config["name"] = awsAccountName
	config["accessKey"] = awsAccessKey
	config["secretKey"] = awsSecretKey
	config["regions"] = awsRegions

	config["regionalBoltNetworkConfigs"] = regionalBoltNetworkConfigs

	apiRequest, err := c.Post("internal", "/aws/account", config, httpTimeout)
	if err != nil {
		return nil, err
	}

	return apiRequest, nil

}

// AWSS3CloudOutRSA configures a new AWS S3 archive target using a RSA Key for encryption.
//
// Valid "awsRegion" choices are:
//
//	ap-south-1,ap-northeast-3, ap-northeast-2, ap-southeast-1, ap-southeast-2, ap-northeast-1, ca-central-1, cn-north-1, cn-northwest-1, eu-central-1, eu-west-1,
//	eu-west-2, eu-west-3, us-west-1, us-east-1, us-east-2, and us-west-2.
//
// Valid "storageClass" choices are:
//
//	ap-standard-1, standard_ia, and reduced_redundancy
//
// The function will return one of the following:
//	- No change required. The '{archiveName}' archive location is already configured on the Rubrik cluster.
//
//	- The full API response for POST /internal/archive/object_store.
func (c *Credentials) AWSS3CloudOutRSA(awsBucketName, storageClass, archiveName, awsRegion, awsAccessKey, awsSecretKey, rsaKey string, timeout ...int) (interface{}, error) {

	httpTimeout := httpTimeout(timeout)

	validAWSRegions := map[string]bool{
		"ap-south-1":     true,
		"ap-northeast-3": true,
		"ap-northeast-2": true,
		"ap-southeast-1": true,
		"ap-southeast-2": true,
		"ap-northeast-1": true,
		"ca-central-1":   true,
		"cn-north-1":     true,
		"cn-northwest-1": true,
		"eu-central-1":   true,
		"eu-west-1":      true,
		"eu-west-2":      true,
		"eu-west-3":      true,
		"us-west-1":      true,
		"us-east-1":      true,
		"us-east-2":      true,
		"us-west-2":      true,
	}

	validStorageClass := map[string]bool{
		"standard":           true,
		"standard_ia":        true,
		"reduced_redundancy": true,
	}

	if validAWSRegions[awsRegion] == false {
		return nil, fmt.Errorf("%s is not a valid AWS Region", awsRegion)
	}

	if validStorageClass[storageClass] == false {
		return nil, fmt.Errorf("%s is not a valid 'storageClass'. Please use 'standard', 'standard_ia', or 'reduced_redundancy'", storageClass)
	}

	config := map[string]string{}
	config["name"] = archiveName
	config["bucket"] = strings.ToLower(awsBucketName)
	config["defaultRegion"] = awsRegion
	config["storageClass"] = strings.ToUpper(storageClass)
	config["accessKey"] = awsAccessKey
	config["secretKey"] = awsSecretKey
	config["objectStoreType"] = "S3"
	config["pemFileContent"] = rsaKey

	// Create a simplified config that only includes the values returned by Rubrik that can be used for idempotence check
	redactedConfig := map[string]interface{}{}
	redactedConfig["name"] = archiveName
	redactedConfig["bucket"] = strings.ToLower(awsBucketName)
	redactedConfig["defaultRegion"] = awsRegion
	redactedConfig["storageClass"] = strings.ToUpper(storageClass)
	redactedConfig["accessKey"] = awsAccessKey
	redactedConfig["objectStoreType"] = "S3"

	archivesOnCluster, err := c.Get("internal", "/archive/object_store", httpTimeout)
	if err != nil {
		return nil, err
	}

	for _, v := range archivesOnCluster.(map[string]interface{})["data"].([]interface{}) {
		archiveDefinition := (v.(interface{}).(map[string]interface{})["definition"])
		delete(archiveDefinition.(map[string]interface{}), "definition")

		archivePresent := reflect.DeepEqual(redactedConfig, archiveDefinition)

		if archivePresent {
			return fmt.Sprintf("No change required. The '%s' archive location is already configured on the Rubrik cluster.", archiveName), nil
		}

		if archiveDefinition.(map[string]interface{})["objectStoreType"] == "S3" && archiveDefinition.(map[string]interface{})["name"] == archiveName {

			return nil, fmt.Errorf("An archive location with the name '%s' already exists. Please enter a unique 'archiveName'", archiveName)
		}

	}

	apiRequest, err := c.Post("internal", "/archive/object_store", config, httpTimeout)
	if err != nil {
		return nil, err
	}

	return apiRequest, nil
}

// AWSS3CloudOutKMS configures a new AWS S3 archive target using a AWS KMS Master Key ID for encryption.
//
// Valid "awsRegion" choices are:
//
//	ap-south-1,ap-northeast-3, ap-northeast-2, ap-southeast-1, ap-southeast-2, ap-northeast-1, ca-central-1, cn-north-1, cn-northwest-1, eu-central-1, eu-west-1,
//	eu-west-2, eu-west-3, us-west-1, us-east-1, us-east-2, and us-west-2.
//
// Valid "storageClass" choices are:
//
//	ap-standard-1, standard_ia, and reduced_redundancy
//
// The function will return one of the following:
//	- No change required. The '{archiveName}' archive location is already configured on the Rubrik cluster.
//
//	- The full API response for POST /internal/archive/object_store/{archiveID}
func (c *Credentials) AWSS3CloudOutKMS(awsBucketName, storageClass, archiveName, awsRegion, awsAccessKey, awsSecretKey, kmsMasterKeyID string, timeout ...int) (interface{}, error) {

	httpTimeout := httpTimeout(timeout)

	validAWSRegions := map[string]bool{
		"ap-south-1":     true,
		"ap-northeast-3": true,
		"ap-northeast-2": true,
		"ap-southeast-1": true,
		"ap-southeast-2": true,
		"ap-northeast-1": true,
		"ca-central-1":   true,
		"cn-north-1":     true,
		"cn-northwest-1": true,
		"eu-central-1":   true,
		"eu-west-1":      true,
		"eu-west-2":      true,
		"eu-west-3":      true,
		"us-west-1":      true,
		"us-east-1":      true,
		"us-east-2":      true,
		"us-west-2":      true,
	}

	validStorageClass := map[string]bool{
		"standard":           true,
		"standard_ia":        true,
		"reduced_redundancy": true,
	}

	if validAWSRegions[awsRegion] == false {
		return nil, fmt.Errorf("%s is not a valid AWS Region", awsRegion)
	}

	if validStorageClass[storageClass] == false {
		return nil, fmt.Errorf("%s is not a valid 'storageClass'. Please use 'standard', 'standard_ia', or 'reduced_redundancy'", storageClass)
	}

	config := map[string]string{}
	config["name"] = archiveName
	config["bucket"] = strings.ToLower(awsBucketName)
	config["defaultRegion"] = awsRegion
	config["storageClass"] = strings.ToUpper(storageClass)
	config["accessKey"] = awsAccessKey
	config["secretKey"] = awsSecretKey
	config["objectStoreType"] = "S3"
	config["kmsMasterKeyId"] = kmsMasterKeyID

	// Create a simplified config that only includes the values returned by Rubrik that can be used for idempotence check
	redactedConfig := map[string]interface{}{}
	redactedConfig["name"] = archiveName
	redactedConfig["bucket"] = strings.ToLower(awsBucketName)
	redactedConfig["defaultRegion"] = awsRegion
	redactedConfig["storageClass"] = strings.ToUpper(storageClass)
	redactedConfig["accessKey"] = awsAccessKey
	redactedConfig["objectStoreType"] = "S3"

	archivesOnCluster, err := c.Get("internal", "/archive/object_store", httpTimeout)
	if err != nil {
		return nil, err
	}

	for _, v := range archivesOnCluster.(map[string]interface{})["data"].([]interface{}) {
		archiveDefinition := (v.(interface{}).(map[string]interface{})["definition"])
		delete(archiveDefinition.(map[string]interface{}), "definition")

		archivePresent := reflect.DeepEqual(redactedConfig, archiveDefinition)

		if archivePresent {
			return fmt.Sprintf("No change required. The '%s' archive location is already configured on the Rubrik cluster.", archiveName), nil
		}

		if archiveDefinition.(map[string]interface{})["objectStoreType"] == "S3" && archiveDefinition.(map[string]interface{})["name"] == archiveName {
			return nil, fmt.Errorf("An archive location with the name '%s' already exists. Please enter a unique 'archiveName'", archiveName)
		}

	}

	apiRequest, err := c.Post("internal", "/archive/object_store", config, httpTimeout)
	if err != nil {
		return nil, err
	}

	return apiRequest, nil
}

// AWSS3CloudOn provides the ability to convert a vSphere virtual machines snapshot, an archived snapshot, or a replica into an Amazon Machine Image (AMI)
// and then launch that AMI into an Elastic Compute Cloud (EC2) instance on an Amazon Virtual Private Cloud (VPC).
//
// The function will return one of the following:
//	- No change required. The '{archiveName}' archive location is already configured for CloudOn.
//
//	- The full API response for PATCH /internal/archive/object_store.
func (c *Credentials) AWSS3CloudOn(archiveName, vpcID, subnetID, securityGroupID string, timeout ...int) (interface{}, error) {

	httpTimeout := httpTimeout(timeout)

	config := map[string]interface{}{}
	config["defaultComputeNetworkConfig"] = map[string]interface{}{}
	config["defaultComputeNetworkConfig"].(map[string]interface{})["subnetId"] = subnetID
	config["defaultComputeNetworkConfig"].(map[string]interface{})["vNetId"] = vpcID
	config["defaultComputeNetworkConfig"].(map[string]interface{})["securityGroupId"] = securityGroupID

	archivesOnCluster, err := c.Get("internal", "/archive/object_store", httpTimeout)
	if err != nil {
		return nil, err
	}

	for _, v := range archivesOnCluster.(map[string]interface{})["data"].([]interface{}) {
		archiveDefinition := (v.(interface{}).(map[string]interface{})["definition"])

		if archiveDefinition.(map[string]interface{})["objectStoreType"] == "S3" && archiveDefinition.(map[string]interface{})["name"] == archiveName {

			archivePresent := reflect.DeepEqual(archiveDefinition.(map[string]interface{})["defaultComputeNetworkConfig"], config["defaultComputeNetworkConfig"])
			if archivePresent {
				return fmt.Sprintf("No change required. The '%s' archive location is already configured for CloudOn.", archiveName), nil

			}

			archiveID := (v.(interface{}).(map[string]interface{})["id"])
			apiRequest, err := c.Patch("internal", fmt.Sprintf("/archive/object_store/%s", archiveID), config, httpTimeout)
			if err != nil {
				return nil, err
			}

			return apiRequest, err

		}

	}
	return nil, fmt.Errorf("The Rubrik cluster does not have an archive location named '%s'", archiveName)

}

// AzureCloudOut configures a new Azure archive target.
//
// Valid "instanceType" choices are:
//
//	default, china, germany, and government
//
// The function will return one of the following:
//	- No change required. The '{archiveName}' archive location is already configured on the Rubrik cluster.
//
//	- The full API response for POST /internal/archive/object_store.
func (c *Credentials) AzureCloudOut(container, azureAccessKey, storageAccountName, archiveName, instanceType, rsaKey string, timeout ...int) (interface{}, error) {

	httpTimeout := httpTimeout(timeout)

	validInstanceTypes := map[string]bool{
		"default":    true,
		"china":      true,
		"germany":    true,
		"government": true,
	}

	if validInstanceTypes[instanceType] == false {
		return nil, fmt.Errorf("'%s' is not a valid Azure Instance Type. Valid choices are 'default', 'china', 'germany', or 'government'", instanceType)
	}

	config := map[string]string{}
	config["name"] = archiveName
	config["bucket"] = container
	config["accessKey"] = storageAccountName
	config["secretKey"] = azureAccessKey
	config["objectStoreType"] = "Azure"
	config["pemFileContent"] = rsaKey
	switch instanceType {
	case "government":
		config["endpoint"] = "core.usgovcloudapi.net"
	case "germany":
		config["endpoint"] = "core.cloudapi.de"
	case "china":
		config["endpoint"] = "core.chinacloudapi.cn"
	}

	// Create a simplified config that only includes the values returned by Rubrik that can be used for idempotence check
	redactedConfig := map[string]interface{}{}
	redactedConfig["objectStoreType"] = "Azure"
	redactedConfig["name"] = archiveName
	redactedConfig["accessKey"] = storageAccountName
	redactedConfig["bucket"] = container
	switch instanceType {
	case "government":
		redactedConfig["endpoint"] = "core.usgovcloudapi.net"
	case "germany":
		redactedConfig["endpoint"] = "core.cloudapi.de"
	case "china":
		redactedConfig["endpoint"] = "core.chinacloudapi.cn"
	}

	archivesOnCluster, err := c.Get("internal", "/archive/object_store", httpTimeout)
	if err != nil {
		return nil, err
	}

	for _, v := range archivesOnCluster.(map[string]interface{})["data"].([]interface{}) {
		archiveDefinition := (v.(interface{}).(map[string]interface{})["definition"])
		delete(archiveDefinition.(map[string]interface{}), "defaultComputeNetworkConfig")
		delete(archiveDefinition.(map[string]interface{}), "isComputeEnabled")
		delete(archiveDefinition.(map[string]interface{}), "isConsolidationEnabled")
		delete(archiveDefinition.(map[string]interface{}), "azureComputeSummary")

		archivePresent := reflect.DeepEqual(redactedConfig, archiveDefinition)

		if archivePresent {
			return fmt.Sprintf("No change required. The '%s' archive location is already configured on the Rubrik cluster.", archiveName), nil
		}

		if archiveDefinition.(map[string]interface{})["objectStoreType"] == "Azure" && archiveDefinition.(map[string]interface{})["name"] == archiveName {
			return nil, fmt.Errorf("An arhive location with the name '%s' already exists. Please enter a unique 'archiveName'", archiveName)
		}

	}

	apiRequest, err := c.Post("internal", "/archive/object_store", config, httpTimeout)
	if err != nil {
		return nil, err
	}

	return apiRequest, nil

}

// AzureCloudOn provides the ability to convert a snapshot, archived snapshot, or replica into a Virtual Hard Disk (VHD). This enables the instantiation
// of the associated virtual machine on the Microsoft Azure cloud platform.
//
// Valid "region" choices are:
// 	westus, westus2, centralus, eastus, eastus2, northcentralus, southcentralus, westcentralus, canadacentral, canadaeast, brazilsouth,
//	northeurope, westeurope, uksouth, ukwest, eastasia, southeastasia, japaneast, japanwest, australiaeast, australiasoutheast, centralindia,
//	southindia, westindia, koreacentral, koreasouth
//
// The function will return one of the following:
//	- No change required. The '{archiveName}' archive location is already configured for CloudOn.
//
//	- The full API response for PATCH /internal/archive/object_store.
func (c *Credentials) AzureCloudOn(archiveName, container, storageAccountName, applicationID, applicationKey, directoryID, region, virtualNetworkID, subnetName, securityGroupID string, timeout ...int) (interface{}, error) {

	httpTimeout := httpTimeout(timeout)

	validRegions := map[string]bool{
		"westus":             true,
		"westus2":            true,
		"centralus":          true,
		"eastus":             true,
		"eastus2":            true,
		"northcentralus":     true,
		"southcentralus":     true,
		"westcentralus":      true,
		"canadacentral":      true,
		"canadaeast":         true,
		"brazilsouth":        true,
		"northeurope":        true,
		"westeurope":         true,
		"uksouth":            true,
		"ukwest":             true,
		"eastasia":           true,
		"southeastasia":      true,
		"japaneast":          true,
		"japanwest":          true,
		"australiaeast":      true,
		"australiasoutheast": true,
		"centralindia":       true,
		"southindia":         true,
		"westindia":          true,
		"koreacentral":       true,
		"koreasouth":         true,
	}

	if validRegions[region] == false {
		return nil, fmt.Errorf("'%s' is not a valid Azure Region", region)
	}

	config := map[string]interface{}{}
	config["name"] = archiveName
	config["objectStoreType"] = "Azure"
	config["isComputeEnabled"] = true

	config["azureComputeSummary"] = map[string]string{}
	config["azureComputeSummary"].(map[string]string)["tenantId"] = directoryID
	config["azureComputeSummary"].(map[string]string)["subscriptionId"] = strings.Split(virtualNetworkID, "/")[2]
	config["azureComputeSummary"].(map[string]string)["clientId"] = applicationID
	config["azureComputeSummary"].(map[string]string)["region"] = region
	config["azureComputeSummary"].(map[string]string)["generalPurposeStorageAccountName"] = storageAccountName
	config["azureComputeSummary"].(map[string]string)["containerName"] = container

	config["azureComputeSecret"] = map[string]string{}
	config["azureComputeSecret"].(map[string]string)["clientSecret"] = applicationKey

	config["defaultComputeNetworkConfig"] = map[string]string{}
	config["defaultComputeNetworkConfig"].(map[string]string)["subnetId"] = subnetName
	config["defaultComputeNetworkConfig"].(map[string]string)["vNetId"] = virtualNetworkID
	config["defaultComputeNetworkConfig"].(map[string]string)["securityGroupId"] = securityGroupID

	// Create a simplified config that only includes the values returned by Rubrik that can be used for idempotence check
	redactedConfig := map[string]interface{}{}
	redactedConfig["name"] = archiveName
	redactedConfig["objectStoreType"] = "Azure"
	redactedConfig["accessKey"] = storageAccountName
	redactedConfig["bucket"] = container
	redactedConfig["isComputeEnabled"] = true

	redactedConfig["azureComputeSummary"] = map[string]string{}
	redactedConfig["azureComputeSummary"].(map[string]string)["tenantId"] = directoryID
	redactedConfig["azureComputeSummary"].(map[string]string)["subscriptionId"] = strings.Split(virtualNetworkID, "/")[2]
	redactedConfig["azureComputeSummary"].(map[string]string)["clientId"] = applicationID
	redactedConfig["azureComputeSummary"].(map[string]string)["region"] = region
	redactedConfig["azureComputeSummary"].(map[string]string)["generalPurposeStorageAccountName"] = storageAccountName
	redactedConfig["azureComputeSummary"].(map[string]string)["containerName"] = container

	redactedConfig["defaultComputeNetworkConfig"] = map[string]string{}
	redactedConfig["defaultComputeNetworkConfig"].(map[string]string)["subnetId"] = subnetName
	redactedConfig["defaultComputeNetworkConfig"].(map[string]string)["vNetId"] = virtualNetworkID
	redactedConfig["defaultComputeNetworkConfig"].(map[string]string)["securityGroupId"] = securityGroupID

	archivesOnCluster, err := c.Get("internal", "/archive/object_store", httpTimeout)
	if err != nil {
		return nil, err
	}

	for _, v := range archivesOnCluster.(map[string]interface{})["data"].([]interface{}) {
		archiveDefinition := (v.(interface{}).(map[string]interface{})["definition"])

		if archiveDefinition.(map[string]interface{})["objectStoreType"] == "Azure" && archiveDefinition.(map[string]interface{})["name"] == archiveName {

			archivePresent := reflect.DeepEqual(archiveDefinition.(map[string]interface{})["defaultComputeNetworkConfig"], config["defaultComputeNetworkConfig"])
			if archivePresent {
				return fmt.Sprintf("No change required. The '%s' archive location is already configured for CloudOn.", archiveName), nil
			}

			archiveID := (v.(interface{}).(map[string]interface{})["id"])
			apiRequest, err := c.Patch("internal", fmt.Sprintf("/archive/object_store/%s", archiveID), config, httpTimeout)
			if err != nil {
				return nil, err
			}

			return apiRequest, nil

		}

	}
	return nil, fmt.Errorf("The Rubrik cluster does not have an archive location named '%s'", archiveName)

}
