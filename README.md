# rubrik-sdk-for-go
Rubrik Cloud Data Management SDK for Go

```go
package main

import (
	"fmt"
	"rubrik-sdk-for-go/src"
)

func main() {

	rubrik := rubrikcdm.ConnectEnv()
	
	// GET
	clusterVersion := rubrik.Get("v1", "/cluster/me")
	fmt.Println(clusterVersion["version"])

	// POST
	config := map[string]string{"slaId": "388a473c-3361-42ab-8f5b-08edb76891f6"}

	onDemandSnapshot := rubrik.Post("v1", "/vmware/vm/VirtualMachine:::fbcb1f51-9520-4227-a68c-6fe145982f48-vm-204969/snapshot", config)
	fmt.Println(onDemandSnapshot)

	// PATCH
	config := map[string]string{"configuredSlaDomainId": "388a473c-3361-42ab-8f5b-08edb76891f6"}

	fileset := rubrik.Patch("v1", "/fileset/Fileset:::b95456e2-7d60-4ed0-af88-648516e139a6", config)
	fmt.Println(fileset)

	// DELETE
	deleteSLA := rubrik.Delete("v1", "/sla_domain/1e76a57b-f96b-483d-91cd-65d0a43a1eed")
	fmt.Println(deleteSLA)

}
```
