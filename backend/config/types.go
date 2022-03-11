package config

type ClusterConfiguration struct {
	Name string `json:"name"`
	Zones []string `json:"zones"`
}

func (cc *ClusterConfiguration) IsValid() bool {
	if cc.Zones == nil { cc.Zones = make([]string, 0) }
	return cc.Name != ""
}

type StorageDeploymentConfiguration struct {
	ClusterName string `json:"cluster_name"`
	Alias string `json:"alias"`
	Endpoint string `json:"endpoint"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	UseSSL bool `json:"use_ssl"`
	ManagementURL string `json:"management_url"`
}

func (sdc *StorageDeploymentConfiguration) IsValid() bool {
	return sdc.ClusterName != "" && sdc.Alias != "" && sdc.Endpoint != "" && sdc.AccessKey != "" && sdc.SecretKey != ""
}

type FaaSDeploymentConfiguration struct {
	ClusterName string `json:"cluster_name"`
	URL string `json:"url"`
}

func (fdc *FaaSDeploymentConfiguration) IsValid() bool {
	return fdc.ClusterName != "" && fdc.URL != ""
}

type BucketConfiguration struct {
	Name string `json:"name"`
	StorageDeploymentAlias string `json:"storage_deployment_alias"`
	AllowedZones []string `json:"allowed_zones"`
	TargetReplicaCount int `json:"target_replica_count"`
}

func (bc *BucketConfiguration) IsValid() bool {
	if bc.AllowedZones == nil { bc.AllowedZones = make([]string, 0) }
	return bc.Name != "" && bc.StorageDeploymentAlias != ""
}

type ServerConfiguration struct {
	Clusters []ClusterConfiguration `json:"clusters"`
	StorageDeployments []StorageDeploymentConfiguration `json:"storage_deployments"`
	FaaSDeployments []FaaSDeploymentConfiguration `json:"faas_deployments"`
	Buckets []BucketConfiguration `json:"buckets"`
}
