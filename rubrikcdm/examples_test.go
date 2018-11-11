package rubrikcdm_test

import (
	"github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm"
)

func ExampleCredentials_ClusterVersion() {
	rubrik := rubrikcdm.ConnectEnv()

	clusterVersion := rubrik.ClusterVersion()
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
