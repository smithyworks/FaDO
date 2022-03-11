package database

import (
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/smithyworks/FaDO/util"
)

// type facilities

type ClusterRecord struct {
	ClusterID int64 `json:"cluster_id"`
	Name string `json:"name"`
}

func ScanClusterRows(rows pgx.Rows) (clusters []ClusterRecord, err error) {
	for rows.Next() {
		var cr ClusterRecord

		err = rows.Scan(
			&cr.ClusterID,
			&cr.Name,
		)
		if err != nil { return clusters, util.ProcessErr(err) }

		clusters = append(clusters, cr)
	}

	return
}

// general query

func QueryClusters(conn DBConn, sql string, args ...interface{}) (clusters []ClusterRecord, err error) {
	rows, err := Query(conn, sql, args...)
	if err != nil { return clusters, util.ProcessErr(err) }
	defer rows.Close()
	
	clusters, err = ScanClusterRows(rows)
	if err != nil { return clusters, util.ProcessErr(err) }

	return
}

func QueryClusterRow(conn DBConn, sql string, args ...interface{}) (cluster ClusterRecord, err error) {
	records, err := QueryClusters(conn, sql, args...)
	if err != nil { return cluster, util.ProcessErr(err) }
	if len(records) != 1 { return cluster, util.ProcessErr(fmt.Errorf("Expected 1 record back, go %v.", len(records))) }
	return records[0], nil
}

// insert

func InsertCluster(conn DBConn, cluster ClusterRecord) (r ClusterRecord, err error) {
	records, err := QueryClusters(conn, "INSERT INTO clusters (name) VALUES ($1) RETURNING *", cluster.Name)
	if err != nil {
		return r, util.ProcessErr(err)
	} else if len(records) != 1 {
		return r, util.ProcessErr(fmt.Errorf("Expected 1 record back, got %v.", len(records)))
	}
	return records[0], err
}
