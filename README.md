# Rubrik SDK for Go

<p></p>
<p align="center">
  <img src="https://user-images.githubusercontent.com/8610203/48332236-55506f00-e610-11e8-9a60-594de963a1ee.png" alt="Rubrik Gopher Logo"/>
</p>

# :hammer: Installation

```go get github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm```

# :mag: Example

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

# :blue_book: Documentation

Here are some resources to get you started! If you find any challenges from this project are not properly documented or are unclear, please [raise an issue](https://github.com/rubrikinc/rubrik-sdk-for-go/issues/new/choose) and let us know! This is a fun, safe environment - don't worry if you're a GitHub newbie! :heart:

* [Quick Start Guide](https://github.com/rubrikinc/rubrik-sdk-for-go/blob/master/docs/quick-start.md)
* [SDK for Go Documentation](https://godoc.org/github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm)
* [Rubrik API Documentation](https://github.com/rubrikinc/api-documentation)
* [VIDEO: Getting Started with the Rubrik SDK for Go](https://youtu.be/xklHJe0u-ZY)
* [BLOG: Introducing the Rubrik SDK for Go](https://www.rubrik.com/blog/rubrik-google-go-sdk/)

# :muscle: How You Can Help

We glady welcome contributions from the community. From updating the documentation to adding additional functions, all ideas are welcome. Thank you in advance for all of your issues, pull requests, and comments! :star:

* [Contributing Guide](CONTRIBUTING.md)
* [Code of Conduct](CODE_OF_CONDUCT.md)

# :pushpin: License

* [MIT License](LICENSE)

# :point_right: About Rubrik Build

We encourage all contributors to become members. We aim to grow an active, healthy community of contributors, reviewers, and code owners. Learn more in our [Welcome to the Rubrik Build Community](https://github.com/rubrikinc/welcome-to-rubrik-build) page.

We'd  love to hear from you! Email us: build@rubrik.com :love_letter:
