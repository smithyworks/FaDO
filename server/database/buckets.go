package database

import (
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/smithyworks/FaDO/util"
)

// type facilities

type BucketRecord struct {
	BucketID int64 `json:"bucket_id"`
	StorageID int64 `json:"storage_id"`
	Name string `json:"name"`
}

func ScanBucketRows(rows pgx.Rows) (buckets []BucketRecord, err error) {
	for rows.Next() {
		var br BucketRecord

		err = rows.Scan(
			&br.BucketID,
			&br.StorageID,
			&br.Name,
		)
		if err != nil { return buckets, util.ProcessErr(err) }

		buckets = append(buckets, br)
	}

	return
}

// general query

func QueryBuckets(conn DBConn, sql string, args ...interface{}) (buckets []BucketRecord, err error) {
	rows, err := Query(conn, sql, args...)
	if err != nil { return buckets, util.ProcessErr(err) }
	defer rows.Close()
	
	buckets, err = ScanBucketRows(rows)
	if err != nil { return buckets, util.ProcessErr(err) }

	return
}

func QueryBucketRow(conn DBConn, sql string, args ...interface{}) (bucket BucketRecord, err error) {
	records, err := QueryBuckets(conn, sql, args...)
	if err != nil { return bucket, util.ProcessErr(err) }
	if len(records) != 1 { return bucket, util.ProcessErr(fmt.Errorf("Expected 1 record back, go %v.", len(records))) }
	return records[0], nil
}

// insert

func InsertBucket(conn DBConn, bucket BucketRecord) (r BucketRecord, err error) {
	records, err := QueryBuckets(conn, `INSERT INTO buckets (storage_id, name) VALUES ($1, $2) RETURNING *`, bucket.StorageID, bucket.Name)
	if err != nil {
		return r, util.ProcessErr(err)
	} else if len(records) != 1 {
		return r, util.ProcessErr(fmt.Errorf("Expected 1 record back, got %v.", len(records)))
	}
	return records[0], err
}
