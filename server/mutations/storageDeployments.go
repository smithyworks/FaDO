package mutations

import (
	"log"
	"strings"
	"time"

	"github.com/minio/madmin-go"
	"github.com/smithyworks/FaDO/database"
	"github.com/smithyworks/FaDO/mc"
	"github.com/smithyworks/FaDO/util"
)

func AddStorageDeployment(conn database.DBConn, storageDeployment database.StorageDeploymentRecord) (err error) {
	shouldUpdate := false
	if sdRecords, err := database.QueryStorageDeployments(conn, "SELECT * FROM storage_deployments WHERE endpoint = $1", storageDeployment.Endpoint); err != nil {
		return util.ProcessErr(err)
	} else if len(sdRecords) == 1 {
		storageDeployment = sdRecords[0]
		shouldUpdate = true
	}

	// get info from MinIO
	if err = GetStorageDeploymentInfo(conn, &storageDeployment); err != nil {
		return util.ProcessErr(err)
	}

	// Insert/update to database
	if shouldUpdate {
		if _, err = database.Exec(conn, "UPDATE storage_deployments SET sqs_arn = $1, minio_deployment_id = $2 WHERE storage_id = $3", storageDeployment.SqsArn, storageDeployment.MinioDeploymentID, storageDeployment.StorageID); err != nil {
			return util.ProcessErr(err)
		}
	} else {
		if storageDeployment, err = database.InsertStorageDeployment(conn, storageDeployment); err != nil {
			return util.ProcessErr(err)
		}
	}

	// Create a map of current buckets
	bucketMap := make(map[string]database.BucketRecord)
	if buckets, err := database.QueryBuckets(conn, "SELECT * FROM buckets"); err != nil {
		return util.ProcessErr(err)
	} else {
		for _, b := range buckets { bucketMap[b.Name] = b }
	}

	// Scan buckets from new deployments
	if client, err := CreateMinioClient(conn, storageDeployment); err != nil {
		return util.ProcessErr(err)
	} else {
		if minioBuckets, err := client.ListBuckets(ctx); err != nil {
			return util.ProcessErr(err)
		} else {
			for _, mb := range minioBuckets {
				// AddMastBucket / AddReplicaBucket
				if br, exists := bucketMap[mb.Name]; exists && br.StorageID != storageDeployment.StorageID {
					if err = AddReplicaBucket(conn, database.ReplicaBucketLocationRecord{BucketID: br.BucketID, StorageID: storageDeployment.StorageID}); err != nil {
						return util.ProcessErr(err)
					}
				} else if !exists {
					if err = AddMasterBucket(conn, database.BucketRecord{StorageID: storageDeployment.StorageID, Name: mb.Name}, 0, []string{}, []int64{}); err != nil {
						return util.ProcessErr(err)
					}
				}
			}
		}
	}

	return
}

func DeleteStorageDeployment(conn database.DBConn, storageDeployement database.StorageDeploymentRecord, permanent bool) (err error) {
	// Delete Master Buckets and replicas if permanent.
	if permanent {
		if masterBuckets, err := database.QueryBuckets(conn, "SELECT * FROM buckets WHERE storage_id = $1", storageDeployement.StorageID); err != nil {
			return util.ProcessErr(err)
		} else {
			for _, b := range masterBuckets {
				if err = EnsureBucketDeletion(conn, storageDeployement.StorageID, b.Name); err != nil {
					util.PrintWarning(err)
					err = nil
				}
			}
		}

		if replicaLocations, err := database.QueryReplicaBucketLocations(conn, "SELECT * FROM replica_bucket_locations WHERE storage_id = $1", storageDeployement.StorageID); err != nil {
			return util.ProcessErr(err)
		} else {
			for _, rl := range replicaLocations {
				if bucket, err := database.QueryBucketRow(conn, "SELECT * FROM buckets WHERE bucket_id = $1", rl.BucketID); err != nil {
					return util.ProcessErr(err)
				} else {
					if err = EnsureBucketDeletion(conn, storageDeployement.StorageID, bucket.Name); err != nil {
						util.PrintWarning(err)
						err = nil
					}
				}
			}
		}
	}
	
	// Delete Replica Buckets.
	var bucketIdsToResolve []int64
	if rows, err := database.Query(conn, "SELECT * FROM replica_bucket_locations WHERE storage_id = $1", storageDeployement.StorageID); err != nil {
		return util.ProcessErr(err)
	} else {
		if bucketStorageRecords, err := database.ScanReplicaBucketLocationRows(rows); err != nil {
			return util.ProcessErr(err)
		} else {
			for _, bs := range bucketStorageRecords {
				bucketIdsToResolve = append(bucketIdsToResolve, bs.BucketID)
				if err = DeleteReplicaBucket(conn, bs); err != nil {
					return util.ProcessErr(err)
				}
			}
		}
	}

	// Delete from database
	if _, err = database.Exec(conn, "DELETE FROM storage_deployments WHERE storage_id = $1", storageDeployement.StorageID); err != nil {
		return util.ProcessErr(err)
	}

	// ResolveBucket Replicas for delete replica masters
	if bucketsToResolve, err := database.QueryBuckets(conn, "SELECT * FROM buckets WHERE bucket_id = ANY($1::int[])", bucketIdsToResolve); err != nil {
		return util.ProcessErr(err)
	} else {
		for _, b := range bucketsToResolve {
			if err = ResolveBucketReplicas(conn, b); err != nil {
				return util.ProcessErr(err)
			}
		}
	}

	return
}

func GetStorageDeploymentInfo(conn database.DBConn, sd *database.StorageDeploymentRecord) (err error) {
	if err = mc.SetAlias(*sd); err != nil {
		return util.ProcessErr(err)
	}

	if adminClient, err := CreateMinioAdminClient(conn, sd); err != nil {
		return util.ProcessErr(err)
	} else {
		if err = mc.SetNotificationTarget(*sd); err != nil {
			return util.ProcessErr(err)
		}

		if err = adminClient.ServiceRestart(ctx); err != nil {
			return util.ProcessErr(err)
		}

		d1 := time.Now()
		log.Printf("INFO: Restarting minio deployment %v.", sd.Alias)
		var info madmin.InfoMessage
		for i := 0; i < 3; i++ {
			time.Sleep(5 * time.Second)
			if info, err = adminClient.ServerInfo(ctx); err == nil {
				break
			}
		}
		if err != nil { return util.ProcessErr(err) }
		d2 := time.Since(d1)
		log.Printf("INFO: Service restarted after %v seconds.", d2.Seconds())

		sd.MinioDeploymentID = info.DeploymentID
		for _, val := range info.SQSARN {
			if strings.HasSuffix(val, "fado:webhook") {
				sd.SqsArn = val
				break
			}
		}
	}

	return util.ProcessErr(err)
}
