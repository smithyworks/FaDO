package database

import (
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/smithyworks/FaDO/util"
)

type BucketReplicationRecord struct {
	BucketID int64 `json:"bucket_id"`
	BucketName string `json:"bucket_name"`
	SrcStorageID int64 `json:"src_storage_id"`
	SrcStorageAlias string `json:"src_storage_alias"`
	SrcMinioDeploymentID string `json:"src_minio_deployment_id"`
	DstStorageID int64 `json:"dst_storage_id"`
	DstStorageAlias string `json:"dst_storage_alias"`
	DstMinioDeploymentID string `json:"dst_minio_deployment_id"`
}

func ScanBucketReplicationRows(rows pgx.Rows) (bucketReplications []BucketReplicationRecord, err error) {
	for rows.Next() {
		var brr BucketReplicationRecord

		err = rows.Scan(
			&brr.BucketID,
			&brr.BucketName,
			&brr.SrcStorageID,
			&brr.SrcStorageAlias,
			&brr.SrcMinioDeploymentID,
			&brr.DstStorageID,
			&brr.DstStorageAlias,
			&brr.DstMinioDeploymentID,
		)
		if err != nil { return bucketReplications, util.ProcessErr(err) }

		bucketReplications = append(bucketReplications, brr)
	}

	return
}

type BucketFaaSDeploymentRecord struct {
	BucketID int64 `json:"bucket_id"`
	BucketName string `json:"bucket_name"`
	FaaSIDs []int64 `json:"faas_ids"`
	FaaSURLs []string `json:"faas_urls"`
}

func ScanBucketFaaSDeploymentRows(rows pgx.Rows) (bucketsFaaSDeployments []BucketFaaSDeploymentRecord, err error) {
	for rows.Next() {
		var bfdr BucketFaaSDeploymentRecord
		bfdr.FaaSIDs = []int64{}
		bfdr.FaaSURLs = []string{}
		ids := pgtype.Int4Array{}
		urls := pgtype.TextArray{}

		err = rows.Scan(
			&bfdr.BucketID,
			&bfdr.BucketName,
			&ids,
			&urls,
		)
		if err != nil { return bucketsFaaSDeployments, util.ProcessErr(err) }

		for _, el := range ids.Elements {
			var id64 int64
			if el.Status == pgtype.Present { el.AssignTo(&id64) }
			bfdr.FaaSIDs = append(bfdr.FaaSIDs, id64)
		}
		for _, el := range urls.Elements {
			var u string
			if el.Status == pgtype.Present { el.AssignTo(&u) }
			bfdr.FaaSURLs = append(bfdr.FaaSURLs, u)
		}

		bucketsFaaSDeployments = append(bucketsFaaSDeployments, bfdr)
	}

	return
}
