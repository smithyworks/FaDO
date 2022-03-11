package database

import (
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/smithyworks/FaDO/util"
)

// Policies

// type facilities

type PolicyRecord struct {
	PolicyID int64 `json:"policy_id"`
	Name string `json:"name"`
	DefaultValue string `json:"default_value"`
}

func ScanPolicyRows(rows pgx.Rows) (policies []PolicyRecord, err error) {
	for rows.Next() {
		var pr PolicyRecord

		err = rows.Scan(
			&pr.PolicyID,
			&pr.Name,
			&pr.DefaultValue,
		)
		if err != nil { return policies, util.ProcessErr(err) }

		policies = append(policies, pr)
	}

	return
}

// general query

func QueryPolicies(conn DBConn, sql string, args ...interface{}) (policies []PolicyRecord, err error) {
	rows, err := Query(conn, sql, args...)
	if err != nil { return policies, util.ProcessErr(err) }
	defer rows.Close()

	policies, err = ScanPolicyRows(rows)
	if err != nil { return policies, util.ProcessErr(err) }

	return
}

func QueryPolicyRow(conn DBConn, sql string, args ...interface{}) (policy PolicyRecord, err error) {
	records, err := QueryPolicies(conn, sql, args...)
	if err != nil { return policy, util.ProcessErr(err) }
	if len(records) != 1 { return policy, util.ProcessErr(fmt.Errorf("Expected 1 record back, go %v.", len(records))) }
	return records[0], nil
}

// Global Policies

type GlobalPolicyRecord struct {
	PolicyID int64 `json:"policy_id"`
	Value string `json:"value"`
}

func QueryGlobalPolicies(conn DBConn, sql string, args ...interface{}) (globalPolicies []GlobalPolicyRecord, err error) {
	rows, err := Query(conn, sql, args...)
	if err != nil { return globalPolicies, util.ProcessErr(err) }

	defer rows.Close()
	for rows.Next() {
		var gpr GlobalPolicyRecord
		err = rows.Scan(&gpr.PolicyID, &gpr.Value)
		if err != nil { return globalPolicies, util.ProcessErr(err) }
		globalPolicies = append(globalPolicies, gpr)
	}

	return
}

func UpsertGlobalPolicy(conn DBConn, gp GlobalPolicyRecord) (r GlobalPolicyRecord, err error) {
	records, err := QueryGlobalPolicies(conn, "INSERT INTO global_policies (policy_id, value) VALUES ($1, $2) ON CONFLICT (policy_id) DO UPDATE SET value = $2 RETURNING *", gp.PolicyID, gp.Value)
	if err != nil {
		return r, util.ProcessErr(err)
	} else if len(records) != 1 {
		return r, util.ProcessErr(fmt.Errorf("Expected 1 record back, got %v.", len(records)))
	}
	return records[0], err
}

// Bucket policies

type BucketPolicyRecord struct {
	BucketID int64 `json:"bucket_id"`
	PolicyID int64 `json:"policy_id"`
	Value string `json:"value"`
}

func QueryBucketsPolicies(conn DBConn, sql string, args ...interface{}) (bucketPolicies []BucketPolicyRecord, err error) {
	rows, err := Query(conn, sql, args...)
	if err != nil { return bucketPolicies, util.ProcessErr(err) }

	defer rows.Close()
	for rows.Next() {
		var bpr BucketPolicyRecord
		err = rows.Scan(&bpr.BucketID, &bpr.PolicyID, &bpr.Value)
		if err != nil { return bucketPolicies, util.ProcessErr(err) }
		bucketPolicies = append(bucketPolicies, bpr)
	}

	return
}

func UpsertBucketPolicy(conn DBConn, bp BucketPolicyRecord) (r BucketPolicyRecord, err error) {
	records, err := QueryBucketsPolicies(conn, "INSERT INTO buckets_policies (bucket_id, policy_id, value) VALUES ($1, $2, $3) ON CONFLICT (bucket_id, policy_id) DO UPDATE SET value = $3 RETURNING *", bp.BucketID, bp.PolicyID, bp.Value)
	if err != nil {
		return r, util.ProcessErr(err)
	} else if len(records) != 1 {
		return r, util.ProcessErr(fmt.Errorf("Expected 1 record back, got %v.", len(records)))
	}
	return records[0], err
}

// Bucket policies

type ClusterPolicyRecord struct {
	ClusterID int64 `json:"cluster_id"`
	PolicyID int64 `json:"policy_id"`
	Value string `json:"value"`
}

func QueryClustersPolicies(conn DBConn, sql string, args ...interface{}) (clusterPolicies []ClusterPolicyRecord, err error) {
	rows, err := Query(conn, sql, args...)
	if err != nil { return clusterPolicies, util.ProcessErr(err) }

	defer rows.Close()
	for rows.Next() {
		var bpr ClusterPolicyRecord
		err = rows.Scan(&bpr.ClusterID, &bpr.PolicyID, &bpr.Value)
		if err != nil { return clusterPolicies, util.ProcessErr(err) }
		clusterPolicies = append(clusterPolicies, bpr)
	}

	return
}

func UpsertClusterPolicy(conn DBConn, cp ClusterPolicyRecord) (r ClusterPolicyRecord, err error) {
	records, err := QueryClustersPolicies(conn, "INSERT INTO clusters_policies (cluster_id, policy_id, value) VALUES ($1, $2, $3) ON CONFLICT (cluster_id, policy_id) DO UPDATE SET value = $3 RETURNING *", cp.ClusterID, cp.PolicyID, cp.Value)
	if err != nil {
		return r, util.ProcessErr(err)
	} else if len(records) != 1 {
		return r, util.ProcessErr(fmt.Errorf("Expected 1 record back, got %v.", len(records)))
	}
	return records[0], err
}

// Abstractions

func GetBucketPolicy(conn DBConn, bucket BucketRecord, policyName string, value interface{}) ( set bool, err error) {
	policy, err := QueryPolicyRow(conn, "SELECT * FROM policies WHERE name = $1", policyName)
	if err != nil { return set, util.ProcessErr(err) }

	err = json.Unmarshal([]byte(policy.DefaultValue), value)
	if err != nil { return set, util.ProcessErr(err) }

	bucketPolicies, err := QueryBucketsPolicies(conn, "SELECT * FROM buckets_policies WHERE bucket_id = $1 AND policy_id = $2", bucket.BucketID, policy.PolicyID)
	if err != nil { return set, util.ProcessErr(err) }

	if len(bucketPolicies) != 1 { return }
	
	err = json.Unmarshal([]byte(bucketPolicies[0].Value), value)
	if err != nil { return set, nil }
	set = true

	return
}

func SetBucketPolicy(conn DBConn, bucket BucketRecord, policyName string, input interface{}) (err error) {
	policy, err := QueryPolicyRow(conn, "SELECT * FROM policies WHERE name = $1", policyName)
	if err != nil { return util.ProcessErr(err) }

	valueBytes, err := json.Marshal(input)
	if err != nil { return util.ProcessErr(err) }

	_, err = UpsertBucketPolicy(conn, BucketPolicyRecord{bucket.BucketID, policy.PolicyID, string(valueBytes)})
	if err != nil { return util.ProcessErr(err) }

	return
}

func DeleteBucketPolicy(conn DBConn, bucket BucketRecord, policyName string) (err error) {
	policy, err := QueryPolicyRow(conn, "SELECT * FROM policies WHERE name = $1", policyName)
	if err != nil { return util.ProcessErr(err) }

	_, err = Exec(conn, "DELETE FROM buckets_policies WHERE bucket_id = $1 AND policy_id = $2", bucket.BucketID, policy.PolicyID)
	if err != nil { return util.ProcessErr(err) }

	return
}

func GetClusterPolicy(conn DBConn, cluster ClusterRecord, policyName string, value interface{}) (set bool, err error) {
	policy, err := QueryPolicyRow(conn, "SELECT * FROM policies WHERE name = $1", policyName)
	if err != nil { return set, util.ProcessErr(err) }

	err = json.Unmarshal([]byte(policy.DefaultValue), value)
	if err != nil { return set, util.ProcessErr(err) }

	clusterPolicies, err := QueryClustersPolicies(conn, "SELECT * FROM clusters_policies WHERE cluster_id = $1 AND policy_id = $2", cluster.ClusterID, policy.PolicyID)
	if err != nil { return set, util.ProcessErr(err) }

	if len(clusterPolicies) != 1 { return }
	
	err = json.Unmarshal([]byte(clusterPolicies[0].Value), value)
	if err != nil { return set, nil }
	set = true

	return
}

func SetClusterPolicy(conn DBConn, cluster ClusterRecord, policyName string, input interface{}) (err error) {
	policy, err := QueryPolicyRow(conn, "SELECT * FROM policies WHERE name = $1", policyName)
	if err != nil { return util.ProcessErr(err) }

	valueBytes, err := json.Marshal(input)
	if err != nil { return util.ProcessErr(err) }

	_, err = UpsertClusterPolicy(conn, ClusterPolicyRecord{cluster.ClusterID, policy.PolicyID, string(valueBytes)})
	if err != nil { return util.ProcessErr(err) }

	return
}

func DeleteClusterPolicy(conn DBConn, cluster ClusterRecord, policyName string) (err error) {
	policy, err := QueryPolicyRow(conn, "SELECT * FROM policies WHERE name = $1", policyName)
	if err != nil { return util.ProcessErr(err) }

	_, err = Exec(conn, "DELETE FROM clusters_policies WHERE cluster_id = $1 AND policy_id = $2", cluster.ClusterID, policy.PolicyID)
	if err != nil { return util.ProcessErr(err) }

	return
}

func GetGlobalPolicy(conn DBConn, policyName string, value interface{}) (err error) {
	policy, err := QueryPolicyRow(conn, "SELECT * FROM policies WHERE name = $1", policyName)
	if err != nil { return util.ProcessErr(err) }

	err = json.Unmarshal([]byte(policy.DefaultValue), value)
	if err != nil { return util.ProcessErr(err) }

	globalPolicies, err := QueryGlobalPolicies(conn, "SELECT * FROM global_policies WHERE policy_id = $1", policy.PolicyID)
	if err != nil { return util.ProcessErr(err) }

	if len(globalPolicies) != 1 { return }
	
	err = json.Unmarshal([]byte(globalPolicies[0].Value), value)
	if err != nil { return nil }

	return
}

func SetGlobalPolicy(conn DBConn, policyName string, input interface{}) (err error) {
	policy, err := QueryPolicyRow(conn, "SELECT * FROM policies WHERE name = $1", policyName)
	if err != nil { return util.ProcessErr(err) }

	valueBytes, err := json.Marshal(input)
	if err != nil { return util.ProcessErr(err) }

	_, err = UpsertGlobalPolicy(conn, GlobalPolicyRecord{policy.PolicyID, string(valueBytes)})
	if err != nil { return util.ProcessErr(err) }

	return
}

func DeleteGlobalPolicy(conn DBConn, policyName string) (err error) {
	policy, err := QueryPolicyRow(conn, "SELECT * FROM policies WHERE name = $1", policyName)
	if err != nil { return util.ProcessErr(err) }

	_, err = Exec(conn, "DELETE FROM global_policies WHERE policy_id = $1", policy.PolicyID)
	if err != nil { return util.ProcessErr(err) }

	return
}