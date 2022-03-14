package mutations

import (
	"github.com/smithyworks/FaDO/database"
	"github.com/smithyworks/FaDO/util"
)

func AddCluster(conn database.DBConn, cluster database.ClusterRecord, zones []string) (err error) {
	// insert into database
	if clusterRecords, err := database.QueryClusters(conn, "SELECT * FROM clusters WHERE name = $1", cluster.Name); err != nil {
		return util.ProcessErr(err)
	} else if len(clusterRecords) == 0 {
		if cluster, err = database.InsertCluster(conn, cluster); err != nil {
			return util.ProcessErr(err)
		} else if zones != nil { // insert cluster/zone record
			err = database.SetClusterPolicy(conn, cluster, "zones", zones)
			if err != nil { return util.ProcessErr(err) }
		}
	}

	return
}

func EditCluster(conn database.DBConn, cluster database.ClusterRecord, zones []string) (err error) {
	if err := database.SetClusterPolicy(conn, cluster, "zones", zones); err != nil {
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

func DeleteCluster(conn database.DBConn, cluster database.ClusterRecord, permanent bool) (err error) {
	// for storage deployments
	//   DeleteStorageDeployment
	if storageDeployments, err := database.QueryStorageDeployments(conn, "SELECT * FROM storage_deployments WHERE cluster_id = $1", cluster.ClusterID); err != nil {
		return util.ProcessErr(err)
	} else {
		for _, sd := range storageDeployments {
			if err = DeleteStorageDeployment(conn, sd, permanent); err != nil {
				return util.ProcessErr(err)
			}
		}
	}

	// for faas deployments
	//   DeleteFaaSDeployment
	if faasDeployments, err := database.QueryFaaSDeployments(conn, "SELECT * FROM faas_deployments WHERE cluster_id = $1", cluster.ClusterID); err != nil {
		return util.ProcessErr(err)
	} else {
		for _, fd := range faasDeployments {
			if err = DeleteFaaSDeployment(conn, fd); err != nil {
				return util.ProcessErr(err)
			}
		}
	}


	// delete cluster from database
	_, err = database.Exec(conn, "DELETE FROM clusters WHERE cluster_id = $1", cluster.ClusterID)

	return util.ProcessErr(err)
}
