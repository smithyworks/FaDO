package database

import (
	"github.com/jackc/pgx/v4"
	"github.com/smithyworks/FaDO/util"
)

// Table replica_bucket_locations

type ReplicaBucketLocationRecord struct {
	BucketID int64 `json:"bucket_id"`
	StorageID int64 `json:"storage_id"`
}

func ScanReplicaBucketLocationRows(rows pgx.Rows) (replicaBucketsStorageDeployments []ReplicaBucketLocationRecord, err error) {
	for rows.Next() {
		var rbsdr ReplicaBucketLocationRecord

		err = rows.Scan(
			&rbsdr.BucketID,
			&rbsdr.StorageID,
		)
		if err != nil { return replicaBucketsStorageDeployments, util.ProcessErr(err) }

		replicaBucketsStorageDeployments = append(replicaBucketsStorageDeployments, rbsdr)
	}

	return
}

func QueryReplicaBucketLocations(conn DBConn, sql string, args ...interface{}) (replicaBucketLocations []ReplicaBucketLocationRecord, err error) {
	rows, err := Query(conn, sql, args...)
	if err != nil { return replicaBucketLocations, util.ProcessErr(err) }
	defer rows.Close()
	
	replicaBucketLocations, err = ScanReplicaBucketLocationRows(rows)
	if err != nil { return replicaBucketLocations, util.ProcessErr(err) }

	return
}
