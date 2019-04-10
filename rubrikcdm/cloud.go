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
	"time"

	"github.com/mitchellh/mapstructure"
)

// ExportEC2Instance corresponds to GET /aws/ec2_instance/{id}/snapshot
type ExportEC2Instance struct {
	HasMore bool `json:"hasMore"`
	Data    []struct {
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
		AccountID              string   `json:"accountId"`
		InstanceID             string   `json:"instanceId"`
		ImageID                string   `json:"imageId"`
		SnapshotVolumeIds      []string `json:"snapshotVolumeIds"`
	} `json:"data"`
	Total int `json:"total"`
}

// CloudObjectStore represents the JSON response for GET /internal/archive/object_store
type CloudObjectStore struct {
	HasMore bool `json:"hasMore"`
	Data    []struct {
		ID         string `json:"id"`
		Definition struct {
			ObjectStoreType             string `json:"objectStoreType"`
			Name                        string `json:"name"`
			AccessKey                   string `json:"accessKey"`
			Bucket                      string `json:"bucket"`
			PemFileContent              string `json:"pemFileContent"`
			KmsMasterKeyID              string `json:"kmsMasterKeyId"`
			DefaultRegion               string `json:"defaultRegion"`
			Endpoint                    string `json:"endpoint"`
			NumBuckets                  int    `json:"numBuckets"`
			IsComputeEnabled            bool   `json:"isComputeEnabled"`
			IsConsolidationEnabled      bool   `json:"isConsolidationEnabled"`
			DefaultComputeNetworkConfig struct {
				SubnetID        string `json:"subnetId"`
				VNetID          string `json:"vNetId"`
				SecurityGroupID string `json:"securityGroupId"`
				ResourceGroupID string `json:"resourceGroupId"`
			} `json:"defaultComputeNetworkConfig"`
			StorageClass        string `json:"storageClass"`
			AzureComputeSummary struct {
				TenantID                         string `json:"tenantId"`
				SubscriptionID                   string `json:"subscriptionId"`
				ClientID                         string `json:"clientId"`
				Region                           string `json:"region"`
				GeneralPurposeStorageAccountName string `json:"generalPurposeStorageAccountName"`
				ContainerName                    string `json:"containerName"`
				Environment                      string `json:"environment"`
			} `json:"azureComputeSummary"`
		} `json:"definition"`
		GlacierStatus struct {
			RetrievalTier   string `json:"retrievalTier"`
			VaultLockStatus struct {
				FileLockPeriodInDays int       `json:"fileLockPeriodInDays"`
				Status               string    `json:"status"`
				ExpiryTime           time.Time `json:"expiryTime"`
			} `json:"vaultLockStatus"`
		} `json:"glacierStatus"`
		ArchivalProxySummary struct {
			Protocol    string `json:"protocol"`
			ProxyServer string `json:"proxyServer"`
			PortNumber  int    `json:"portNumber"`
			UserName    string `json:"userName"`
		} `json:"archivalProxySummary"`
		ComputeProxySummary struct {
			Protocol    string `json:"protocol"`
			ProxyServer string `json:"proxyServer"`
			PortNumber  int    `json:"portNumber"`
			UserName    string `json:"userName"`
		} `json:"computeProxySummary"`
		ReaderLocationSummary struct {
			State         string    `json:"state"`
			RefreshedTime time.Time `json:"refreshedTime"`
		} `json:"readerLocationSummary"`
	} `json:"data"`
	Total int `json:"total"`
}

// JobStatus represents the JSON response for DELETE /internal/archive/location/{id}
type JobStatus struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	Progress  int    `json:"progress"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	NodeID    string `json:"nodeId"`
	Error     struct {
		Message string `json:"message"`
	} `json:"error"`
	Links []struct {
		Href string `json:"href"`
		Rel  string `json:"rel"`
	} `json:"links"`
}

// CurrentArchiveLocations represents the JSON response for GET /internal/archive/location
type CurrentArchiveLocations struct {
	HasMore bool `json:"hasMore"`
	Data    []struct {
		ID              string `json:"id"`
		Name            string `json:"name"`
		LocationType    string `json:"locationType"`
		IsActive        bool   `json:"isActive"`
		IPAddress       string `json:"ipAddress"`
		Bucket          string `json:"bucket"`
		OwnershipStatus string `json:"ownershipStatus"`
	} `json:"data"`
	Total int `json:"total"`
}

// UpdateArchiveLocations represents the JSON response for PATCH /internal/archive/location/{id}
type UpdateArchiveLocations struct {
	ID         string `json:"id"`
	Definition struct {
		ObjectStoreType             string `json:"objectStoreType"`
		Name                        string `json:"name"`
		AccessKey                   string `json:"accessKey"`
		Bucket                      string `json:"bucket"`
		PemFileContent              string `json:"pemFileContent"`
		KmsMasterKeyID              string `json:"kmsMasterKeyId"`
		DefaultRegion               string `json:"defaultRegion"`
		Endpoint                    string `json:"endpoint"`
		NumBuckets                  int    `json:"numBuckets"`
		IsComputeEnabled            bool   `json:"isComputeEnabled"`
		IsConsolidationEnabled      bool   `json:"isConsolidationEnabled"`
		DefaultComputeNetworkConfig struct {
			SubnetID        string `json:"subnetId"`
			VNetID          string `json:"vNetId"`
			SecurityGroupID string `json:"securityGroupId"`
			ResourceGroupID string `json:"resourceGroupId"`
		} `json:"defaultComputeNetworkConfig"`
		StorageClass        string `json:"storageClass"`
		AzureComputeSummary struct {
			TenantID                         string `json:"tenantId"`
			SubscriptionID                   string `json:"subscriptionId"`
			ClientID                         string `json:"clientId"`
			Region                           string `json:"region"`
			GeneralPurposeStorageAccountName string `json:"generalPurposeStorageAccountName"`
			ContainerName                    string `json:"containerName"`
			Environment                      string `json:"environment"`
		} `json:"azureComputeSummary"`
	} `json:"definition"`
	GlacierStatus struct {
		RetrievalTier   string `json:"retrievalTier"`
		VaultLockStatus struct {
			FileLockPeriodInDays int       `json:"fileLockPeriodInDays"`
			Status               string    `json:"status"`
			ExpiryTime           time.Time `json:"expiryTime"`
		} `json:"vaultLockStatus"`
	} `json:"glacierStatus"`
	ArchivalProxySummary struct {
		Protocol    string `json:"protocol"`
		ProxyServer string `json:"proxyServer"`
		PortNumber  int    `json:"portNumber"`
		UserName    string `json:"userName"`
	} `json:"archivalProxySummary"`
	ComputeProxySummary struct {
		Protocol    string `json:"protocol"`
		ProxyServer string `json:"proxyServer"`
		PortNumber  int    `json:"portNumber"`
		UserName    string `json:"userName"`
	} `json:"computeProxySummary"`
	ReaderLocationSummary struct {
		State         string    `json:"state"`
		RefreshedTime time.Time `json:"refreshedTime"`
	} `json:"readerLocationSummary"`
}

// UpdateAWSNative
type UpdateAWSNative struct {
	Name                       string   `json:"name"`
	AccessKey                  string   `json:"accessKey"`
	Regions                    []string `json:"regions"`
	RegionalBoltNetworkConfigs []struct {
		Region          string `json:"region"`
		VNetID          string `json:"vNetId"`
		SubnetID        string `json:"subnetId"`
		SecurityGroupID string `json:"securityGroupId"`
	} `json:"regionalBoltNetworkConfigs"`
	DisasterRecoveryArchivalLocationID string `json:"disasterRecoveryArchivalLocationId"`
	ID                                 string `json:"id"`
	ConfiguredSLADomainID              string `json:"configuredSlaDomainId"`
	ConfiguredSLADomainName            string `json:"configuredSlaDomainName"`
	PrimaryClusterID                   string `json:"primaryClusterId"`
}

// CurrentAWSAccount represents the JSON response for GET /aws/account
type CurrentAWSAccount struct {
	HasMore bool `json:"hasMore"`
	Data    []struct {
		ID               string `json:"id"`
		Name             string `json:"name"`
		PrimaryClusterID string `json:"primaryClusterId"`
		Status           string `json:"status"`
	} `json:"data"`
	Total int `json:"total"`
}

// CurrentAWSAccountID represents the JSON response for GET /aws/account/{id}
type CurrentAWSAccountID struct {
	Name                       string   `json:"name"`
	AccessKey                  string   `json:"accessKey"`
	Regions                    []string `json:"regions"`
	RegionalBoltNetworkConfigs []struct {
		Region          string `json:"region"`
		VNetID          string `json:"vNetId"`
		SubnetID        string `json:"subnetId"`
		SecurityGroupID string `json:"securityGroupId"`
	} `json:"regionalBoltNetworkConfigs"`
	DisasterRecoveryArchivalLocationID string `json:"disasterRecoveryArchivalLocationId"`
	ID                                 string `json:"id"`
	ConfiguredSLADomainID              string `json:"configuredSlaDomainId"`
	ConfiguredSLADomainName            string `json:"configuredSlaDomainName"`
	PrimaryClusterID                   string `json:"primaryClusterId"`
}

// CloudOn represents the JSON response for PATCH /internal/archive/object_store/{id}
type CloudOn struct {
	ID         string `json:"id"`
	Definition struct {
		ObjectStoreType             string `json:"objectStoreType"`
		Name                        string `json:"name"`
		AccessKey                   string `json:"accessKey"`
		Bucket                      string `json:"bucket"`
		PemFileContent              string `json:"pemFileContent"`
		KmsMasterKeyID              string `json:"kmsMasterKeyId"`
		DefaultRegion               string `json:"defaultRegion"`
		Endpoint                    string `json:"endpoint"`
		NumBuckets                  int    `json:"numBuckets"`
		IsComputeEnabled            bool   `json:"isComputeEnabled"`
		IsConsolidationEnabled      bool   `json:"isConsolidationEnabled"`
		DefaultComputeNetworkConfig struct {
			SubnetID        string `json:"subnetId"`
			VNetID          string `json:"vNetId"`
			SecurityGroupID string `json:"securityGroupId"`
			ResourceGroupID string `json:"resourceGroupId"`
		} `json:"defaultComputeNetworkConfig"`
		StorageClass        string `json:"storageClass"`
		AzureComputeSummary struct {
			TenantID                         string `json:"tenantId"`
			SubscriptionID                   string `json:"subscriptionId"`
			ClientID                         string `json:"clientId"`
			Region                           string `json:"region"`
			GeneralPurposeStorageAccountName string `json:"generalPurposeStorageAccountName"`
			ContainerName                    string `json:"containerName"`
			Environment                      string `json:"environment"`
		} `json:"azureComputeSummary"`
	} `json:"definition"`
	GlacierStatus struct {
		RetrievalTier   string `json:"retrievalTier"`
		VaultLockStatus struct {
			FileLockPeriodInDays int    `json:"fileLockPeriodInDays"`
			Status               string `json:"status"`
			ExpiryTime           string `json:"expiryTime"`
		} `json:"vaultLockStatus"`
	} `json:"glacierStatus"`
	ArchivalProxySummary struct {
		Protocol    string `json:"protocol"`
		ProxyServer string `json:"proxyServer"`
		PortNumber  int    `json:"portNumber"`
		UserName    string `json:"userName"`
	} `json:"archivalProxySummary"`
	ComputeProxySummary struct {
		Protocol    string `json:"protocol"`
		ProxyServer string `json:"proxyServer"`
		PortNumber  int    `json:"portNumber"`
		UserName    string `json:"userName"`
	} `json:"computeProxySummary"`
	ReaderLocationSummary struct {
		State         string `json:"state"`
		RefreshedTime string `json:"refreshedTime"`
	} `json:"readerLocationSummary"`
}

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

	httpTimeout := httpTimeout(timeout)

	minimumClusterVersion := c.ClusterVersionCheck(4.2, httpTimeout)
	if minimumClusterVersion != nil {
		return nil, minimumClusterVersion
	}

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

	// Convert the API Response (map[string]interface{}) to a struct
	var addAccount JobStatus
	mapErr := mapstructure.Decode(apiRequest, &addAccount)
	if mapErr != nil {
		return nil, mapErr
	}

	status, err := c.JobStatus(addAccount.Links[0].Href, httpTimeout)
	if err != nil {
		return nil, err
	}

	return status, nil

}

// ExportEC2Instance exports the latest snapshot of the specified EC2 instance.
//
// The dateTime should be in the following format:  "Month:Day:Year Hour:Minute AM/PM". Ex. 04-09-2019 05:56 PM. You may also use "latest" to export the last
// snapshot taken.
//
// Valid "awsRegion" choices are:
//
//	ap-south-1,ap-northeast-3, ap-northeast-2, ap-southeast-1, ap-southeast-2, ap-northeast-1, ca-central-1, cn-north-1, cn-northwest-1, eu-central-1, eu-west-1,
//	eu-west-2, eu-west-3, us-west-1, us-east-1, us-east-2, and us-west-2.
//
// Valid "instanceType" choices are:
//
// a1.medium, a1.large, a1.xlarge, a1.2xlarge, a1.4xlarge, m4.large, m4.xlarge, m4.2xlarge, m4.4xlarge, m4.10xlarge, m4.16xlarge, m5.large, m5.xlarge, m5.2xlarge,
// m5.4xlarge, m5.12xlarge, m5.24xlarge,  m5a.large, m5a.xlarge, m5a.2xlarge, m5a.4xlarge, m5a.12xlarge, m5a.24xlarge, m5d.large, m5d.xlarge, m5d.2xlarge, m5d.4xlarge,
// m5d.12xlarge, m5d.24xlarge, t2.nano, t2.micro, t2.small, t2.medium, t2.large, t2.xlarge, t2.2xlarge, t3.nano, t3.micro, t3.small, t3.medium, t3.large, t3.xlarge,
// t3.2xlarge, c4.large, c4.xlarge, c4.2xlarge, c4.4xlarge, c4.8xlarge, c5.large, c5.xlarge, c5.2xlarge, c5.4xlarge, c5.9xlarge, c5.18xlarge, c5d.xlarge, c5d.2xlarge,
// c5d.4xlarge, c5d.9xlarge, c5d.18xlarge, c5n.large, c5n.xlarge, c5n.2xlarge, c5n.4xlarge, c5n.9xlarge, c5n.18xlarge, r4.large, r4.xlarge, r4.2xlarge, r4.4xlarge,
// r4.8xlarge, r4.16xlarge, r5.large, r5.xlarge, r5.2xlarge, r5.4xlarge, r5.12xlarge, r5.24xlarge, r5a.large, r5a.xlarge, r5a.2xlarge, r5a.4xlarge, r5a.12xlarge,
// r5a.24xlarge, r5d.large, r5d.xlarge, r5d.2xlarge, r5d.4xlarge, r5d.12xlarge, r5d.24xlarge, x1.16xlarge, x1.32xlarge, x1e.xlarge, x1e.2xlarge, x1e.4xlarge, x1e.8xlarge,
// x1e.16xlarge, x1e.32xlarge, z1d.large, z1d.xlarge, z1d.2xlarge, z1d.3xlarge, z1d.6xlarge, z1d.12xlarge, d2.xlarge, d2.2xlarge, d2.4xlarge, d2.8xlarge, h1.2xlarge,
// h1.4xlarge,  h1.8xlarge, h1.16xlarge, i3.large, i3.xlarge, i3.2xlarge, i3.4xlarge, i3.8xlarge, i3.16xlarge, f1.2xlarge, f1.4xlarge, f1.16xlarge, g3s.xlarge, g3.4xlarge,
// g3.8xlarge, g3.16xlarge, p2.xlarge, p2.8xlarge, p2.16xlarge, p3.2xlarge, p3.8xlarge, p3.16xlarge, and p3dn.24xlarge.
func (c *Credentials) ExportEC2Instance(instanceID, exportedInstanceName, instanceType, awsRegion, subnetID, securityGroupID, dateTime string, waitForCompletion bool, timeout ...int) (interface{}, error) {

	httpTimeout := httpTimeout(timeout)

	minimumClusterVersion := c.ClusterVersionCheck(4.2, httpTimeout)
	if minimumClusterVersion != nil {
		return nil, minimumClusterVersion
	}

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

	validInstanceTypes := map[string]bool{
		"a1.medium":     true,
		"a1.large":      true,
		"a1.xlarge":     true,
		"a1.2xlarge":    true,
		"a1.4xlarge":    true,
		"m4.large":      true,
		"m4.xlarge":     true,
		"m4.2xlarge":    true,
		"m4.4xlarge":    true,
		"m4.10xlarge":   true,
		"m4.16xlarge":   true,
		"m5.large":      true,
		"m5.xlarge":     true,
		"m5.2xlarge":    true,
		"m5.4xlarge":    true,
		"m5.12xlarge":   true,
		"m5.24xlarge":   true,
		"m5a.large":     true,
		"m5a.xlarge":    true,
		"m5a.2xlarge":   true,
		"m5a.4xlarge":   true,
		"m5a.12xlarge":  true,
		"m5a.24xlarge":  true,
		"m5d.large":     true,
		"m5d.xlarge":    true,
		"m5d.2xlarge":   true,
		"m5d.4xlarge":   true,
		"m5d.12xlarge":  true,
		"m5d.24xlarge":  true,
		"t2.nano":       true,
		"t2.micro":      true,
		"t2.small":      true,
		"t2.medium":     true,
		"t2.large":      true,
		"t2.xlarge":     true,
		"t2.2xlarge":    true,
		"t3.nano":       true,
		"t3.micro":      true,
		"t3.small":      true,
		"t3.medium":     true,
		"t3.large":      true,
		"t3.xlarge":     true,
		"t3.2xlarge":    true,
		"c4.large":      true,
		"c4.xlarge":     true,
		"c4.2xlarge":    true,
		"c4.4xlarge":    true,
		"c4.8xlarge":    true,
		"c5.large":      true,
		"c5.xlarge":     true,
		"c5.2xlarge":    true,
		"c5.4xlarge":    true,
		"c5.9xlarge":    true,
		"c5.18xlarge":   true,
		"c5d.xlarge":    true,
		"c5d.2xlarge":   true,
		"c5d.4xlarge":   true,
		"c5d.9xlarge":   true,
		"c5d.18xlarge":  true,
		"c5n.large":     true,
		"c5n.xlarge":    true,
		"c5n.2xlarge":   true,
		"c5n.4xlarge":   true,
		"c5n.9xlarge":   true,
		"c5n.18xlarge":  true,
		"r4.large":      true,
		"r4.xlarge":     true,
		"r4.2xlarge":    true,
		"r4.4xlarge":    true,
		"r4.8xlarge":    true,
		"r4.16xlarge":   true,
		"r5.large":      true,
		"r5.xlarge":     true,
		"r5.2xlarge":    true,
		"r5.4xlarge":    true,
		"r5.12xlarge":   true,
		"r5.24xlarge":   true,
		"r5a.large":     true,
		"r5a.xlarge":    true,
		"r5a.2xlarge":   true,
		"r5a.4xlarge":   true,
		"r5a.12xlarge":  true,
		"r5a.24xlarge":  true,
		"r5d.large":     true,
		"r5d.xlarge":    true,
		"r5d.2xlarge":   true,
		"r5d.4xlarge":   true,
		"r5d.12xlarge":  true,
		"r5d.24xlarge":  true,
		"x1.16xlarge":   true,
		"x1.32xlarge":   true,
		"x1e.xlarge":    true,
		"x1e.2xlarge":   true,
		"x1e.4xlarge":   true,
		"x1e.8xlarge":   true,
		"x1e.16xlarge":  true,
		"x1e.32xlarge":  true,
		"z1d.large":     true,
		"z1d.xlarge":    true,
		"z1d.2xlarge":   true,
		"z1d.3xlarge":   true,
		"z1d.6xlarge":   true,
		"z1d.12xlarge":  true,
		"d2.xlarge":     true,
		"d2.2xlarge":    true,
		"d2.4xlarge":    true,
		"d2.8xlarge":    true,
		"h1.2xlarge":    true,
		"h1.4xlarge":    true,
		"h1.8xlarge":    true,
		"h1.16xlarge":   true,
		"i3.large":      true,
		"i3.xlarge":     true,
		"i3.2xlarge":    true,
		"i3.4xlarge":    true,
		"i3.8xlarge":    true,
		"i3.16xlarge":   true,
		"f1.2xlarge":    true,
		"f1.4xlarge":    true,
		"f1.16xlarge":   true,
		"g3s.xlarge":    true,
		"g3.4xlarge":    true,
		"g3.8xlarge":    true,
		"g3.16xlarge":   true,
		"p2.xlarge":     true,
		"p2.8xlarge":    true,
		"p2.16xlarge":   true,
		"p3.2xlarge":    true,
		"p3.8xlarge":    true,
		"p3.16xlarge":   true,
		"p3dn.24xlarge": true,
	}

	if validInstanceTypes[instanceType] == false {
		return nil, fmt.Errorf("'%s' is not a valid AWS Instance Type", instanceType)
	}

	if validAWSRegions[awsRegion] == false {
		return "", fmt.Errorf("%s is not a valid AWS Region", awsRegion)
	}

	objectID, err := c.ObjectID(instanceID, "ec2", httpTimeout)
	if err != nil {
		return nil, err

	}

	allSnapshots, err := c.Get("internal", fmt.Sprintf("/aws/ec2_instance/%s/snapshot", objectID), httpTimeout)
	if err != nil {
		return nil, err
	}

	// Convert the API Response (map[string]interface{}) to a struct
	var snapshot ExportEC2Instance
	mapErr := mapstructure.Decode(allSnapshots, &snapshot)
	if mapErr != nil {
		return nil, mapErr
	}

	var snapshotID string
	if dateTime == "latest" {
		snapshotID = snapshot.Data[snapshot.Total-1].ID
	} else {
		snapshotTime, err := c.DateTimeConversion(dateTime)
		if err != nil {
			return "", err
		}

		for _, snapshot := range snapshot.Data {
			// The API records EC2 Snapshots down to the second so we have to remove the :00Z from the end to check if the converted dateTime
			// is present in the returned dates
			if strings.Contains(snapshot.Date, strings.TrimSuffix(snapshotTime, ":00Z")) {
				snapshotID = snapshot.ID
			}

		}

		if snapshotID == "" {
			return "", fmt.Errorf("The EC2 Instance '%s' does not have a snapshot take on '%s'", instanceID, dateTime)
		}

	}

	config := map[string]string{}
	config["instanceName"] = exportedInstanceName
	config["instanceType"] = instanceType
	config["region"] = awsRegion
	config["subnetId"] = subnetID
	config["securityGroupId"] = securityGroupID

	exportInstance, err := c.Post("internal", fmt.Sprintf("/aws/ec2_instance/snapshot/%s/export", snapshotID), config, httpTimeout)
	if err != nil {
		return nil, err
	}

	// Convert the API Response (map[string]interface{}) to a struct
	var export JobStatus
	exportMapErr := mapstructure.Decode(exportInstance, &export)
	if exportMapErr != nil {
		return nil, exportMapErr
	}

	if waitForCompletion == true {

		status, err := c.JobStatus(export.Links[0].Href, httpTimeout)
		if err != nil {
			return nil, err
		}

		return status, nil

	}

	return export, nil

}

// RemoveAWSAccount deletes the specific AWS account from the Rubrik clsuter and waits for the job to complete before returning the JobStatus API response.
func (c *Credentials) RemoveAWSAccount(awsAccountName string, deleteExsitingSnapshots bool, timeout ...int) (interface{}, error) {

	httpTimeout := httpTimeout(timeout)

	awsAccountSummary, err := c.AWSAccountSummary(awsAccountName, httpTimeout)
	if err != nil {
		return nil, err
	}

	deleteAPIRequest, err := c.Delete("internal", fmt.Sprintf("/aws/account/%s?delete_existing_snapshots=%t", awsAccountSummary.ID, deleteExsitingSnapshots), httpTimeout)
	if err != nil {
		return nil, err
	}

	// Convert the API Response (map[string]interface{}) to a struct
	var deleteAccount JobStatus
	mapErr := mapstructure.Decode(deleteAPIRequest, &deleteAccount)
	if mapErr != nil {
		return nil, mapErr
	}

	status, err := c.JobStatus(deleteAccount.Links[0].Href)
	if err != nil {
		return nil, err
	}

	return status, nil
}

// UpdateAWSNativeAccount updates the configuration of a AWS Native account. The following values, from PATCH /internal/aws/account/{id} are options for the config:
//  {
//   "name": "string",
//   "accessKey": "string",
//   "secretKey": "string",
//   "regions": [
//     "string"
//   ],
//   "regionalBoltNetworkConfigs": [
//     {
//       "region": "string",
//       "vNetId": "string",
//       "subnetId": "string",
//       "securityGroupId": "string"
//     }
//   ],
//   "disasterRecoveryArchivalLocationId": "string"
// }
func (c *Credentials) UpdateAWSNativeAccount(archiveName string, config map[string]interface{}, timeout ...int) (*UpdateAWSNative, error) {

	httpTimeout := httpTimeout(timeout)

	awsAccountSummary, err := c.AWSAccountSummary(archiveName, httpTimeout)
	if err != nil {
		return nil, err
	}

	patchAPIRequest, err := c.Patch("internal", fmt.Sprintf("/aws/account/%s", awsAccountSummary.ID), config, httpTimeout)
	if err != nil {
		return nil, err
	}

	// Convert the API Response (map[string]interface{}) to a struct
	var patchArchive UpdateAWSNative
	mapErr := mapstructure.Decode(patchAPIRequest, &patchArchive)
	if mapErr != nil {
		return nil, mapErr
	}

	return &patchArchive, nil
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
//	standard, standard_ia, and reduced_redundancy
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
		return "", fmt.Errorf("%s is not a valid AWS Region", awsRegion)
	}

	if validStorageClass[storageClass] == false {
		return "", fmt.Errorf("%s is not a valid 'storageClass'. Please use 'standard', 'standard_ia', or 'reduced_redundancy'", storageClass)
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

	archivesOnCluster, err := c.CloudObjectStore(httpTimeout)
	if err != nil {
		return "", err
	}

	for _, v := range archivesOnCluster.Data {

		// Create a map of the current configuration for easy comparison
		compareRedactedConfig := map[string]interface{}{}
		compareRedactedConfig["objectStoreType"] = v.Definition.ObjectStoreType
		compareRedactedConfig["name"] = v.Definition.Name
		compareRedactedConfig["accessKey"] = v.Definition.AccessKey
		compareRedactedConfig["bucket"] = v.Definition.Bucket
		compareRedactedConfig["defaultRegion"] = v.Definition.DefaultRegion
		compareRedactedConfig["storageClass"] = v.Definition.StorageClass

		archivePresent := reflect.DeepEqual(redactedConfig, compareRedactedConfig)

		if archivePresent {
			return fmt.Sprintf("No change required. The '%s' archive location is already configured on the Rubrik cluster.", archiveName), nil
		}

		if v.Definition.ObjectStoreType == "S3" && v.Definition.Name == archiveName {
			return "", fmt.Errorf("An archive location with the name '%s' already exists. Please enter a unique 'archiveName'", archiveName)
		}

	}

	apiRequest, err := c.Post("internal", "/archive/object_store", config, httpTimeout)
	if err != nil {
		return "", err
	}

	status, err := c.JobStatus(fmt.Sprintf("https://%s/api/internal/archive/location/job/connect/%s", c.NodeIP, apiRequest.(map[string]interface{})["jobInstanceId"].(string)), httpTimeout)
	if err != nil {
		return nil, err
	}

	return status, nil

}

// CloudObjectStore retrieves all archive locations configured on the Rubik cluster.
func (c *Credentials) CloudObjectStore(timeout ...int) (*CloudObjectStore, error) {

	httpTimeout := httpTimeout(timeout)

	apiArchivesOnCluster, err := c.Get("internal", "/archive/object_store", httpTimeout)
	if err != nil {
		return nil, err
	}

	// Convert the API Response (map[string]interface{}) to a struct
	var archivesOnCluster CloudObjectStore
	mapErr := mapstructure.Decode(apiArchivesOnCluster, &archivesOnCluster)
	if mapErr != nil {
		return nil, mapErr
	}

	return &archivesOnCluster, nil

}

// AWSAccountSummary retrives all information from an AWS Native Account.
func (c *Credentials) AWSAccountSummary(awsAccountName string, timeout ...int) (*CurrentAWSAccountID, error) {

	httpTimeout := httpTimeout(timeout)

	apiAWSAccounts, err := c.Get("internal", "/aws/account", httpTimeout)
	if err != nil {
		return nil, err
	}

	// Convert the API Response (map[string]interface{}) to a struct
	var awsAccountsOnCluster CurrentAWSAccount
	mapErr := mapstructure.Decode(apiAWSAccounts, &awsAccountsOnCluster)
	if mapErr != nil {
		return nil, mapErr
	}

	var accountID string
	for _, v := range awsAccountsOnCluster.Data {
		if v.Name == awsAccountName {
			accountID = v.ID
		}
	}
	if len(accountID) == 0 {
		return nil, fmt.Errorf("The %s AWS Native Account was not found on the Rubrik cluster", awsAccountName)
	}

	apiAWSAccountsID, err := c.Get("internal", fmt.Sprintf("/aws/account/%s", accountID), httpTimeout)
	if err != nil {
		return nil, err
	}

	// Convert the API Response (map[string]interface{}) to a struct
	var awsAccountID CurrentAWSAccountID
	idMapErr := mapstructure.Decode(apiAWSAccountsID, &awsAccountID)
	if mapErr != nil {
		return nil, idMapErr
	}

	return &awsAccountID, nil

}

// RemoveArchiveLocation deletes the archival location from the SLA Domains that reference it and expire all snapshots at the archival location
func (c *Credentials) RemoveArchiveLocation(archiveName string, timeout ...int) (*JobStatus, error) {

	httpTimeout := httpTimeout(timeout)

	// Search the Rubrik cluster for all current archive locations
	currentArchivesRequest, err := c.Get("internal", fmt.Sprintf("/archive/location?name=%s", archiveName), httpTimeout)
	if err != nil {
		return nil, err
	}

	// Convert the API Response (map[string]interface{}) to a struct
	var currentArchive CurrentArchiveLocations
	currentArchiveMapErr := mapstructure.Decode(currentArchivesRequest, &currentArchive)
	if currentArchiveMapErr != nil {
		return nil, currentArchiveMapErr

	}

	// Get the ID of the archive location
	var archiveID string
	for _, v := range currentArchive.Data {
		if v.Name == archiveName {
			archiveID = v.ID
		}
	}

	if archiveID == "" {
		return nil, fmt.Errorf("No change required. The Rubrik cluster does not contain a archive location named '%s'", archiveName)
	}

	// Pause archive activity on the archive location before deleting
	_, pauseErr := c.Post("internal", fmt.Sprintf("/archive/location/%s/owner/pause", archiveID), httpTimeout)
	if pauseErr != nil {
		// If the archive location is already paused do not return an error message
		if strings.Contains(pauseErr.Error(), "already paused") != true {
			return nil, pauseErr
		}

	}

	deleteAPIRequest, err := c.Delete("internal", fmt.Sprintf("/archive/location/%s", archiveID), httpTimeout)
	if err != nil {
		return nil, err
	}

	// Convert the API Response (map[string]interface{}) to a struct
	var deleteArchive JobStatus
	mapErr := mapstructure.Decode(deleteAPIRequest, &deleteArchive)
	if mapErr != nil {
		return nil, mapErr
	}

	return &deleteArchive, nil
}

// UpdateCloudArchiveLocation updates the configuration of a the Cloud Archival Location. The following values, from PATCH /internal/object_store/{id} are options for the config:
//  {
//     "name": "string",
//     "accessKey": "string",
//     "secretKey": "string",
//     "endpoint": "string",
//     "numBuckets": 0,
//     "isComputeEnabled": true,
//     "isConsolidationEnabled": true,
//     "defaultComputeNetworkConfig": {
//       "subnetId": "string",
//       "vNetId": "string",
//       "securityGroupId": "string",
//       "resourceGroupId": "string"
//     },
//     "storageClass": "string",
//     "glacierConfig": {
//       "retrievalTier": "BulkRetrieval",
//       "vaultLockPolicy": {
//         "fileLockPeriodInDays": 0
//       }
//     },
//     "azureComputeSummary": {
//       "tenantId": "string",
//       "subscriptionId": "string",
//       "clientId": "string",
//       "region": "string",
//       "generalPurposeStorageAccountName": "string",
//       "containerName": "string",
//       "environment": "AZURE"
//     },
//     "azureComputeSecret": {
//     "  clientSecret": "string"
//     },
//     "archivalProxyConfig": {
//       "protocol": "HTTP",
//       "proxyServer": "string",
//       "portNumber": 0,
//       "userName": "string",
//       "password": "string"
//     },
//     "computeProxyConfig": {
//       "protocol": "HTTP",
//       "proxyServer": "string",
//       "portNumber": 0,
//       "userName": "string",
//       "password": "string"
//     }
//   }
func (c *Credentials) UpdateCloudArchiveLocation(archiveName string, config map[string]interface{}, timeout ...int) (*UpdateArchiveLocations, error) {

	httpTimeout := httpTimeout(timeout)

	// Search the Rubrik cluster for all current archive locations
	currentArchivesRequest, err := c.Get("internal", fmt.Sprintf("/archive/location?name=%s", archiveName), httpTimeout)
	if err != nil {
		return nil, err
	}

	// Convert the API Response (map[string]interface{}) to a struct
	var currentArchive CurrentArchiveLocations
	currentArchiveMapErr := mapstructure.Decode(currentArchivesRequest, &currentArchive)
	if currentArchiveMapErr != nil {
		return nil, currentArchiveMapErr

	}

	// Get the ID of the archive location
	var archiveID string
	for _, v := range currentArchive.Data {
		if v.Name == archiveName {
			archiveID = v.ID
		}
	}
	if archiveID == "" {
		return nil, fmt.Errorf("No change required. The Rubrik cluster does not contain a archive location named '%s'", archiveName)
	}

	patchAPIRequest, err := c.Patch("internal", fmt.Sprintf("/archive/object_store/%s", archiveID), config, httpTimeout)
	if err != nil {
		return nil, err
	}

	// Convert the API Response (map[string]interface{}) to a struct
	var patchArchive UpdateArchiveLocations
	mapErr := mapstructure.Decode(patchAPIRequest, &patchArchive)
	if mapErr != nil {
		return nil, mapErr
	}

	return &patchArchive, nil
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
//	standard, standard_ia, and reduced_redundancy
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
		delete(archiveDefinition.(map[string]interface{}), "isComputeEnabled")
		delete(archiveDefinition.(map[string]interface{}), "isConsolidationEnabled")

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

	status, err := c.JobStatus(fmt.Sprintf("https://%s/api/internal/archive/location/job/connect/%s", c.NodeIP, apiRequest.(map[string]interface{})["jobInstanceId"].(string)), httpTimeout)
	if err != nil {
		return nil, err
	}

	return status, nil

}

// AWSS3CloudOn provides the ability to convert a vSphere virtual machines snapshot, an archived snapshot, or a replica into an Amazon Machine Image (AMI)
// and then launch that AMI into an Elastic Compute Cloud (EC2) instance on an Amazon Virtual Private Cloud (VPC).
//
// The function will return one of the following:
//	- No change required. The '{archiveName}' archive location is already configured for CloudOn.
//
//	- The full API response for PATCH /internal/archive/object_store.
func (c *Credentials) AWSS3CloudOn(archiveName, vpcID, subnetID, securityGroupID string, timeout ...int) (*CloudOn, error) {

	httpTimeout := httpTimeout(timeout)

	config := map[string]interface{}{}
	config["defaultComputeNetworkConfig"] = map[string]string{}
	config["defaultComputeNetworkConfig"].(map[string]string)["vNetId"] = vpcID
	config["defaultComputeNetworkConfig"].(map[string]string)["subnetId"] = subnetID
	config["defaultComputeNetworkConfig"].(map[string]string)["securityGroupId"] = securityGroupID

	archivesOnCluster, err := c.CloudObjectStore(httpTimeout)
	if err != nil {
		return nil, err
	}

	for _, v := range archivesOnCluster.Data {

		if v.Definition.ObjectStoreType == "S3" && v.Definition.Name == archiveName {

			archivePresent := reflect.DeepEqual(v.Definition.DefaultComputeNetworkConfig, config["defaultComputeNetworkConfig"])
			if archivePresent {
				return nil, fmt.Errorf("No change required. The '%s' archive location is already configured for CloudOn", archiveName)
			}

			apiRequest, err := c.Patch("internal", fmt.Sprintf("/archive/object_store/%s", v.ID), config, httpTimeout)
			if err != nil {
				return nil, err
			}

			// Convert the API Response (map[string]interface{}) to a struct
			var apiResponse CloudOn
			mapErr := mapstructure.Decode(apiRequest, &apiResponse)
			if mapErr != nil {
				return nil, mapErr
			}
			return &apiResponse, err

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
			return nil, fmt.Errorf("An archive location with the name '%s' already exists. Please enter a unique 'archiveName'", archiveName)
		}

	}

	apiRequest, err := c.Post("internal", "/archive/object_store", config, httpTimeout)
	if err != nil {
		return "", err
	}

	status, err := c.JobStatus(fmt.Sprintf("https://%s/api/internal/archive/location/job/connect/%s", c.NodeIP, apiRequest.(map[string]interface{})["jobInstanceId"].(string)), httpTimeout)
	if err != nil {
		return nil, err
	}

	return status, nil

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
func (c *Credentials) AzureCloudOn(archiveName, container, storageAccountName, applicationID, applicationKey, directoryID, region, virtualNetworkID, subnetName, securityGroupID string, timeout ...int) (*CloudOn, error) {

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
	config["isComputeEnabled"] = true

	config["azureComputeSummary"] = map[string]string{}
	config["azureComputeSummary"].(map[string]string)["tenantId"] = directoryID
	config["azureComputeSummary"].(map[string]string)["subscriptionId"] = strings.Split(virtualNetworkID, "/")[2]
	config["azureComputeSummary"].(map[string]string)["clientId"] = applicationID
	config["azureComputeSummary"].(map[string]string)["region"] = region
	config["azureComputeSummary"].(map[string]string)["generalPurposeStorageAccountName"] = storageAccountName
	config["azureComputeSummary"].(map[string]string)["containerName"] = container
	config["azureComputeSummary"].(map[string]string)["environment"] = "AZURE"

	config["azureComputeSecret"] = map[string]string{}
	config["azureComputeSecret"].(map[string]string)["clientSecret"] = applicationKey

	config["defaultComputeNetworkConfig"] = map[string]string{}
	config["defaultComputeNetworkConfig"].(map[string]string)["subnetId"] = subnetName
	config["defaultComputeNetworkConfig"].(map[string]string)["vNetId"] = virtualNetworkID
	config["defaultComputeNetworkConfig"].(map[string]string)["securityGroupId"] = securityGroupID
	config["defaultComputeNetworkConfig"].(map[string]string)["resourceGroupId"] = strings.Split(virtualNetworkID, "/")[4]

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
				return nil, fmt.Errorf("No change required. The '%s' archive location is already configured for CloudOn", archiveName)
			}

			archiveID := (v.(interface{}).(map[string]interface{})["id"])
			apiRequest, err := c.Patch("internal", fmt.Sprintf("/archive/object_store/%s", archiveID), config, httpTimeout)
			if err != nil {
				return nil, err
			}

			// Convert the API Response (map[string]interface{}) to a struct
			var apiResponse CloudOn
			mapErr := mapstructure.Decode(apiRequest, &apiResponse)
			if mapErr != nil {
				return nil, mapErr
			}

			return &apiResponse, nil

		}

	}
	return nil, fmt.Errorf("The Rubrik cluster does not have an archive location named '%s'", archiveName)

}
