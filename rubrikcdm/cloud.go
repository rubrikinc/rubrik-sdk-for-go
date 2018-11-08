package rubrikcdm

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

// AddAWSNativeAccount
func (c *Credentials) AddAWSNativeAccount(awsAccountName, awsAccessKey, awsSecretKey string, awsRegions []string, regionalBoltNetworkConfigs interface{}, timeout ...int) interface{} {

	c.ClusterVersionCheck(4.2)

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
			log.Fatalf(fmt.Sprintf("Error: '%s' is not a valid AWS Region.", region))
		}

	}

	cloudNativeOnCluster := c.Get("internal", "/aws/account", httpTimeout).(map[string]interface{})["data"]

	if len(cloudNativeOnCluster.([]interface{})) > 0 {
		var currentAWSAccountName string
		var currentAWSConfigID string
		for _, v := range cloudNativeOnCluster.([]interface{}) {
			currentAWSAccountName = (v.(interface{}).(map[string]interface{})["name"].(string))
			currentAWSConfigID = (v.(interface{}).(map[string]interface{})["id"].(string))

		}

		// TODO - Add additional logic checks for the region and bolt configs
		currentAccessKey := c.Get("internal", fmt.Sprintf("/aws/account/%s", currentAWSConfigID), httpTimeout).(interface{}).(map[string]interface{})["accessKey"].(string)
		if currentAccessKey == awsAccessKey {
			return fmt.Sprintf("No change required. Cloud native source with access key '%s' is already configured on the Rubrik cluster.", awsAccessKey)
		}

		if currentAWSAccountName == awsAccountName {
			log.Fatalf(fmt.Sprintf("Error: A Cloud native source with name '%s' already exists. Please enter a unique 'awsAccountName'.", awsAccountName))
		}

	}

	config := map[string]interface{}{}
	config["name"] = awsAccountName
	config["accessKey"] = awsAccessKey
	config["secretKey"] = awsSecretKey
	config["regions"] = awsRegions

	config["regionalBoltNetworkConfigs"] = regionalBoltNetworkConfigs

	return c.Post("internal", "/aws/account", config, httpTimeout)

}

// AWSS3CloudOutRSA configures a new AWS S3 archive target using RSA Key for Encryption.
func (c *Credentials) AWSS3CloudOutRSA(awsBucketName, storageClass, archiveName, awsRegion, awsAccessKey, awsSecretKey, rsaKey string, timeout ...int) interface{} {

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
		log.Fatalf("Error: %s is not a valid AWS Region.", awsRegion)
	}

	if validStorageClass[storageClass] == false {
		log.Fatalf("Error: %s is not a valid 'storageClass'. Please use 'standard', 'standard_ia', or 'reduced_redundancy'.", storageClass)
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

	archivesOnCluster := c.Get("internal", "/archive/object_store", httpTimeout).(map[string]interface{})["data"]

	for _, v := range archivesOnCluster.([]interface{}) {
		archiveDefinition := (v.(interface{}).(map[string]interface{})["definition"])
		delete(archiveDefinition.(map[string]interface{}), "definition")

		archivePresent := reflect.DeepEqual(redactedConfig, archiveDefinition)

		if archivePresent {
			return fmt.Sprintf("No change required. The '%s' archive location is already configured on the Rubrik cluster.", archiveName)
		}

		if archiveDefinition.(map[string]interface{})["objectStoreType"] == "S3" && archiveDefinition.(map[string]interface{})["name"] == archiveName {

			log.Fatalf(fmt.Sprintf("Error: An archive location with the name '%s' already exists. Please enter a unique 'archiveName'", archiveName))

		}

	}

	return c.Post("internal", "/archive/object_store", config, httpTimeout)
}

// AWSS3CloudOutKMS configures a new AWS S3 archive target using a AWS KMS Master Key ID for Encryption.
func (c *Credentials) AWSS3CloudOutKMS(awsBucketName, storageClass, archiveName, awsRegion, awsAccessKey, awsSecretKey, kmsMasterKeyID string, timeout ...int) interface{} {

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
		log.Fatalf("Error: %s is not a valid AWS Region.", awsRegion)
	}

	if validStorageClass[storageClass] == false {
		log.Fatalf("Error: %s is not a valid 'storageClass'. Please use 'standard', 'standard_ia', or 'reduced_redundancy'.", storageClass)
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

	archivesOnCluster := c.Get("internal", "/archive/object_store", httpTimeout).(map[string]interface{})["data"]

	for _, v := range archivesOnCluster.([]interface{}) {
		archiveDefinition := (v.(interface{}).(map[string]interface{})["definition"])
		delete(archiveDefinition.(map[string]interface{}), "definition")

		archivePresent := reflect.DeepEqual(redactedConfig, archiveDefinition)

		if archivePresent {
			return fmt.Sprintf("No change required. The '%s' archive location is already configured on the Rubrik cluster.", archiveName)
		}

		if archiveDefinition.(map[string]interface{})["objectStoreType"] == "S3" && archiveDefinition.(map[string]interface{})["name"] == archiveName {

			log.Fatalf(fmt.Sprintf("Error: An archive location with the name '%s' already exists. Please enter a unique 'archiveName'", archiveName))

		}

	}

	return c.Post("internal", "/archive/object_store", config, httpTimeout)
}

// AWSS3CloudOn
func (c *Credentials) AWSS3CloudOn(archiveName, vpcID, subnetID, securityGroupID string, timeout ...int) interface{} {

	httpTimeout := httpTimeout(timeout)

	config := map[string]interface{}{}
	config["defaultComputeNetworkConfig"] = map[string]interface{}{}
	config["defaultComputeNetworkConfig"].(map[string]interface{})["subnetId"] = subnetID
	config["defaultComputeNetworkConfig"].(map[string]interface{})["vNetId"] = vpcID
	config["defaultComputeNetworkConfig"].(map[string]interface{})["securityGroupId"] = securityGroupID

	archivesOnCluster := c.Get("internal", "/archive/object_store", httpTimeout).(map[string]interface{})["data"]

	for _, v := range archivesOnCluster.([]interface{}) {
		archiveDefinition := (v.(interface{}).(map[string]interface{})["definition"])

		if archiveDefinition.(map[string]interface{})["objectStoreType"] == "S3" && archiveDefinition.(map[string]interface{})["name"] == archiveName {

			archivePresent := reflect.DeepEqual(archiveDefinition.(map[string]interface{})["defaultComputeNetworkConfig"], config["defaultComputeNetworkConfig"])
			if archivePresent {
				return fmt.Sprintf("No change required. The '%s' archive location is already configured for CloudOn.", archiveName)

			}

			archiveID := (v.(interface{}).(map[string]interface{})["id"])
			return c.Patch("internal", fmt.Sprintf("/archive/object_store/%s", archiveID), config, httpTimeout)

		}

	}
	log.Fatalf(fmt.Sprintf("Error: The Rubrik cluster does not have an archive location named '%s'.", archiveName))
	return ""

}

// AzureCloudOut
func (c *Credentials) AzureCloudOut(container, azureAccessKey, storageAccountName, archiveName, instanceType, rsaKey string, timeout ...int) interface{} {

	httpTimeout := httpTimeout(timeout)

	validInstanceTypes := map[string]bool{
		"default":    true,
		"china":      true,
		"germany":    true,
		"government": true,
	}

	if validInstanceTypes[instanceType] == false {
		log.Fatalf(fmt.Sprintf("Error: '%s' is not a valid Azure Instance Type. Valid choices are 'default', 'china', 'germany', or 'government'", instanceType))
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

	archivesOnCluster := c.Get("internal", "/archive/object_store", httpTimeout).(map[string]interface{})["data"]

	for _, v := range archivesOnCluster.([]interface{}) {
		archiveDefinition := (v.(interface{}).(map[string]interface{})["definition"])
		delete(archiveDefinition.(map[string]interface{}), "defaultComputeNetworkConfig")
		delete(archiveDefinition.(map[string]interface{}), "isComputeEnabled")
		delete(archiveDefinition.(map[string]interface{}), "isConsolidationEnabled")
		delete(archiveDefinition.(map[string]interface{}), "azureComputeSummary")

		archivePresent := reflect.DeepEqual(redactedConfig, archiveDefinition)

		if archivePresent {
			return fmt.Sprintf("No change required. The '%s' archive location is already configured on the Rubrik cluster.", archiveName)
		}

		if archiveDefinition.(map[string]interface{})["objectStoreType"] == "Azure" && archiveDefinition.(map[string]interface{})["name"] == archiveName {

			log.Fatalf(fmt.Sprintf("Error: An arhive location with the name '%s' already exists. Please enter a unique 'archiveName'", archiveName))

		}

	}

	return c.Post("internal", "/archive/object_store", config, httpTimeout)

}

// AzureCloudOn
func (c *Credentials) AzureCloudOn(archiveName, container, storageAccountName, applicationID, applicationKey, directoryID, region, virtualNetworkID, subnetName, securityGroupID string, timeout ...int) interface{} {

	httpTimeout := httpTimeout(timeout)

	fmt.Println(httpTimeout)

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
		log.Fatalf(fmt.Sprintf("Error: '%s' is not a valid Azure Region.", region))
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

	archivesOnCluster := c.Get("internal", "/archive/object_store", httpTimeout).(map[string]interface{})["data"]

	for _, v := range archivesOnCluster.([]interface{}) {
		archiveDefinition := (v.(interface{}).(map[string]interface{})["definition"])

		if archiveDefinition.(map[string]interface{})["objectStoreType"] == "Azure" && archiveDefinition.(map[string]interface{})["name"] == archiveName {

			archivePresent := reflect.DeepEqual(archiveDefinition.(map[string]interface{})["defaultComputeNetworkConfig"], config["defaultComputeNetworkConfig"])
			if archivePresent {
				return fmt.Sprintf("No change required. The '%s' archive location is already configured for CloudOn.", archiveName)

			}

			archiveID := (v.(interface{}).(map[string]interface{})["id"])
			return c.Patch("internal", fmt.Sprintf("/archive/object_store/%s", archiveID), config, httpTimeout)

		}

	}
	log.Fatalf(fmt.Sprintf("Error: The Rubrik cluster does not have an archive location named '%s'.", archiveName))
	return ""

}
