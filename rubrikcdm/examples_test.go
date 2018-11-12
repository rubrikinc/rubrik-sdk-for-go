package rubrikcdm_test

import (
	"io/ioutil"
	"os"

	"github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm"
)

func ExampleCredentials_ClusterVersion() {
	rubrik := rubrikcdm.ConnectEnv()

	clusterVersion := rubrik.ClusterVersion()
}

func ExampleCredentials_BeginManagedVolumeSnapshot() {
	rubrik := rubrikcdm.ConnectEnv()

	mvName := "GoSDK"

	beginMV := rubrik.BeginManagedVolumeSnapshot(mvName)
}

func ExampleCredentials_PauseSnapshot() {
	rubrik := rubrikcdm.ConnectEnv()

	vmName := "vm01"

	pauseVM := rubrik.PauseSnapshot(vmName, "vmware")
}

func ExampleCredentials_OnDemandSnapshotVM() {
	rubrik := rubrikcdm.ConnectEnv()

	vmName := "ansible-node01"
	sla := "current"

	vmSnapshot := rubrik.OnDemandSnapshotVM(vmName, "vmware", sla)
}

func ExampleCredentials_OnDemandSnapshotPhysical() {
	rubrik := rubrikcdm.ConnectEnv()

	hostname := "vm01"
	slaName := "current"
	fileset := "C_Drive"
	hostOS := "Windows"

	hostSnapshot := rubrik.OnDemandSnapshotPhysical(hostname, slaName, fileset, hostOS)
}

func ExampleCredentials_ResumeSnapshot() {
	rubrik := rubrikcdm.ConnectEnv()

	vmName := "vm01"

	resumeVM := rubrik.ResumeSnapshot(vmName, "vmware")
}

func ExampleCredentials_GetSLAObjects() {
	rubrik := rubrikcdm.ConnectEnv()

	slaName := "Gold"

	getObjSLA := rubrik.GetSLAObjects(slaName, "vmware")

}

func ExampleCredentials_EndManagedVolumeSnapshot() {
	rubrik := rubrikcdm.ConnectEnv()

	mvName := "GoSDK"
	slaName := "Gold"

	endMV := rubrik.EndManagedVolumeSnapshot(mvName, slaName)
}

func ExampleCredentials_AssignSLA() {
	rubrik := rubrikcdm.ConnectEnv()

	objectName := "vm01"
	slaName := "Bronze"

	assignSLA := rubrik.AssignSLA(objectName, "vmware", slaName)
}

func ExampleCredentials_ConfigureTimezone() {
	rubrik := rubrikcdm.ConnectEnv()

	confTimezone := rubrik.ConfigureTimezone("America/Los_Angeles")
}

func ExampleCredentials_ConfigureVLAN() {
	rubrik := rubrikcdm.ConnectEnv()

	vlanIPs := map[string]string{}
	vlanIPs["RVM157S018901"] = "192.168.100.100"
	vlanIPs["RVM157S018902"] = "192.168.100.101"
	vlanIPs["RVM157S018903"] = "192.168.100.102"
	vlanIPs["RVM157S0189014"] = "192.168.100.103"

	vlan := 100
	netmask := "255.255.255.0"

	configVLAN := rubrik.ConfigureVLAN(netmask, vlan, vlanIPs)
}

func ExampleCredentials_AddvCenter() {
	rubrik := rubrikcdm.ConnectEnv()

	vCenterIP := "demogosdk.lab"
	vCenterUsername := "go"
	vCenterPassword := "sdk"
	vmLinking := true

	addVcenter := rubrik.AddvCenter(vCenterIP, vCenterUsername, vCenterPassword, vmLinking)
}

func ExampleCredentials_AddvCenterWithCert() {
	rubrik := rubrikcdm.ConnectEnv()

	vCenterIP := "demogosdk.lab"
	vCenterUsername := "go"
	vCenterPassword := "sdk"
	readcaCertificate, _ := ioutil.ReadFile("ca_cert")
	caCert := string(readcaCertificate)
	vmLinking := true

	addVcenterWithCert := rubrik.AddvCenterWithCert(vCenterIP, vCenterUsername, vCenterPassword, caCert, vmLinking)
}

func ExampleCredentials_ConfigureSMTPSettings() {
	rubrik := rubrikcdm.ConnectEnv()

	hostname := "smtp.GOSDK.lab"
	port := 100
	fromEmail := "gosdk@rubrik.com"
	username := "go"
	password := "sdk"
	encryption := "NONE"

	smtpConfig := rubrik.ConfigureSMTPSettings(hostname, fromEmail, username, password, encryption, port)

	confTimezone := rubrik.ConfigureTimezone("America/Los_Angeles")
}

func ExampleCredentials_ConfigureSearchDomain() {
	rubrik := rubrikcdm.ConnectEnv()

	searchDomains := []string{"gosdk.lab"}

	searchDomainConfig := rubrik.ConfigureSearchDomain(searchDomains)

}

func ExampleCredentials_ObjectID() {
	rubrik := rubrikcdm.ConnectEnv()

	slaName := "Gold"

	slaID := c.ObjectID(slaName, "sla")
}

func ExampleCredentials_Bootstrap() {

	bootstrapNode := "10.77.16.239"
	rubrik := rubrikcdm.Connect(bootstrapNode, "", "")

	clusterName := "Go-SDK"
	adminEmail := "gosdk@rubrik.com"
	adminPassword := "RubrikGoSDK"
	managementGateway := "10.77.16.1"
	managementSubnetMask := "255.255.252.0"
	dnsSearchDomain := []string{"gosdk.lab"}
	dnsNameServers := []string{}
	ntpServers := []string{"192.21.10.21", "192.21.10.22"}
	enableEncryption := true // set to false for a Cloud Cluster
	waitForCompletion := true

	nodeConfig := map[string]string{}
	nodeConfig["RVM157S018901"] = bootstrapNode
	nodeConfig["RVM157S018902"] = "10.77.16.56"
	nodeConfig["RVM157S018903"] = "10.77.16.198"
	nodeConfig["RVM157S018904"] = "10.77.16.81"

	bootstrap := rubrik.Bootstrap(clusterName, adminEmail, adminPassword, managementGateway, managementSubnetMask, dnsSearchDomain, dnsNameServers, ntpServers, nodeConfig, enableEncryption, waitForCompletion)

}

func ExampleCredentials_ConfigureDNSServers() {
	rubrik := rubrikcdm.ConnectEnv()

	dnsServers := []string{"192.21.10.50", "192.21.10.51"}

	dnsConfig := rubrik.ConfigureDNSServers(dnsServers)
}

func ExampleCredentials_ConfigureSyslog() {
	rubrik := rubrikcdm.ConnectEnv()

	syslogIP := "192.21.11.29"
	syslogProtocol := "UDP"
	syslogPort := 514

	syslog := rubrik.ConfigureSyslog(syslogIP, syslogProtocol, syslogPort)
}

func ExampleCredentials_ConfigureNTP() {
	rubrik := rubrikcdm.ConnectEnv()

	ntpServers := []string{"192.21.10.21", "192.21.10.22"}

	ntp := rubrik.ConfigureNTP(ntpServers)
}

func ExampleCredentials_ClusterNodeIP() {
	rubrik := rubrikcdm.ConnectEnv()

	clusterVersion := rubrik.ClusterNodeIP()
}

func ExampleCredentials_EndUserAuthorization() {
	rubrik := rubrikcdm.ConnectEnv()

	vmName := "vm01"
	endUser := "user01"

	endUserAuth := rubrik.EndUserAuthorization(vmName, endUser, "vmware")
}

func ExampleCredentials_ClusterVersionCheck() {
	rubrik := rubrikcdm.ConnectEnv()

	clusterVersion := rubrik.ClusterVersionCheck(4.2)
}

func ExampleCredentials_Get() {
	rubrik := rubrikcdm.ConnectEnv()

	clusterInfo := rubrik.Get("v1", "/cluster/me")
}

func ExampleCredentials_Post() {
	rubrik := rubrikcdm.ConnectEnv()

	config := map[string]string{}
	config["slaId"] = "388a473c-3361-42ab-8f5b-08edb76891f6"

	onDemandSnapshot := rubrik.Post("v1", "/vmware/vm/VirtualMachine:::fbcb1f51-9520-4227-a68c-6fe145982f48-vm-204969/snapshot", config)
}

func ExampleCredentials_Patch() {
	rubrik := rubrikcdm.ConnectEnv()

	config := map[string]string{}
	config["configuredSlaDomainId"] = "388a473c-3361-42ab-8f5b-08edb76891f6"

	fileset := rubrik.Patch("v1", "/fileset/Fileset:::b95456e2-7d60-4ed0-af88-648516e139a6", config)
}

func ExampleCredentials_Delete() {
	rubrik := rubrikcdm.ConnectEnv()

	deleteSLA := rubrik.Delete("v1", "/sla_domain/388a473c-3361-42ab-8f5b-08edb76891f6")
}

func ExampleCredentials_AddAWSNativeAccount() {
	rubrik := rubrikcdm.ConnectEnv()

	awsAccountName := "GO SDK Demo" // This is the name that will be displayed in the Rubrik UI
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsRegions := []string{"us-east-1"}

	usEast1 := map[string]string{}
	usEast1["region"] = "us-east-1"
	usEast1["vNetId"] = "vpc-11a44968"
	usEast1["subnetId"] = "subnet-3ac58e06"
	usEast1["securityGroupId"] = "sg-9ba90ee5"
	boltConfig := []interface{}{usEast1}

	addAWSNative := rubrik.AddAWSNativeAccount(awsAccountName, awsAccessKey, awsSecretKey, awsRegions, boltConfig)
}

func ExampleCredentials_AWSS3CloudOutRSA() {
	rubrik := rubrikcdm.ConnectEnv()

	awsBucket := "rubrikgosdk"
	storageClass := "standard"
	archiveName := "AWS:S3:GoSDK"
	awsRegion := "us-east-1"
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	readRSAKey, _ := ioutil.ReadFile("rsa_key.pem")
	rsaKey := string(readRSAKey)

	awsCloudOut := rubrik.AWSS3CloudOut(awsBucket, storageClass, archiveName, awsRegion, awsAccessKey, awsSecretKey, rsaKey)
}

func ExampleCredentials_AWSS3CloudOutKMS() {
	rubrik := rubrikcdm.ConnectEnv()

	awsBucket := "rubrikgosdk"
	storageClass := "standard"
	archiveName := "AWS:S3:GoSDK"
	awsRegion := "us-east-1"
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	kmsMasterKeyID := os.Getenv("AWS_MASTER_KEY_ID")

	awsCloudOut := rubrik.AWSS3CloudOut(awsBucket, storageClass, archiveName, awsRegion, awsAccessKey, awsSecretKey, kmsMasterKeyID)
}

func ExampleCredentials_AWSS3CloudOn() {
	rubrik := rubrikcdm.ConnectEnv()

	archiveName := "AWS:S3:GoSDK"
	vpcID := "vpc-28e32931"
	subnetID := "subnet-3ae87e92"
	securityGroupID := "sg-9ba32ff8"

	awsCloudOn := rubrik.AWSS3CloudOn(archiveName, vpcID, subnetID, securityGroupID)
}

func ExampleCredentials_AzureCloudOut() {
	rubrik := rubrikcdm.ConnectEnv()

	container := "gosdk"
	azureAccessKey := os.Getenv("AZURE_ACCESS_KEY")
	storageAccountName := "rubrikgosdk"
	archiveName := "Azure:gosdk"
	instanceType := "default"
	readRSAKey, _ := ioutil.ReadFile("rsa_key.pem")
	rsaKey := string(readRSAKey)

	azureCloudOut := rubrik.AzureCloudOut(container, azureAccessKey, storageAccountName, archiveName, instanceType, rsaKey)
}

func ExampleCredentials_AzureCloudOn() {
	rubrik := rubrikcdm.ConnectEnv()

	archiveName := "Azure:GoSDK"
	container := "gosdk"
	storageAccountName := "rubrikgosdk"
	applicationID := os.Getenv("AZURE_APP_ID")
	applicationKey := os.Getenv("AZURE_APP_KEY")

	directoryID := os.Getenv("AZURE_DIRECTORTY_ID")
	region := "westus2"
	virtualNetworkID := os.Getenv("AZURE_VNET_ID")
	subnetName := os.Getenv("AZURE_SUBNET")
	securityGroupID := os.Getenv("AZURE_SECURITY_GROUP")

	azureCloudOn := rubrik.AzureCloudOn(archiveName, container, storageAccountName, applicationID, applicationKey, directoryID, region, virtualNetworkID, subnetName, securityGroupID)

}
