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

package rubrikcdm_test

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm"
)

func ExampleCredentials_ExportEC2Instance() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	instanceID := "i-0123456789abcdefg"
	exportInstanceName := "Go SDK"
	instanceType := "m4.large"
	awsRegion := "us-east-2"
	subnetID := "subnet-0123456789abcdefg"
	securityGroupID := "sg-0123456789abcdefg"
	dateTime := "04-09-2019 05:56 PM"
	waitForCompletion := true

	exportEC2, err := rubrik.ExportEC2Instance(instanceID, exportInstanceName, instanceType, awsRegion, subnetID, securityGroupID, dateTime, waitForCompletion)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleCredentials_ClusterVersion() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	clusterVersion, err := rubrik.ClusterVersion()
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleCredentials_BeginManagedVolumeSnapshot() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	mvName := "GoSDK"

	beginMV, err := rubrik.BeginManagedVolumeSnapshot(mvName)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleCredentials_PauseSnapshot() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	vmName := "vm01"

	pauseVM, err := rubrik.PauseSnapshot(vmName, "VMware")
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleCredentials_OnDemandSnapshotVM() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	vmName := "ansible-node01"
	sla := "current"

	vmSnapshot, err := rubrik.OnDemandSnapshotVM(vmName, "VMware", sla)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleCredentials_OnDemandSnapshotPhysical() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	hostname := "vm01"
	slaName := "current"
	fileset := "C_Drive"
	hostOS := "Windows"

	hostSnapshot, err := rubrik.OnDemandSnapshotPhysical(hostname, slaName, fileset, hostOS)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleCredentials_ResumeSnapshot() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	vmName := "vm01"

	resumeVM, err := rubrik.ResumeSnapshot(vmName, "VMware")
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleCredentials_GetSLAObjects() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	slaName := "Gold"

	getObjSLA, err := rubrik.GetSLAObjects(slaName, "VMware")
	if err != nil {
		log.Fatal(err)
	}

}

func ExampleCredentials_EndManagedVolumeSnapshot() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	mvName := "GoSDK"
	slaName := "Gold"

	endMV, err := rubrik.EndManagedVolumeSnapshot(mvName, slaName)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleCredentials_AssignSLA() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	objectName := "vm01"
	slaName := "Bronze"

	assignSLA, err := rubrik.AssignSLA(objectName, "vmware", slaName)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleCredentials_ConfigureTimezone() {
	rubrik, err := rubrikcdm.ConnectEnv()

	confTimezone, err := rubrik.ConfigureTimezone("America/Los_Angeles")
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleCredentials_ConfigureVLAN() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	vlanIPs := map[string]string{}
	vlanIPs["RVM157S018901"] = "192.168.100.100"
	vlanIPs["RVM157S018902"] = "192.168.100.101"
	vlanIPs["RVM157S018903"] = "192.168.100.102"
	vlanIPs["RVM157S0189014"] = "192.168.100.103"

	vlan := 100
	netmask := "255.255.255.0"

	configVLAN, err := rubrik.ConfigureVLAN(netmask, vlan, vlanIPs)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleCredentials_AddvCenter() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	vCenterIP := "vcsa.rubrikgosdk.lab"
	vCenterUsername := "go"
	vCenterPassword := "sdk"
	vmLinking := true

	addVcenter, err := rubrik.AddvCenter(vCenterIP, vCenterUsername, vCenterPassword, vmLinking)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleCredentials_AddvCenterWithCert() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	vCenterIP := "vcsa.rubrikgosdk.lab"
	vCenterUsername := "go"
	vCenterPassword := "sdk"
	readcaCertificate, _ := ioutil.ReadFile("ca_cert")
	caCert := string(readcaCertificate)
	vmLinking := true

	addVcenterWithCert, err := rubrik.AddvCenterWithCert(vCenterIP, vCenterUsername, vCenterPassword, caCert, vmLinking)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleCredentials_ConfigureSMTPSettings() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	hostname := "smtp.rubrikgosdk.lab"
	port := 100
	fromEmail := "gosdk@rubrikgosdk.lab"
	username := "go"
	password := "sdk"
	encryption := "NONE"

	smtpConfig, err := rubrik.ConfigureSMTPSettings(hostname, fromEmail, username, password, encryption, port)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleCredentials_ConfigureSearchDomain() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	searchDomains := []string{"rubrikgosdk.lab"}

	searchDomainConfig, err := rubrik.ConfigureSearchDomain(searchDomains)
	if err != nil {
		log.Fatal(err)
	}

}

func ExampleCredentials_ObjectID() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	slaName := "Gold"
	timeout := 15

	slaID, err := rubrik.ObjectID(slaName, "sla", timeout)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleCredentials_Bootstrap() {

	bootstrapNode := "192.168.101.100"
	rubrik := rubrikcdm.Connect(bootstrapNode, "", "")

	clusterName := "Go-SDK"
	adminEmail := "gosdk@rubrikgosdk.lab"
	adminPassword := "RubrikGoSDK"
	managementGateway := "192.168.101.1"
	managementSubnetMask := "255.255.255.0"
	dnsSearchDomain := []string{"rubrikgosdk.lab"}
	dnsNameServers := []string{"192.168.100.5", "192.168.100.6"}
	ntpServers := map[string]interface{}{}
	ntpServers["ntpServer1"] = map[string]interface{}{}
	ntpServers["ntpServer1"].(map[string]interface{})["IP"] = "192.168.100.5"
	ntpServers["ntpServer2"] = map[string]interface{}{}
	ntpServers["ntpServer2"].(map[string]interface{})["IP"] = "192.168.100.6"	enableEncryption := true // set to false for a Cloud Cluster
	waitForCompletion := true

	nodeConfig := map[string]string{}
	nodeConfig["RVM1234567890"] = bootstrapNode
	nodeConfig["RVM1234567891"] = "192.168.101.101"
	nodeConfig["RVM1234567892"] = "192.168.101.102"
	nodeConfig["RVM1234567893"] = "192.168.101.103"

	bootstrap, err := rubrik.Bootstrap(clusterName, adminEmail, adminPassword, managementGateway, managementSubnetMask, dnsSearchDomain, dnsNameServers, ntpServers, nodeConfig, enableEncryption, waitForCompletion)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleCredentials_Bootstrap_with_Secure_NTP() {

	bootstrapNode := "192.168.101.100"
	rubrik := rubrikcdm.Connect(bootstrapNode, "", "")

	clusterName := "Go-SDK"
	adminEmail := "gosdk@rubrikgosdk.lab"
	adminPassword := "RubrikGoSDK"
	managementGateway := "192.168.101.1"
	managementSubnetMask := "255.255.255.0"
	dnsSearchDomain := []string{"rubrikgosdk.lab"}
	dnsNameServers := []string{"192.168.100.5", "192.168.100.6"}
	ntpServers := map[string]interface{}{}
	ntpServers["ntpServer1"] = map[string]interface{}{}
	ntpServers["ntpServer1"].(map[string]interface{})["IP"] = "192.168.100.5"
	ntpServers["ntpServer1"].(map[string]interface{})["key"] = "key-material-for-ntpserver-1"
	ntpServers["ntpServer1"].(map[string]interface{})["keyId"] = 0
	ntpServers["ntpServer1"].(map[string]interface{})["keyType"] = "MD5"
	ntpServers["ntpServer2"] = map[string]interface{}{}
	ntpServers["ntpServer2"].(map[string]interface{})["IP"] = "192.168.100.6"
	ntpServers["ntpServer2"].(map[string]interface{})["keyId"] = 1
	ntpServers["ntpServer2"].(map[string]interface{})["key"] = "key-material-for-ntpserver-2"
	ntpServers["ntpServer2"].(map[string]interface{})["keyType"] = "SHA1"
	enableEncryption := true // set to false for a Cloud Cluster
	waitForCompletion := true

	nodeConfig := map[string]string{}
	nodeConfig["RVM1234567890"] = bootstrapNode
	nodeConfig["RVM1234567891"] = "192.168.101.101"
	nodeConfig["RVM1234567892"] = "192.168.101.102"
	nodeConfig["RVM1234567893"] = "192.168.101.103"

	bootstrap, err := rubrik.Bootstrap(clusterName, adminEmail, adminPassword, managementGateway, managementSubnetMask, dnsSearchDomain, dnsNameServers, ntpServers, nodeConfig, enableEncryption, waitForCompletion)
	if err != nil {
		log.Fatal(err)
	}

}

func ExampleCredentials_BootstrapAws() {

	bootstrapNode := "192.168.102.100"
	rubrik := rubrikcdm.Connect(bootstrapNode, "", "")

	clusterName := "Go-SDK"
	adminEmail := "gosdk@rubrikgosdk.lab"
	adminPassword := "RubrikGoSDK"
	managementGateway := "192.168.102.1"
	managementSubnetMask := "255.255.255.0"
	dnsSearchDomain := []string{"rubrikgosdk.lab"}
	dnsNameServers := []string{"192.168.100.5", "192.168.100.6"}
	ntpServers := map[string]interface{}{}
	ntpServers["ntpServer1"] = map[string]interface{}{}
	ntpServers["ntpServer1"].(map[string]interface{})["IP"] = "192.168.100.5"
	ntpServers["ntpServer2"] = map[string]interface{}{}
	ntpServers["ntpServer2"].(map[string]interface{})["IP"] = "192.168.100.6"	
	bucketName := "s3-bucket-for-cces-aws"
	waitForCompletion := true
	enableEncryption := false // set to false for a Cloud Cluster

	nodeConfig := map[string]string{}
	nodeConfig["CCESAWSNODE1"] = bootstrapNode
	nodeConfig["CCESAWSNODE2"] = "192.168.102.101"
	nodeConfig["CCESAWSNODE3"] = "192.168.102.102"

	_, err := rubrik.BootstrapCcesAws(clusterName, adminEmail, adminPassword, managementGateway, managementSubnetMask, dnsSearchDomain, dnsNameServers, ntpServers, nodeConfig, enableEncryption, bucketName, waitForCompletion)
	if err != nil {
		log.Fatal(err)
	}

}

func ExampleCredentials_BootstrapAzure() {

	bootstrapNode := "192.168.103.100"
	rubrik := rubrikcdm.Connect(bootstrapNode, "", "")

	clusterName := "Go-SDK"
	adminEmail := "gosdk@rubrikgosdk.lab"
	adminPassword := "RubrikGoSDK"
	managementGateway := "192.168.103.1"
	managementSubnetMask := "255.255.255.0"
	dnsSearchDomain := []string{"rubrikgosdk.lab"}
	dnsNameServers := []string{"192.168.100.5", "192.168.100.6"}
	ntpServers := map[string]interface{}{}
	ntpServers["ntpServer1"] = map[string]interface{}{}
	ntpServers["ntpServer1"].(map[string]interface{})["IP"] = "192.168.100.5"
	ntpServers["ntpServer2"] = map[string]interface{}{}
	ntpServers["ntpServer2"].(map[string]interface{})["IP"] = "192.168.100.6"
	connectionString := "DefaultEndpointsProtocol=https;AccountName=storageaccountforccesazuregosdk;AccountKey=abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890abcdefghijklm==;EndpointSuffix=core.windows.net"
	containerName := "container-for-cces-azure-gosdk"
	enableEncryption := false // set to false for a Cloud Cluster
	waitForCompletion := true

	nodeConfig := map[string]string{}
	nodeConfig["CCESAZURENODE1"] = bootstrapNode
	nodeConfig["CCESAZURENODE2"] = "192.168.103.101"
	nodeConfig["CCESAZURENODE3"] = "192.168.103.102"

	_, err := rubrik.BootstrapCcesAzure(clusterName, adminEmail, adminPassword, managementGateway, managementSubnetMask, dnsSearchDomain, dnsNameServers, ntpServers, nodeConfig, enableEncryption, connectionString, containerName, waitForCompletion)
	if err != nil {
		log.Fatal(err)
	}

}

func ExampleCredentials_ConfigureDNSServers() {
	rubrik, err := rubrikcdm.ConnectEnv()

	dnsServers := []string{"192.168.100.5", "192.168.100.6"}

	dnsConfig, err := rubrik.ConfigureDNSServers(dnsServers)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleCredentials_ConfigureSyslog() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	syslogIP := "192.168.100.7"
	syslogProtocol := "UDP"
	syslogPort := 514

	syslog, err := rubrik.ConfigureSyslog(syslogIP, syslogProtocol, syslogPort)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleCredentials_RegisterCluster() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	support_portal_username := "gosdk@rubrikgosdk.lab"
	support_portal_password := "GoDummyPassword"

	register, err := rubrik.RegisterCluster(support_portal_username, support_portal_password)
	if err != nil {
		log.Fatal(err)
	}
}

// NOTE: This code is broken for releases of Rubrik CDM v5.0 and later.
func ExampleCredentials_ConfigureNTP() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	ntpServers := []string{"192.168.100.5", "192.168.100.6"}

	ntp, err := rubrik.ConfigureNTP(ntpServers)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleCredentials_ClusterNodeIP() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	clusterVersion, err := rubrik.ClusterNodeIP()
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleCredentials_EndUserAuthorization() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	vmName := "vm01"
	endUser := "user01"

	endUserAuth, err := rubrik.EndUserAuthorization(vmName, endUser, "VMware")
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleCredentials_ClusterVersionCheck() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	clusterVersion := rubrik.ClusterVersionCheck(4.2)
	if clusterVersion != nil {
		log.Fatal(err)
	}
}

func ExampleCredentials_Get() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	clusterInfo, err := rubrik.Get("v1", "/cluster/me")
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleCredentials_Post() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	config := map[string]string{}
	config["slaId"] = "388a473c-3361-42ab-8f5b-08edb76891f6"

	onDemandSnapshot, err := rubrik.Post("v1", "/VMware/vm/VirtualMachine:::fbcb1f51-9520-4227-a68c-6fe145982f48-vm-204969/snapshot", config)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleCredentials_Patch() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	config := map[string]string{}
	config["configuredSlaDomainId"] = "388a473c-3361-42ab-8f5b-08edb76891f6"

	fileset, err := rubrik.Patch("v1", "/fileset/Fileset:::b95456e2-7d60-4ed0-af88-648516e139a6", config)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleCredentials_Delete() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	deleteSLA, err := rubrik.Delete("v1", "/sla_domain/388a473c-3361-42ab-8f5b-08edb76891f6")
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleCredentials_AddAWSNativeAccount() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	awsAccountName := "GO SDK Demo" // This is the name that will be displayed in the Rubrik UI
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsRegions := []string{"us-east-1"}

	usEast1 := map[string]string{}
	usEast1["region"] = "us-east-1"
	usEast1["vNetId"] = "vpc-01234567"
	usEast1["subnetId"] = "subnet-01234567"
	usEast1["securityGroupId"] = "sg-01234567"
	boltConfig := []interface{}{usEast1}

	addAWSNative, err := rubrik.AddAWSNativeAccount(awsAccountName, awsAccessKey, awsSecretKey, awsRegions, boltConfig)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleCredentials_RefreshvCenter() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	vcenter_hostname := "vcsa.rubrikgosdk.lab"

	refresh, err := rubrik.RefreshvCenter(vcenter_hostname)
	if err != nil {
		log.Fatal(err)

	}
}

func ExampleCredentials_RemoveAWSAccount() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	awsAccountName := "GO SDK Demo" // This is the name that will be displayed in the Rubrik UI
	deleteSnapshots := true

	removeAWSAccount, err := rubrik.RemoveAWSAccount(awsAccountName, deleteSnapshots)
	if err != nil {
		log.Fatal(err)
	}

}

func ExampleCredentials_UpdateAWSNativeAccount() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	awsAccountName := "GO SDK Demo" // This is the name that will be displayed in the Rubrik UI
	newAWSAccountName := "GO SDK"

	config := make(map[string]interface{})
	config["name"] = newAWSAccountName

	updateAWSAccount, err := rubrik.UpdateAWSNativeAccount(awsAccountName, config)
	if err != nil {
		log.Fatal(err)
	}

}

func ExampleCredentials_CloudObjectStore() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	archiveLocations, err := rubrik.CloudObjectStore()
	if err != nil {
		log.Fatal(err)
	}

}

func ExampleCredentials_AWSAccountSummary() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	awasAccountName := "Go SDK"

	awsSummary, err := rubrik.AWSAccountSummary(awasAccountName)
	if err != nil {
		log.Fatal(err)
	}

}

func ExampleCredentials_RemoveArchiveLocation() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	archiveName := "Go SDK"

	removeArchive, err := rubrik.RemoveArchiveLocation(archiveName)
	if err != nil {
		log.Fatal(err)
	}

}

func ExampleCredentials_UpdateCloudArchiveLocation() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	archiveName := "Go SDK"
	storageClass := "standard"

	config := make(map[string]interface{})
	config["storageClass"] = storageClass

	updateArchive, err := rubrik.UpdateCloudArchiveLocation(archiveName, config)
	if err != nil {
		log.Fatal(err)
	}

}

func ExampleCredentials_AWSS3CloudOutRSA() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	awsBucket := "rubrikgosdk"
	storageClass := "standard"
	archiveName := "AWS:S3:GoSDK"
	awsRegion := "us-east-1"
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	readRSAKey, _ := ioutil.ReadFile("rsa_key.pem")
	rsaKey := string(readRSAKey)

	awsCloudOut, err := rubrik.AWSS3CloudOutRSA(awsBucket, storageClass, archiveName, awsRegion, awsAccessKey, awsSecretKey, rsaKey)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleCredentials_AWSS3CloudOutKMS() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	awsBucket := "rubrikgosdk"
	storageClass := "standard"
	archiveName := "AWS:S3:GoSDK"
	awsRegion := "us-east-1"
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	kmsMasterKeyID := os.Getenv("AWS_MASTER_KEY_ID")

	awsCloudOut, err := rubrik.AWSS3CloudOutKMS(awsBucket, storageClass, archiveName, awsRegion, awsAccessKey, awsSecretKey, kmsMasterKeyID)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleCredentials_AWSS3CloudOn() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	archiveName := "AWS:S3:GoSDK"
	vpcID := "vpc-01234567"
	subnetID := "subnet-01234567"
	securityGroupID := "sg-01234567"

	awsCloudOn, err := rubrik.AWSS3CloudOn(archiveName, vpcID, subnetID, securityGroupID)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleCredentials_AzureCloudOut() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	container := "gosdk"
	azureAccessKey := os.Getenv("AZURE_ACCESS_KEY")
	storageAccountName := "rubrikgosdk"
	archiveName := "Azure:gosdk"
	instanceType := "default"
	readRSAKey, _ := ioutil.ReadFile("rsa_key.pem")
	rsaKey := string(readRSAKey)

	azureCloudOut, err := rubrik.AzureCloudOut(container, azureAccessKey, storageAccountName, archiveName, instanceType, rsaKey)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleCredentials_AzureCloudOn() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	archiveName := "Azure:GoSDK"
	container := "gosdk"
	storageAccountName := "rubrikgosdk"
	applicationID := os.Getenv("AZURE_APP_ID")
	applicationKey := os.Getenv("AZURE_APP_KEY")

	directoryID := os.Getenv("AZURE_DIRECTORY_ID")
	region := "westus2"
	virtualNetworkID := os.Getenv("AZURE_VNET_ID")
	subnetName := os.Getenv("AZURE_SUBNET")
	securityGroupID := os.Getenv("AZURE_SECURITY_GROUP")

	azureCloudOn, err := rubrik.AzureCloudOn(archiveName, container, storageAccountName, applicationID, applicationKey, directoryID, region, virtualNetworkID, subnetName, securityGroupID)
	if err != nil {
		log.Fatal(err)
	}

}

func ExampleCredentials_RecoverFileDownload() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	hostName := "rubrik-sql01rubrikgosdk.lab"
	fileset := "fileset01"
	hostOS := "Linux"
	dateTime := "04-17-2020 12:49 PM"
	filePath := "/rubrik/bu-tar01/test01"

	fileDownload, err := rubrik.RecoverFileDownload(hostName, fileset, hostOS, filePath, dateTime)
	if err != nil {
		log.Fatal(err)
	}
}
