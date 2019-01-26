# Rubrik SDK for Go

<p></p>
<p align="center">
  <img src="https://user-images.githubusercontent.com/8610203/48332236-55506f00-e610-11e8-9a60-594de963a1ee.png" alt="Rubrik Gopher Logo"/>
</p>

## Installation

```go get github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm```

## Quick Start

[Quick Start Guide](https://github.com/rubrikinc/rubrik-sdk-for-go/blob/master/docs/quick-start.md)

## Documentation

[SDK Documentation](https://godoc.org/github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm)

[Rubrik API Documentation](https://github.com/rubrikinc/api-documentation)

## Example

```go
package main

import (
	"fmt"
        "log"
	
	"github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm"
)

func main() {

	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}
	
	// GET the Rubrik cluster Version
	clusterSummary, err := rubrik.Get("v1", "/cluster/me")
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println(clusterSummary.(map[string]interface{})["version"])

	// Simplified Function to determine the Rubrik cluster version
	clusterVersion, err := rubrik.ClusterVersion()
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println(clusterVersion)

}
```




