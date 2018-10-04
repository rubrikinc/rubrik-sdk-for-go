# rubrik-sdk-for-go
Rubrik Cloud Data Management SDK for Go

```go
package main

import (
	"client"
	"fmt"
)

func main() {

	rubrik := client.Connect{
		NodeIP:   "172.21.8.53",
		Username: "yourUsername",
		Password: "yourPassword",
	}
	// GET
	clusterVersion := rubrik.Get("v1", "/cluster/me")
	fmt.Println(clusterVersion["version"])

	// POST
	mapConfig := map[string]string{"slaId": "388a473c-3361-42ab-8f5b-08edb76891f6"}
	fmt.Println(mapConfig)

	config, _ := json.Marshal(mapConfig)
	fmt.Println(config)

	snapshot := rubrik.Post("v1", "/vmware/vm/VirtualMachine:::fbcb1f51-9520-4227-a68c-6fe145982f48-vm-204969/snapshot", config)
	fmt.Println(snapshot)

	// PATCH
	mapConfig := map[string]string{"configuredSlaDomainId": "388a473c-3361-42ab-8f5b-08edb76891f6"}
	fmt.Println(mapConfig)

	config, _ := json.Marshal(mapConfig)
	fmt.Println(config)

	snapshot := rubrik.Patch("v1", "/fileset/Fileset:::b95456e2-7d60-4ed0-af88-648516e139a6", config)
	fmt.Println(snapshot)

	// DELETE
	deleteSLA := rubrik.Delete("v1", "/sla_domain/00f81582-4b06-40a6-9583-7d69db5864a6", 60)
	fmt.Println(deleteSLA)

}
```
