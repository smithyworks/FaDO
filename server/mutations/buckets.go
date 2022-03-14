package mutations

import (
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/notification"
	"github.com/smithyworks/FaDO/database"
	"github.com/smithyworks/FaDO/mc"
	"github.com/smithyworks/FaDO/util"
)

func AddMasterBucket(conn database.DBConn, bucket database.BucketRecord, targetReplicaCount int, allowedZones []string, replicaStorageDeployments []int64) (err error) {
	// Check for existence of bucket
	if bucketRecords, err := database.QueryBuckets(conn, "SELECT * FROM buckets WHERE name = $1", bucket.Name); err != nil {
		return util.ProcessErr(err)
	} else {
		if len(bucketRecords) == 0 {
			// Insert bucket into database
			if bucket, err = database.InsertBucket(conn, bucket); err != nil {
				return util.ProcessErr(err)
			}
			if allowedZones != nil {
				err = database.SetBucketPolicy(conn, bucket, "zones", allowedZones)
				if err != nil { return util.ProcessErr(err) }
			}
			if targetReplicaCount != 0 {
				err = database.SetBucketPolicy(conn, bucket, "target_replica_count", targetReplicaCount)
				if err != nil { return util.ProcessErr(err) }
			}
			if replicaStorageDeployments != nil && len(replicaStorageDeployments) > 0 {
				err = database.SetBucketPolicy(conn, bucket, "replica_locations", replicaStorageDeployments)
			}
		} else {
			bucket = bucketRecords[0]
		}
	}
	
	// Ensure bucket exists in minio
	if err = EnsureBucketCreation(conn, bucket.StorageID, bucket.Name); err != nil {
		return util.ProcessErr(err)
	}

	// Set up notifications
	if err = SetupBucketNotifications(conn, bucket); err != nil {
		return util.ProcessErr(err)
	}

	// Update LB config
	if err := ConfigureLoadBalancer(conn); err != nil {
		util.PrintWarning(err)
	}

	if err = TrackBucketObjects(conn, bucket); err != nil {
		return util.ProcessErr(err)
	}

	if err = SetBucketReplicas(conn, bucket, replicaStorageDeployments); err != nil {
		return util.ProcessErr(err)
	}

	if err = ResolveBucketReplicas(conn, bucket); err != nil {
		return util.ProcessErr(err)
	}

	return
}

func SetBucketReplicas(conn database.DBConn, bucket database.BucketRecord, replicaStorageDeployments []int64) (err error) {
	rows, err := database.Query(conn, "SELECT storage_ids FROM existing_bucket_locations WHERE bucket_id = $1", bucket.BucketID)
    if err != nil { util.PrintErr(err); return }
	defer rows.Close()

	existingStorageIds := make([]int64, 0)
	if rows.Next() {
		rows.Scan(&existingStorageIds)
		rows.Close()
	}

	for _, sid := range replicaStorageDeployments {
		if !util.HasInt(existingStorageIds, sid) {
			err = AddReplicaBucket(conn, database.ReplicaBucketLocationRecord{BucketID: bucket.BucketID, StorageID: sid})
			if err != nil { return util.ProcessErr(err) }
		}
	}
	for _, sid := range existingStorageIds {
		if !util.HasInt(replicaStorageDeployments, sid) {
			err = DeleteReplicaBucket(conn, database.ReplicaBucketLocationRecord{BucketID: bucket.BucketID, StorageID: sid})
			if err != nil { return util.ProcessErr(err) }
		}
	}

	return
}

func ResolveBucketReplicas(conn database.DBConn, bucket database.BucketRecord) (err error) {
	var locations []int64
	if ok, err := database.GetBucketPolicy(conn, bucket, "replica_locations", &locations); ok && err == nil {
		err = SetBucketReplicas(conn, bucket, locations)
		if err != nil { return util.ProcessErr(err) } else { return nil }
	} // nothing to do

	var bucketZones []string
	_, err = database.GetBucketPolicy(conn, bucket, "zones", &bucketZones)
	if err != nil { return util.ProcessErr(err) }

	// get bucket's home cluster for exclusion

	homeStorageDeployment, err := database.QueryStorageDeploymentRow(conn, "SELECT * FROM storage_deployments WHERE storage_id = $1", bucket.StorageID)
	if err != nil { return util.ProcessErr(err) }
	homeCluster, err := database.QueryClusterRow(conn, "SELECT * FROM clusters WHERE cluster_id = $1", homeStorageDeployment.ClusterID)
	if err != nil { return util.ProcessErr(err) }

	clusters, err := database.QueryClusters(conn, "SELECT * FROM clusters")
	if err != nil { return util.ProcessErr(err) }

	var possibleCluster []int64
	for _, c := range clusters {
		var clusterZones []string
		_, err = database.GetClusterPolicy(conn, c, "zones", &clusterZones)
		if err != nil { return util.ProcessErr(err) }

		overlap := false
		if homeCluster.ClusterID != c.ClusterID {
			if len(bucketZones) == 0 { overlap = true } else {
				for _, z := range clusterZones { if util.HasString(bucketZones, z) { overlap = true; break } }
			}
		}

		if overlap { possibleCluster = append(possibleCluster, c.ClusterID) }
	}

	var targetReplicaCount int
	_, err = database.GetBucketPolicy(conn, bucket, "target_replica_count", &targetReplicaCount)
	if err != nil { return util.ProcessErr(err) }

	storageDeployments, err := database.QueryStorageDeployments(conn, "SELECT * FROM storage_deployments WHERE cluster_id = ANY($1) LIMIT $2", possibleCluster, targetReplicaCount)
	if err != nil { return util.ProcessErr(err) }

	var storageIDs []int64
	for _, sd := range storageDeployments { storageIDs = append(storageIDs, sd.StorageID) }

	err = SetBucketReplicas(conn, bucket, storageIDs)
	if err != nil { return util.ProcessErr(err) }
	
	return
}

func AddReplicaBucket(conn database.DBConn, bucketStorage database.ReplicaBucketLocationRecord) (err error) {
	// check for existence, insert into database
	if replicaBucketsStorageDeployments, err := database.QueryReplicaBucketLocations(conn, "SELECT * FROM replica_bucket_locations WHERE bucket_id = $1 AND storage_id = $2", bucketStorage.BucketID, bucketStorage.StorageID); err != nil {
		return util.ProcessErr(err)
	} else if len(replicaBucketsStorageDeployments) == 0 {
		if _, err = database.Exec(conn, "INSERT INTO replica_bucket_locations (bucket_id, storage_id) VALUES ($1, $2)", bucketStorage.BucketID, bucketStorage.StorageID); err != nil {
			return util.ProcessErr(err)
		}
	}

	var bucketName string
	if rows, err := database.Query(conn, "SELECT name FROM buckets WHERE bucket_id = $1", bucketStorage.BucketID); err != nil {
		return util.ProcessErr(err)
	} else {
		rows.Next(); err = rows.Scan(&bucketName); rows.Close()
		if err != nil { return util.ProcessErr(err) }
	}

	// create bucket in minio
	if err = EnsureBucketCreation(conn, bucketStorage.StorageID, bucketName); err != nil {
		return util.ProcessErr(err)
	}

	// mirror from master
	if rows, err := database.Query(conn, "SELECT * FROM bucket_replications WHERE bucket_id = $1 AND dst_storage_id = $2", bucketStorage.BucketID, bucketStorage.StorageID); err != nil {
		return util.ProcessErr(err)
	} else {
		if bucketReplications, err := database.ScanBucketReplicationRows(rows); err != nil {
			return util.ProcessErr(err)
		} else if len(bucketReplications) == 1 {
			br := bucketReplications[0]
			if err = mc.Mirror(br.SrcStorageAlias, br.BucketName, br.DstStorageAlias, br.BucketName); err != nil {
				return util.ProcessErr(err)
			}
		}
	}

	// update LB config
	if err := ConfigureLoadBalancer(conn); err != nil {
		util.PrintWarning(err)
	}
	
	return
}

func EditMasterBucket(conn database.DBConn, bucket database.BucketRecord, targetReplicaCount int, zones []string, replicaStorageDeployments []int64) (err error) {
	err = database.SetBucketPolicy(conn, bucket, "target_replica_count", targetReplicaCount)
	if err != nil { return util.ProcessErr(err) }
	if zones != nil {
		err = database.SetBucketPolicy(conn, bucket, "zones", zones)
		if err != nil { return util.ProcessErr(err) }
	}

	if replicaStorageDeployments == nil || len(replicaStorageDeployments) == 0 {
		err = database.DeleteBucketPolicy(conn, bucket, "replica_locations")
		if err != nil { return util.ProcessErr(err) }
	} else {
		err = database.SetBucketPolicy(conn, bucket, "replica_locations", replicaStorageDeployments)
		if err != nil { return util.ProcessErr(err) }
	}

	// ResolveBucketReplicas
	if err = ResolveBucketReplicas(conn, bucket); err != nil {
		return util.ProcessErr(err)
	}

	return
}

func DeleteMasterBucket(conn database.DBConn, bucket database.BucketRecord) (err error) {
	// bucket replica / storage deployment associations
	var bucketStorageRecords []database.ReplicaBucketLocationRecord
	if bucketStorageRecords, err = database.QueryReplicaBucketLocations(conn, "SELECT * FROM replica_bucket_locations WHERE bucket_id = $1", bucket.BucketID); err != nil {
		return util.ProcessErr(err)
	}


	// for replicas
	//   DeleteReplicaBucket
	for _, bucketStorage := range bucketStorageRecords {
		if err = DeleteReplicaBucket(conn, bucketStorage); err != nil {
			return util.ProcessErr(err)
		}
	}

	// delete from MinIO
	if err = EnsureBucketDeletion(conn, bucket.StorageID, bucket.Name); err != nil {
		util.PrintWarning(err)
		err = nil
	}

	// delete from database
	if _, err = database.Exec(conn, "DELETE FROM buckets WHERE bucket_id = $1", bucket.BucketID); err != nil {
		return util.ProcessErr(err)
	}

	// update LB config
	if err := ConfigureLoadBalancer(conn); err != nil {
		util.PrintWarning(err)
	}

	return util.ProcessErr(err)
}

func DeleteReplicaBucket(conn database.DBConn, bucketStorage database.ReplicaBucketLocationRecord) (err error) {
	// Fetch bucket for which we want to delete the replica
	bucket, err := database.QueryBucketRow(conn, "SELECT * FROM buckets WHERE bucket_id = $1", bucketStorage.BucketID)
	if err != nil { return util.ProcessErr(err) }

	// delete from database
	if _, err := database.Exec(conn, "DELETE FROM replica_bucket_locations WHERE bucket_id = $1 AND storage_id = $2", bucketStorage.BucketID, bucketStorage.StorageID); err != nil {
		return util.ProcessErr(err)
	}

	// delete from MinIO
	if err = EnsureBucketDeletion(conn, bucketStorage.StorageID, bucket.Name); err != nil {
		util.PrintWarning(err)
		err = nil
	}

	// update LB config
	if err := ConfigureLoadBalancer(conn); err != nil {
		util.PrintWarning(err)
	}

	return util.ProcessErr(err)
}

func EnsureBucketCreation(conn database.DBConn, storage_id int64, bucketName string) (err error) {
	// Ensure bucket exists in minio
	if client, err := CreateMinioClient(conn, storage_id); err != nil {
		return util.ProcessErr(err)
	} else {
		if exists, err := client.BucketExists(ctx, bucketName); err != nil {
			return util.ProcessErr(err)
		} else if !exists {
			if err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
				return util.ProcessErr(err)
			}
		}
	}

	return
}

func EnsureBucketDeletion(conn database.DBConn, storage_id int64, bucketName string) (err error) {
		// Ensure bucket is deleted in minio
		if client, err := CreateMinioClient(conn, storage_id); err != nil {
			return util.ProcessErr(err)
		} else {
			if exists, err := client.BucketExists(ctx, bucketName); err != nil {
				return util.ProcessErr(err)
			} else if exists {
				if err = client.RemoveBucketWithOptions(ctx, bucketName, minio.BucketOptions{ForceDelete: true}); err != nil {
					return util.ProcessErr(err)
				}
			}
		}
	
		return
}

func SetupBucketNotifications(conn database.DBConn, bucket database.BucketRecord) (err error) {
	storageDeployment, err := database.QueryStorageDeploymentRow(conn, "SELECT * FROM storage_deployments WHERE storage_id = $1", bucket.StorageID)
	if err != nil { return util.ProcessErr(err) }

	if client, err := CreateMinioClient(conn, storageDeployment); err != nil {
		return util.ProcessErr(err)
	} else {
		arnTokens := strings.Split(storageDeployment.SqsArn, ":")
		queueArn := notification.NewArn(arnTokens[1], arnTokens[2], arnTokens[3], arnTokens[4], arnTokens[5])
		queueConfig := notification.NewConfig(queueArn)
		queueConfig.AddEvents(notification.ObjectCreatedAll, notification.ObjectRemovedAll)
	
		config := notification.Configuration{}
		config.AddQueue(queueConfig)
	
		if err = client.SetBucketNotification(ctx, bucket.Name, config); err != nil {
			return util.ProcessErr(err)
		}
	}

	return util.ProcessErr(err)
}

func ReplicateBucket(conn database.DBConn, bucketName string) (err error) {
	var bucketReplications []database.BucketReplicationRecord
	if rows, err := database.Query(conn, "SELECT * FROM bucket_replications WHERE bucket_name = $1", bucketName); err != nil {
		return util.ProcessErr(err)
	} else {
		if bucketReplications, err = database.ScanBucketReplicationRows(rows); err != nil {
			return util.ProcessErr(err)
		}
	}

	if len(bucketReplications) < 1 { return }

	for _, brr := range bucketReplications {
		if err = mc.Mirror(brr.SrcStorageAlias, brr.BucketName, brr.DstStorageAlias, brr.BucketName); err != nil {
			return util.ProcessErr(err)
		}
	}

	return
}

