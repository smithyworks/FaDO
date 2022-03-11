package mutations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/smithyworks/FaDO/cli"
	"github.com/smithyworks/FaDO/database"
	"github.com/smithyworks/FaDO/util"
)

func ConfigureLoadBalancer(conn database.DBConn) (err error) {
	if conn == nil {
		log.Println("cl2")
		if c, err := database.Acquire(); err != nil {
			return util.ProcessErr(err)
		} else {
			defer c.Release()
			conn = c
		}
	}

	// Get default policy
	var policy, matchHeader string
	err = database.GetGlobalPolicy(conn, "lb_policy", &policy)
	if err != nil { return util.ProcessErr(err) }
	err = database.GetGlobalPolicy(conn, "lb_match_header", &matchHeader)
	if err != nil { return util.ProcessErr(err) }
	port := cli.Input.LBPort
	host := cli.Input.LBDomain
	routes, err := GenerateRoutes(conn, policy, matchHeader)
	if err != nil { return util.ProcessErr(err) }

	var serversFromCaddy map[string]database.LoadBalancerServerConfig
	if resp, err := http.Get(fmt.Sprintf("%v/config/apps/http/servers/", cli.Input.CaddyAdminURL)); err != nil {
		return util.ProcessErr(err)
	} else {
		err = json.NewDecoder(resp.Body).Decode(&serversFromCaddy)
		if err != nil { return util.ProcessErr(err) }
	}

	var serverName string
	var serverConfig database.LoadBalancerServerConfig
	portPart := fmt.Sprintf(":%v", port)
	log.Print("Host", host)
	if port == "" {
		return util.ProcessErr(fmt.Errorf("Load balancing port must be specified!"))
	} else if host == "" {
		serverFound := false
		for k, v := range serversFromCaddy {
			serverFound = true
			if util.HasString(v.Listen, portPart){
				serverName = k
				serverConfig = database.LoadBalancerServerConfig{
					Listen: []string{portPart},
					Routes: routes,
				}
				break
			}
		}
		if !serverFound {
			serverName = "FaDO_LB"
			serverConfig = database.LoadBalancerServerConfig{
				Listen: []string{portPart},
				Routes: routes,
			}
		}
		serversFromCaddy[serverName] = serverConfig
	} else {
		serverFound := false
		for k, v := range serversFromCaddy {
			if util.HasString(v.Listen, portPart) {
				serverFound = true
				serverName = k
				serverConfig = v

				routeFound := false
				terminal := true
				for i, r := range serverConfig.Routes {
					if util.HasString(r.Match[0].Host, host) {
						routeFound = true
						serverConfig.Routes[i] = database.LoadBalancerRouteConfig{
							Handle: []database.HandleConfig{
								{
									Handler: "subroute",
									Routes: routes,
								},
							},
							Match: []database.MatchConfig{
								{
									Host: []string{host},
								},
							},
							Terminal: &terminal,
						}
						break
					}
				}

				if !routeFound {
					serverConfig.Routes = append(serverConfig.Routes, database.LoadBalancerRouteConfig{
						Handle: []database.HandleConfig{
							{
								Handler: "subroute",
								Routes: routes,
							},
						},
						Match: []database.MatchConfig{
							{
								Host: []string{host},
							},
						},
						Terminal: &terminal,
					})
				}

				break
			}
		}
		if !serverFound {
			serverName = "FaDO_LB"
			serverConfig = database.LoadBalancerServerConfig{
				Listen: []string{portPart},
				Routes: []database.LoadBalancerRouteConfig{
					{
						Handle: []database.HandleConfig{
							{
								Handler: "subroute",
								Routes: routes,
							},
						},
						Match: []database.MatchConfig{
							{
								Host: []string{host},
							},
						},
					},
				},
				Terminal: true,
			}	
		}
		serversFromCaddy[serverName] = serverConfig
	}

	configJSON, err := json.Marshal(serversFromCaddy)
	if err != nil { return util.ProcessErr(err) }

	go postToCaddyWithDelay(fmt.Sprintf("%v/config/apps/http/servers/", cli.Input.CaddyAdminURL), bytes.NewBuffer(configJSON))

	return nil
}

func postToCaddyWithDelay(url string, bytes io.Reader) {
	time.Sleep(100 * time.Millisecond)
	http.Post(url, "application/json", bytes)
}

func GenerateRoutes(conn database.DBConn, policy, matchHeader string) (routes []database.LoadBalancerRouteConfig, err error) {
	// Get bucket and faas associations
	rows, err := database.Query(conn, "SELECT * FROM buckets_faas_deployments")
	if err != nil { return routes, util.ProcessErr(err) }
	bucketsFaaSDeployments, err := database.ScanBucketFaaSDeploymentRows(rows)
	if err != nil { return routes, util.ProcessErr(err) }

	// Get eventual route overrides
	var routeOverridesMap map[string]database.LoadBalancerRouteSettings
	err = database.GetGlobalPolicy(conn, "lb_route_overrides", &routeOverridesMap);
	if err != nil { return routes, util.ProcessErr(err) }
	newRoutesOverridesMap := make(map[string]database.LoadBalancerRouteSettings)

	routesMap := make(map[string]database.LoadBalancerRouteSettings)
	for _, bfd := range bucketsFaaSDeployments {
		rs, isOverridden := routeOverridesMap[bfd.BucketName]
		if isOverridden {
			rs.BucketName = bfd.BucketName
			newRoutesOverridesMap[bfd.BucketName] = rs
		} else {
			rs.Policy = policy
			rs.Upstreams = util.MakeStringSet(bfd.FaaSURLs)
			rs.BucketName = bfd.BucketName
		}

		matcher := database.MatchConfig{
			Header: map[string][]string{
				matchHeader: {rs.BucketName},
			},
		}

		upstreams := make([]database.UpstreamConfig, 0)
		for _, fe := range rs.Upstreams {
			if fe != "" { upstreams = append(upstreams, database.UpstreamConfig{Dial: fe}) }
		}

		handler := database.HandleConfig{
			Handler: "reverse_proxy",
			LoadBalancing: &database.LoadBalancingConfig{
				SelectionPolicy: database.SelectionPolicyConfig{
					Policy: rs.Policy,
				},
			},
			Upstreams: upstreams,
		}

		route := database.LoadBalancerRouteConfig{
			Handle: []database.HandleConfig{handler},
			Match: []database.MatchConfig{matcher},
		}
		
		routes = append(routes, route)
		routesMap[rs.BucketName] = rs
	}

	if err := database.SetGlobalPolicy(conn, "lb_routes", routesMap); err != nil {
		return routes, util.ProcessErr(err)
	}
	if err := database.SetGlobalPolicy(conn, "lb_route_overrides", newRoutesOverridesMap); err != nil {
		return routes, util.ProcessErr(err)
	}

	return
}