package database

import (
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/smithyworks/FaDO/util"
)

// type facilities

type ObjectRecord struct {
	ObjectID int64 `json:"object_id"`
	BucketID int64 `json:"bucket_id"`
	Name string `json:"name"`
}

func ScanObjectRows(rows pgx.Rows) (objects []ObjectRecord, err error) {
	for rows.Next() {
		var or ObjectRecord

		err = rows.Scan(
			&or.ObjectID,
			&or.BucketID,
			&or.Name,
		)
		if err != nil { return objects, util.ProcessErr(err) }

		objects = append(objects, or)
	}

	return
}

// general query

func QueryObjects(conn DBConn, sql string, args ...interface{}) (objects []ObjectRecord, err error) {
	rows, err := Query(conn, sql, args...)
	if err != nil { return objects, util.ProcessErr(err) }
	defer rows.Close()
	
	objects, err = ScanObjectRows(rows)
	if err != nil { return objects, util.ProcessErr(err) }

	return
}

func QueryObjectRow(conn DBConn, sql string, args ...interface{}) (object ObjectRecord, err error) {
	records, err := QueryObjects(conn, sql, args...)
	if err != nil { return object, util.ProcessErr(err) }
	if len(records) != 1 { return object, util.ProcessErr(fmt.Errorf("Expected 1 record back, go %v.", len(records))) }
	return records[0], nil
}

// insert

func InsertObject(conn DBConn, obj ObjectRecord) (r ObjectRecord, err error) {
	records, err := QueryObjects(conn, "INSERT INTO objects (bucket_id, name) VALUES ($1, $2) RETURNING *", obj.BucketID, obj.Name)
	if err != nil {
		return r, util.ProcessErr(err)
	} else if len(records) != 1 {
		return r, util.ProcessErr(fmt.Errorf("Expected 1 record back, got %v.", len(records)))
	}
	return records[0], err
}
