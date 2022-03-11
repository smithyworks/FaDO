import { Paper } from "@mui/material";
import React from "react";
import Page from "../../components/Page";

import "./index.css";

export default function DashboardView({ resources }) {
  return (
    <Page>
      <div className="Dashboard">
        <div className="resource-header">
          <h2>Dashboard</h2>
        </div>
        <Paper
          component="a"
          className="Dashboard-card"
          href="/clusters"
          style={{ backgroundColor: "#9B2226", color: "white" }}
        >
          <div className="Dashboard-card-number">{resources?.clusters?.length ?? "-"}</div>
          <div className="Dashboard-card-name">Clusters</div>
        </Paper>
        <Paper
          component="a"
          className="Dashboard-card"
          href="/faas"
          style={{ backgroundColor: "#AE2012", color: "white" }}
        >
          <div className="Dashboard-card-number">{resources?.faas_deployments?.length ?? "-"}</div>
          <div className="Dashboard-card-name">FaaS Deployments</div>
        </Paper>
        <Paper
          component="a"
          className="Dashboard-card"
          href="/storage"
          style={{ backgroundColor: "#BB3E03", color: "white" }}
        >
          <div className="Dashboard-card-number">{resources?.storage_deployments?.length ?? "-"}</div>
          <div className="Dashboard-card-name">Storage Deployments</div>
        </Paper>
        <Paper
          component="a"
          className="Dashboard-card"
          href="/buckets"
          style={{ backgroundColor: "#CA6702", color: "white" }}
        >
          <div className="Dashboard-card-number">{resources?.buckets?.length ?? "-"}</div>
          <div className="Dashboard-card-name">Buckets</div>
        </Paper>
        <Paper
          component="a"
          className="Dashboard-card"
          href="/objects"
          style={{ backgroundColor: "#EE9B00", color: "white" }}
        >
          <div className="Dashboard-card-number">{resources?.objects?.length ?? "-"}</div>
          <div className="Dashboard-card-name">Objects</div>
        </Paper>
      </div>
    </Page>
  );
}
