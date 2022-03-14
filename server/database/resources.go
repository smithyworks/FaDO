package database

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/smithyworks/FaDO/cli"
	"github.com/smithyworks/FaDO/util"
)


type ResourceCollection struct {
	Policies []PolicyRecord `json:"policies"`
	GlobalPolicies []GlobalPolicyRecord `json:"global_policies"`
	Clusters []ClusterRecord `json:"clusters"`
	ClustersPolicies []ClusterPolicyRecord `json:"clusters_policies"`
	FaaSDeployments []FaaSDeploymentRecord `json:"faas_deployments"`
	StorageDeployments []StorageDeploymentRecord `json:"storage_deployments"`
	Buckets []BucketRecord `json:"buckets"`
	BucketsPolicies []BucketPolicyRecord `json:"buckets_policies"`
	ReplicaBucketsLocations []ReplicaBucketLocationRecord `json:"replica_bucket_locations"`
	Objects []ObjectRecord `json:"objects"`
	LoadBalancerConfig map[string]LoadBalancerServerConfig `json:"load_balancer_config"`
	LoadBalancerHost string `json:"load_balancer_host"`
	LoadBalancerPort string `json:"load_balancer_port"`
	LoadBalancerMatchHeader string `json:"load_balancer_match_header"`
	LoadBalancerPolicy string `json:"load_balancer_policy"`
	LoadBalancerRoutes map[string]LoadBalancerRouteSettings `json:"load_balancer_routes"`
	LoadBalancerRouteOverrides map[string]LoadBalancerRouteSettings `json:"load_balancer_route_overrides"`
}

func QueryResources() (resources ResourceCollection, err error) {
	if conn, err := Acquire(); err != nil {
		return resources, util.ProcessErr(err)
	} else {
		defer conn.Release()
		if resources.Policies, err = QueryPolicies(conn, "SELECT * FROM policies"); err != nil {
			return resources, util.ProcessErr(err)
		} else if resources.GlobalPolicies, err = QueryGlobalPolicies(conn, "SELECT * FROM global_policies"); err != nil {
			return resources, util.ProcessErr(err)
		} else if resources.Clusters, err = QueryClusters(conn, "SELECT * FROM clusters ORDER BY cluster_id ASC, name"); err != nil {
			return resources, util.ProcessErr(err)
		} else if resources.ClustersPolicies, err = QueryClustersPolicies(conn, "SELECT * FROM clusters_policies"); err != nil {
			return resources, util.ProcessErr(err)
		} else if resources.FaaSDeployments, err = QueryFaaSDeployments(conn, "SELECT * FROM faas_deployments ORDER BY cluster_id ASC, url"); err != nil {
			return resources, util.ProcessErr(err)
		} else if resources.StorageDeployments, err = QueryStorageDeployments(conn, "SELECT * FROM storage_deployments ORDER BY cluster_id ASC, alias"); err != nil {
			return resources, util.ProcessErr(err)
		} else if resources.Buckets, err = QueryBuckets(conn, "SELECT * FROM buckets ORDER BY storage_id ASC, name"); err != nil {
			return resources, util.ProcessErr(err)
		} else if resources.BucketsPolicies, err = QueryBucketsPolicies(conn, "SELECT * FROM buckets_policies"); err != nil {
			return resources, util.ProcessErr(err)
		} else if resources.ReplicaBucketsLocations, err = QueryReplicaBucketLocations(conn, "SELECT * FROM replica_bucket_locations"); err != nil {
			return resources, util.ProcessErr(err)
		} else if resources.Objects, err = QueryObjects(conn, "SELECT * FROM objects ORDER BY bucket_id ASC, name"); err != nil {
			return resources, util.ProcessErr(err)
		}

		resources.LoadBalancerHost = cli.Input.LBDomain
		resources.LoadBalancerPort = cli.Input.LBPort
		
		var str string
		err := GetGlobalPolicy(conn, "lb_match_header", &str)
		if err != nil {
			return resources, util.ProcessErr(err)
		} else { resources.LoadBalancerMatchHeader = str }
		err = GetGlobalPolicy(conn, "lb_policy", &str)
		if err != nil {
			return resources, util.ProcessErr(err)
		} else { resources.LoadBalancerPolicy = str }
		
		var r map[string]LoadBalancerRouteSettings
		err = GetGlobalPolicy(conn, "lb_routes", &r)
		if err != nil {
			return resources, util.ProcessErr(err)
		} else { resources.LoadBalancerRoutes = r }
		var ro map[string]LoadBalancerRouteSettings
		err = GetGlobalPolicy(conn, "lb_route_overrides", &ro)
		if err != nil {
			return resources, util.ProcessErr(err)
		} else { resources.LoadBalancerRouteOverrides = ro }
		
		if resp, err := http.Get(fmt.Sprintf("%v/config/apps/http/servers/", cli.Input.CaddyAdminURL)); err != nil {
			return resources, util.ProcessErr(err)
		} else {
			err = json.NewDecoder(resp.Body).Decode(&resources.LoadBalancerConfig)
			if err != nil { return resources, util.ProcessErr(err) }
		}
	}

	return
}