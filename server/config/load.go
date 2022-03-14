package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/smithyworks/FaDO/database"
	"github.com/smithyworks/FaDO/mutations"
	"github.com/smithyworks/FaDO/util"
)

var ctx = context.Background()

func ReadConfigurationFile(filePath string) (pc ServerConfiguration, err error) {
	data, err := os.ReadFile(filePath)
	if err != nil { return pc, util.ProcessErr(err) }

	err = json.Unmarshal(data, &pc)

	return
}

func LoadConfigFromFile(configFilePath string) (err error) {
	tx, err := database.Begin()
	if err != nil { return util.ProcessErr(err) }
	defer tx.Rollback(ctx)

	pc, err := ReadConfigurationFile(configFilePath)
	if err != nil { return util.ProcessErr(err) }

	if len(pc.Clusters) < 1 { return util.ProcessErr(fmt.Errorf("Configuration invalid, no clusters were defined.")) }

	for _, c := range pc.Clusters {
		if !c.IsValid() { return util.ProcessErr(fmt.Errorf("Cluster configuration is invalid. %+v", c)) }

		if err = mutations.AddCluster(tx, database.ClusterRecord{Name: c.Name}, c.Zones); err != nil {
			return util.ProcessErr(err)
		}
	}

	for _, f := range pc.FaaSDeployments {
		if !f.IsValid() { return util.ProcessErr(fmt.Errorf("FaaS deployment configuration is invalid. %+v", f)) }

		newFaaSRecord := database.FaaSDeploymentRecord{URL: f.URL}
		if cluster, err := database.QueryClusterRow(tx, "SELECT * FROM clusters WHERE name = $1", f.ClusterName); err != nil {
			return util.ProcessErr(err)
		} else {
			newFaaSRecord.ClusterID = cluster.ClusterID
		}

		if err = mutations.AddFaaSDeployment(tx, newFaaSRecord); err != nil {
			return util.ProcessErr(err)
		}
	}

	for _, s := range pc.StorageDeployments {
		if !s.IsValid() { return util.ProcessErr(fmt.Errorf("Storage Deployment configuration is invalid. %+v", s)) }

		newStorageRecord := database.StorageDeploymentRecord{Alias: s.Alias, Endpoint: s.Endpoint, AccessKey: s.AccessKey, SecretKey: s.SecretKey, UseSSL: s.UseSSL, ManagementURL: s.ManagementURL}
		if cluster, err := database.QueryClusterRow(tx, "SELECT * FROM clusters WHERE name = $1", s.ClusterName); err != nil {
			return util.ProcessErr(err)
		} else {
			newStorageRecord.ClusterID = cluster.ClusterID
		}

		if err = mutations.AddStorageDeployment(tx, newStorageRecord); err != nil {
			return util.ProcessErr(err)
		}
	}

	for _, b := range pc.Buckets {
		if !b.IsValid() { return util.ProcessErr(fmt.Errorf("Bucket configuration is invalid. %+v", b)) }

		newBucketRecord := database.BucketRecord{Name: b.Name}
		if storageDeployement, err := database.QueryStorageDeploymentRow(tx, "SELECT * FROM storage_deployments WHERE alias = $1", b.StorageDeploymentAlias); err != nil {
			return util.ProcessErr(err)
		} else {
			newBucketRecord.StorageID = storageDeployement.StorageID
		}

		if err = mutations.AddMasterBucket(tx, newBucketRecord, b.TargetReplicaCount, b.AllowedZones, []int64{}); err != nil {
			return util.ProcessErr(err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return util.ProcessErr(err)
	}

	return
}