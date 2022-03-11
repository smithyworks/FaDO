import React, { useEffect, useState } from "react";
import { BrowserRouter as Router, Switch, Route, Link, Redirect } from "react-router-dom";
import { useLocation } from "react-router";
import {
  FunctionsOutlined,
  StorageOutlined,
  DashboardOutlined,
  CloudOutlined,
  InboxOutlined,
  DescriptionOutlined,
  CallSplitOutlined,
} from "@mui/icons-material";

import DashboardView from "./views/Dashboard";
import ClustersView from "./views/Clusters";
import FaaSView from "./views/FaaS";
import StorageView from "./views/Storage";
import BucketsView from "./views/Buckets";
import ObjectsView from "./views/Objects";
import LoadBalancingView from "./views/LoadBalancing";

import "./App.css";
import api from "./api";
import { CssBaseline, Typography } from "@mui/material";

function prepareData(resources) {
  /*
    policies
    global_policies
    clusters
    clusters_policies
    faas_deployments
    storage_deployments
    buckets
    buckets_policies
    replica_bucket_locations
    objects
   */
  console.log(resources);

  // policies
  const policies = resources?.policies ?? [];
  const policyMap = policies.reduce((acc, curr) => {
    acc[curr.policy_id] = curr;
    return acc;
  }, {});
  const global_policies = resources?.global_policies ?? [];
  global_policies.forEach((gp) => {
    gp.name = policyMap[gp?.policy_id]?.name;
  });
  const clusters_policies = resources?.clusters_policies ?? [];
  clusters_policies.forEach((cp) => {
    cp.name = policyMap[cp?.policy_id]?.name;
  });
  const buckets_policies = resources?.buckets_policies ?? [];
  buckets_policies.forEach((cp) => {
    cp.name = policyMap[cp?.policy_id]?.name;
  });

  // Clusters
  const clusters = resources?.clusters ?? [];
  clusters.forEach((cz) => {
    cz.zones = [];
    cz.faas_deployments = [];
    cz.storage_deployments = [];
  });
  const clusterMap = clusters.reduce((acc, curr) => {
    acc[curr.cluster_id] = curr;
    return acc;
  }, {});
  // Clusters Zones
  const clustersZones = clusters_policies.filter((cz) => cz.name === "zones") ?? [];
  clustersZones.forEach((cz) => (clusterMap[cz.cluster_id].zones = JSON.parse(cz.value)));

  // FaaS Deployments
  const faas_deployments = resources?.faas_deployments ?? [];
  faas_deployments.forEach((fd) => {
    fd.cluster = clusterMap[fd.cluster_id];
    clusterMap[fd.cluster_id].faas_deployments.push(fd);
  });

  // Storage Deployments
  const storage_deployments = resources?.storage_deployments ?? [];
  storage_deployments.forEach((sd) => {
    sd.cluster = clusterMap[sd.cluster_id];
    sd.buckets = [];
    sd.replica_buckets = [];
  });
  const storageDeploymentMap = storage_deployments.reduce((acc, curr) => {
    clusterMap[curr.cluster_id].storage_deployments.push(curr);
    acc[curr.storage_id] = curr;
    return acc;
  }, {});

  // Buckets
  const buckets = resources?.buckets ?? [];
  buckets.forEach((b) => {
    b.storage_deployment = storageDeploymentMap[b.storage_id];
    b.allowed_zones = [];
    b.replica_storage_deployments = [];
    b.objects = [];
    b.target_replica_count = 0;
    b.replication_overridden = false;
  });
  const bucketMap = buckets.reduce((acc, curr) => {
    storageDeploymentMap[curr.storage_id].buckets.push(curr);
    acc[curr.bucket_id] = curr;
    return acc;
  }, {});
  // Buckets Allowed Zones
  const bucketsAllowedZones = buckets_policies.filter((bp) => bp.name === "zones");
  bucketsAllowedZones.forEach((baz) => (bucketMap[baz.bucket_id].allowed_zones = JSON.parse(baz.value)));
  // target replica count
  const bucketsTargetReplicaCounts = buckets_policies.filter((bp) => bp.name === "target_replica_count");
  bucketsTargetReplicaCounts.forEach((tc) => (bucketMap[tc.bucket_id].target_replica_count = JSON.parse(tc.value)));
  // replication override
  const replicaLocationPolicies = buckets_policies.filter((bp) => bp.name === "replica_locations");
  replicaLocationPolicies.forEach((tc) => (bucketMap[tc.bucket_id].replication_overridden = true));
  // Replica Buckets Storage Deployments
  const replicaBucketLocations = resources?.replica_bucket_locations ?? [];
  replicaBucketLocations.forEach((rbl) => {
    storageDeploymentMap[rbl.storage_id].replica_buckets.push(bucketMap[rbl.bucket_id]);
    bucketMap[rbl.bucket_id].replica_storage_deployments.push(storageDeploymentMap[rbl.storage_id]);
  });

  const objects = resources?.objects ?? [];
  objects.forEach((o) => {
    bucketMap[o.bucket_id].objects.push(o);
    o.bucket = bucketMap[o.bucket_id];
  });

  return {
    policies,
    global_policies,
    clusters_policies,
    clusters,
    faas_deployments,
    storage_deployments,
    buckets,
    buckets_policies,
    objects,
    load_balancer_config: resources.load_balancer_config ?? {},
    load_balancer_settings: {
      policy: resources?.load_balancer_policy ?? "-",
      match_header: resources?.load_balancer_match_header ?? "-",
      host: resources?.load_balancer_host ?? "-",
      port: resources?.load_balancer_port ?? "-",
    },
    load_balancer_routes: resources?.load_balancer_routes ?? {},
    load_balancer_route_overrides: resources?.load_balancer_route_overrides ?? {},
  };
}

function AppLink({ name, to, icon: IconComponent }) {
  const loc = useLocation();
  const active = loc.pathname === to;
  return (
    <Link className={`App-link ${active ? "active" : ""}`} to={to}>
      <IconComponent className="App-name-icon" />
      <Typography
        display="inline-block"
        style={{ fontSize: "1.1rem", fontWeight: 400, position: "relative", bottom: -2 }}
      >
        {name}
      </Typography>
    </Link>
  );
}

function App() {
  const [resources, setResources] = useState({ deployments: [], buckets: [], objects: [] });

  useEffect(() => {
    async function fetchResources() {
      try {
        const resources = await api.listResources();
        setResources(prepareData(resources));
      } catch (err) {
        console.log(err);
      }
    }
    fetchResources();
  }, []);

  return (
    <Router>
      <div className="App">
        <div className="App-sidebar">
          <div className="App-name">FaDO</div>
          <div className="App-subtitle">The Function and Data Orchestrator</div>
          <AppLink name="Dashboard" icon={DashboardOutlined} to="/" />
          <AppLink name="Clusters" icon={CloudOutlined} to="/clusters" />
          <AppLink name="FaaS Deployments" icon={FunctionsOutlined} to="/faas" />
          <AppLink name="Storage Deployments" icon={StorageOutlined} to="/storage" />
          <AppLink name="Buckets" icon={InboxOutlined} to="/buckets" />
          <AppLink name="Objects" icon={DescriptionOutlined} to="/objects" />
          <AppLink name="Load Balancing" icon={CallSplitOutlined} to="/lb" />
        </div>
        <div className="App-content">
          <Switch>
            <Route exact path="/">
              <DashboardView resources={resources} />
            </Route>
            <Route path="/clusters">
              <ClustersView resources={resources} setResources={(r) => setResources(prepareData(r))} />
            </Route>
            <Route path="/faas">
              <FaaSView resources={resources} setResources={(r) => setResources(prepareData(r))} />
            </Route>
            <Route path="/storage">
              <StorageView resources={resources} setResources={(r) => setResources(prepareData(r))} />
            </Route>
            <Route path="/buckets">
              <BucketsView resources={resources} setResources={(r) => setResources(prepareData(r))} />
            </Route>
            <Route path="/objects">
              <ObjectsView resources={resources} setResources={(r) => setResources(prepareData(r))} />
            </Route>
            <Route path="/lb">
              <LoadBalancingView resources={resources} setResources={(r) => setResources(prepareData(r))} />
            </Route>
            <Route path="*">
              <Redirect to="/" />
            </Route>
          </Switch>
        </div>
      </div>
      <CssBaseline />
    </Router>
  );
}

export default App;
