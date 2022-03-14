package mutations

import (
	"github.com/minio/minio-go/v7"
	"github.com/smithyworks/FaDO/database"
	"github.com/smithyworks/FaDO/util"
)

func TrackBucketObjects(conn database.DBConn, bucket database.BucketRecord) (err error) {
	// Fetch currently tracked objects and make a convenient map
	databaseObjectMap := make(map[string]database.ObjectRecord)
	if objectRecords, err := database.QueryObjects(conn, "SELECT * FROM objects WHERE bucket_id = $1", bucket.BucketID); err != nil {
		return util.ProcessErr(err)
	} else {
		for _, or := range objectRecords { databaseObjectMap[or.Name] = or }
	}

	// List out objects from minio
	latestObjects := make([]string, 0)
	if client, err := CreateMinioClient(conn, bucket.StorageID); err != nil {
		return util.ProcessErr(err)
	} else {
		for o := range client.ListObjects(ctx, bucket.Name, minio.ListObjectsOptions{Recursive: true}) {
			latestObjects = append(latestObjects, o.Key)
		}
	}

	// Go through all the minio objects and make sure they are tracked in the database,
	for _, oName := range latestObjects {
		_, oExists := databaseObjectMap[oName]
		if oExists {
			delete(databaseObjectMap, oName)
		} else {
			newO := database.ObjectRecord{BucketID: bucket.BucketID, Name: oName}
			if _, err = database.InsertObject(conn, newO); err != nil {
				return util.ProcessErr(err)
			}
		}
	}

	// Clean up any objects that are no longer in minio
	for _, obj := range databaseObjectMap {
		if _, err = database.Exec(conn, "DELETE FROM objects WHERE object_id = $1", obj.ObjectID); err != nil {
			return util.ProcessErr(err)
		}
	}

	return
}

func DeleteObject(conn database.DBConn, object database.ObjectRecord) (err error) {
	bucket, err := database.QueryBucketRow(conn, "SELECT * FROM buckets WHERE bucket_id = $1", object.BucketID)
	if err != nil { return util.ProcessErr(err) }

	client, err := CreateMinioClient(conn, bucket.StorageID)
	if err != nil { return util.ProcessErr(err) }

	err = client.RemoveObject(ctx, bucket.Name, object.Name, minio.RemoveObjectOptions{ForceDelete: true})
	if err != nil { return util.ProcessErr(err) }

	_, err = database.Exec(conn, "DELETE FROM objects WHERE object_id = $1", object.ObjectID)
	if err != nil { return util.ProcessErr(err) }

	return
}