package database

import (
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/smithyworks/FaDO/util"
)

// type facilities

type FaaSDeploymentRecord struct {
	FaaSID int64 `json:"faas_id"`
	ClusterID int64 `json:"cluster_id"`
	URL string `json:"url"`
}

func ScanFaaSDeploymentRows(rows pgx.Rows) (faasDeployments []FaaSDeploymentRecord, err error) {
	for rows.Next() {
		var fd FaaSDeploymentRecord

		err = rows.Scan(
			&fd.FaaSID,
			&fd.ClusterID,
			&fd.URL,
		)
		if err != nil { return faasDeployments, util.ProcessErr(err) }

		faasDeployments = append(faasDeployments, fd)
	}

	return
}

// general query

func QueryFaaSDeployments(conn DBConn, sql string, args ...interface{}) (faasDeployments []FaaSDeploymentRecord, err error) {
	rows, err := Query(conn, sql, args...)
	if err != nil { return faasDeployments, util.ProcessErr(err) }
	defer rows.Close()
	
	faasDeployments, err = ScanFaaSDeploymentRows(rows)
	if err != nil { return faasDeployments, util.ProcessErr(err) }

	return
}

func QueryFaaSDeploymentRow(conn DBConn, sql string, args ...interface{}) (faasDeployment FaaSDeploymentRecord, err error) {
	records, err := QueryFaaSDeployments(conn, sql, args...)
	if err != nil { return faasDeployment, util.ProcessErr(err) }
	if len(records) != 1 { return faasDeployment, util.ProcessErr(fmt.Errorf("Expected 1 record back, go %v.", len(records))) }
	return records[0], nil
}

// insert

func InsertFaaSDeployment(conn DBConn, faas FaaSDeploymentRecord) (r FaaSDeploymentRecord, err error) {
	records, err := QueryFaaSDeployments(conn, "INSERT INTO faas_deployments (cluster_id, url) VALUES ($1, $2) RETURNING *", faas.ClusterID, faas.URL)
	if err != nil {
		return r, util.ProcessErr(err)
	} else if len(records) != 1 {
		return r, util.ProcessErr(fmt.Errorf("Expected 1 record back, got %v.", len(records)))
	}
	return records[0], err
}
