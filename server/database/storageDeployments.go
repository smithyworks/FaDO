package database

import (
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/smithyworks/FaDO/util"
)

// type facilities

type StorageDeploymentRecord struct {
	StorageID int64 `json:"storage_id"`
	ClusterID int64 `json:"cluster_id"`
	MinioDeploymentID string `json:"minio_deployment_id"`
	Alias string `json:"alias"`
	Endpoint string `json:"endpoint"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	UseSSL bool `json:"use_ssl"`
	SqsArn string `json:"sqs_arn"`
	ManagementURL string `json:"management_url"`
}

func ScanStorageDeploymentRows(rows pgx.Rows) (storageDeployments []StorageDeploymentRecord, err error) {
	for rows.Next() {
		var sd StorageDeploymentRecord

		err = rows.Scan(
			&sd.StorageID,
			&sd.ClusterID,
			&sd.MinioDeploymentID,
			&sd.Alias,
			&sd.Endpoint,
			&sd.AccessKey,
			&sd.SecretKey,
			&sd.UseSSL,
			&sd.SqsArn,
			&sd.ManagementURL,
		)
		if err != nil { return storageDeployments, util.ProcessErr(err) }

		storageDeployments = append(storageDeployments, sd)
	}

	return
}

// general query

func QueryStorageDeployments(conn DBConn, sql string, args ...interface{}) (storageDeployments []StorageDeploymentRecord, err error) {
	rows, err := Query(conn, sql, args...)
	if err != nil { return storageDeployments, util.ProcessErr(err) }
	defer rows.Close()
	
	storageDeployments, err = ScanStorageDeploymentRows(rows)
	if err != nil { return storageDeployments, util.ProcessErr(err) }

	return
}

func QueryStorageDeploymentRow(conn DBConn, sql string, args ...interface{}) (storageDeployment StorageDeploymentRecord, err error) {
	records, err := QueryStorageDeployments(conn, sql, args...)
	if err != nil { return storageDeployment, util.ProcessErr(err) }
	if len(records) != 1 { return storageDeployment, util.ProcessErr(fmt.Errorf("Expected 1 record back, go %v.", len(records))) }
	return records[0], nil
}

// insert

func InsertStorageDeployment(conn DBConn, sd StorageDeploymentRecord) (r StorageDeploymentRecord, err error) {
	records, err := QueryStorageDeployments(conn, "INSERT INTO storage_deployments (cluster_id, minio_deployment_id, alias, endpoint, access_key, secret_key, use_ssl, sqs_arn, management_url) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *",
		sd.ClusterID, sd.MinioDeploymentID, sd.Alias, sd.Endpoint, sd.AccessKey, sd.SecretKey, sd.UseSSL, sd.SqsArn, sd.ManagementURL)
	if err != nil {
		return r, util.ProcessErr(err)
	} else if len(records) != 1 {
		return r, util.ProcessErr(fmt.Errorf("Expected 1 record back, got %v.", len(records)))
	}
	return records[0], err
}
