package rubrikcdm

func (c *Credentials) ClusterVersion() interface{} {
	cluster := c.Get("v1", "/cluster/me")
	return cluster["version"]

}
