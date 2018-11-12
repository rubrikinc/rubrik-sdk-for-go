# rubrik-sdk-for-go
Rubrik Cloud Data Management SDK for Go

<p></p>
<p align="center">
  <img src="https://user-images.githubusercontent.com/8610203/48332236-55506f00-e610-11e8-9a60-594de963a1ee.png" alt="Rubrik Gopher Logo"/>
</p>

# Installation

```go get github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm```

# Documentation

[https://godoc.org/github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm](https://godoc.org/github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm)

# Example

```go
package main

import (
	"fmt"
	"github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm"
)

func main() {

	rubrik := rubrikcdm.ConnectEnv()
	
	// GET the Rubrik cluster Version
	clusterSummary := rubrik.Get("v1", "/cluster/me")
	fmt.Println(clusterSummary["version"])

	// Simplified Function to determine the Rubrik cluster version
	clusterVersion := rubrik.ClusterVersion()
	fmt.Println(clusterVersion)

}
```

# AUTHOR INFORMATION

<p></p>
<p align="center">
  <img src="https://user-images.githubusercontent.com/8610203/37415009-6f9cf416-2778-11e8-8b56-052a8e41c3c8.png" alt="Rubrik Ranger Logo"/>
</p>

