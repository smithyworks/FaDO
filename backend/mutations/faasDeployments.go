package mutations

import (
	"github.com/smithyworks/FaDO/database"
	"github.com/smithyworks/FaDO/util"
)

func AddFaaSDeployment(conn database.DBConn, faasDeployment database.FaaSDeploymentRecord) (err error) {
	if faasDeployments, err := database.QueryFaaSDeployments(conn, "SELECT * FROM faas_deployments WHERE faas_id = $1", faasDeployment.FaaSID); err != nil {
		return util.ProcessErr(err)
	} else if len(faasDeployments) == 0 {
		if _, err = database.InsertFaaSDeployment(conn, faasDeployment); err != nil {
			return util.ProcessErr(err)
		}
	} else {
		// record already exists, no need to proceed
		return nil
	}

	// for all master buckets
	//   ResolveBucketReplicas
	if buckets, err := database.QueryBuckets(conn, "SELECT * FROM buckets"); err != nil {
		return util.ProcessErr(err)
	} else {
		for _, b := range buckets {
			if err = ResolveBucketReplicas(conn, b); err != nil {
				return util.ProcessErr(err)
			}
		}
	}

	return
}

func EditFaaSDeployment(conn database.DBConn, faasDeployment database.FaaSDeploymentRecord) (err error) {
	// Can only update endpoint

	if _, err = database.Exec(conn, "UPDATE faas_deployments SET url = $1 WHERE faas_id = $2", faasDeployment.URL, faasDeployment.FaaSID); err != nil {
		return util.ProcessErr(err)
	}

	// for all master buckets
	//   ResolveBucketReplicas
	if buckets, err := database.QueryBuckets(conn, "SELECT * FROM buckets"); err != nil {
		return util.ProcessErr(err)
	} else {
		for _, b := range buckets {
			if err = ResolveBucketReplicas(conn, b); err != nil {
				return util.ProcessErr(err)
			}
		}
	}

	return
}

func DeleteFaaSDeployment(conn database.DBConn, faasDeployment database.FaaSDeploymentRecord) (err error) {
	if _, err = database.Exec(conn, "DELETE FROM faas_deployments WHERE faas_id = $1", faasDeployment.FaaSID); err != nil {
		return util.ProcessErr(err)
	}

	// for all master buckets
	//   ResolveBucketReplicas
	if buckets, err := database.QueryBuckets(conn, "SELECT * FROM buckets"); err != nil {
		return util.ProcessErr(err)
	} else {
		for _, b := range buckets {
			if err = ResolveBucketReplicas(conn, b); err != nil {
				return util.ProcessErr(err)
			}
		}
	}

	return
}
