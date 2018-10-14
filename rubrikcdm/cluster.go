package rubrikcdm

// ClusterVersion returns the CDM version of the Rubrik cluster
func (c *Credentials) ClusterVersion() interface{} {
	apiRequest := c.Get("v1", "/cluster/me")
	return apiRequest["version"]
}

// ClusterNodeIP returns a slice of all Node IPs in the Rubrik cluster
func (c *Credentials) ClusterNodeIP() []string {
	apiRequest := c.Get("internal", "/cluster/me/node")

	var nodeList []string

	for _, v := range apiRequest["data"].([]interface{}) {
		nodeList = append(nodeList, v.(interface{}).(map[string]interface{})["ipAddress"].(string))
	}

	return nodeList
}

// ClusterNodeName returns the name of all nodes in the Rubrik cluster
func (c *Credentials) ClusterNodeName() []string {
	apiRequest := c.Get("internal", "/cluster/me/node")

	var nodeName []string

	for _, v := range apiRequest["data"].([]interface{}) {
		nodeName = append(nodeName, v.(interface{}).(map[string]interface{})["id"].(string))
	}

	return nodeName
}
