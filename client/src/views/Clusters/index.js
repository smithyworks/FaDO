import React, { useState } from "react";
import {
  Button,
  TextField,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  Divider,
  Typography,
  Checkbox,
} from "@mui/material";
import { AddOutlined, ExpandMore, DeleteOutlined, EditOutlined } from "@mui/icons-material";

import Page from "../../components/Page";
import api from "../../api";

import "./index.css";
import { useLocation } from "react-router";
import ResourceDialog from "../../components/ResourceDialog";

function AddClusterDialog({ open, onClose, resources, setResources }) {
  const [name, setName] = useState("");
  const [zonesStr, setZonesStr] = useState("");

  function onOk() {
    const data = {
      cluster: { name },
      zones:
        zonesStr.trim() === ""
          ? []
          : zonesStr
              .trim()
              .split(",")
              .map((s) => s.trim()),
    };

    api
      .addCluster(data)
      .then((resources) => setResources(resources))
      .catch((err) => console.log(err));

    setName("");
    setZonesStr("");

    onClose();
  }

  return (
    <ResourceDialog title="Add a new cluster." open={open} onClose={onClose} onOk={onOk}>
      <TextField
        size="small"
        label="Name"
        value={name}
        onChange={(e) => setName(e.target.value)}
        fullWidth
        margin="normal"
        variant="standard"
      />
      <TextField
        size="small"
        label="Zones (Comma Separated List)"
        value={zonesStr}
        onChange={(e) => setZonesStr(e.target.value)}
        fullWidth
        margin="normal"
        variant="standard"
      />
    </ResourceDialog>
  );
}

function EditClusterDialog({ open, onClose, cluster, resources, setResources }) {
  const [name, setName] = useState(cluster?.name ?? "-");
  const [zonesStr, setZonesStr] = useState(cluster?.zones?.join(", ") ?? "");

  if (!cluster) return null;

  function onOk() {
    const data = {
      cluster: { cluster_id: cluster?.cluster_id, name },
      zones: zonesStr.trim() === "" ? [] : zonesStr.split(",").map((s) => s.trim()),
    };

    api
      .editCluster(data)
      .then((resources) => setResources(resources))
      .catch((err) => console.log(err));

    setName("");
    setZonesStr("");

    onClose();
  }

  return (
    <ResourceDialog title="Edit cluster." open={open} onClose={onClose} onOk={onOk}>
      <TextField size="small" label="Name" value={name} fullWidth margin="normal" variant="standard" disabled />
      <TextField
        size="small"
        label="Zones (Comma Separated List)"
        value={zonesStr}
        onChange={(e) => setZonesStr(e.target.value)}
        fullWidth
        margin="normal"
        variant="standard"
      />
    </ResourceDialog>
  );
}

function DeleteClusterDialog({ open, onClose, cluster, setResources }) {
  const [permanent, setPermanent] = useState(false);

  if (!cluster) return null;

  function onOk() {
    api
      .deleteCluster(cluster, permanent)
      .then((resources) => setResources(resources))
      .catch((err) => console.log(err));

    onClose();
  }

  const clusterName = <Typography style={{ padding: "0 10px" }}>{cluster?.name ?? ""}</Typography>;
  const faasNames = cluster?.faas_deployments?.map((fd, i) => (
    <Typography style={{ padding: "0 10px" }} key={i}>
      {fd.url}
    </Typography>
  ));
  const storageNames = [],
    bucketNames = [],
    objectNames = [];
  cluster?.storage_deployments?.forEach((sd, i) => {
    storageNames.push(
      <Typography style={{ padding: "0 10px" }} key={i}>
        {sd?.alias}
      </Typography>
    );
    sd?.buckets?.forEach((b, j) => {
      bucketNames.push(
        <Typography style={{ padding: "0 10px" }} key={j}>
          {b?.name}
        </Typography>
      );
      b?.objects?.forEach((o, k) => {
        objectNames.push(
          <Typography style={{ padding: "0 10px" }} key={k}>
            {b?.name}/{o?.name}
          </Typography>
        );
      });
    });
  });

  return (
    <ResourceDialog title="The following resources will be deleted:" open={open} onClose={onClose} onOk={onOk}>
      <Typography variant="subtitle1">Clusters</Typography>
      {clusterName}
      {faasNames.length > 0 && (
        <Typography variant="subtitle1" style={{ marginTop: 10 }}>
          FaaS Deployments
        </Typography>
      )}
      {faasNames}
      {storageNames.length > 0 && (
        <Typography variant="subtitle1" style={{ marginTop: 10 }}>
          Storage Deployments
        </Typography>
      )}
      {storageNames}
      {bucketNames.length > 0 && (
        <Typography variant="subtitle1" style={{ marginTop: 10 }}>
          Buckets
        </Typography>
      )}
      {bucketNames}
      {objectNames.length > 0 && (
        <Typography variant="subtitle1" style={{ marginTop: 10 }}>
          Objects
        </Typography>
      )}
      {objectNames}
      {bucketNames.length > 0 && (
        <Typography style={{ marginTop: 10 }}>
          <Checkbox
            value={permanent}
            onChange={(e) => setPermanent(e.target.checked)}
            style={{ margin: "-5px -5px 0 -10px" }}
          />{" "}
          Permanently delete buckets and objects from storage.
        </Typography>
      )}
    </ResourceDialog>
  );
}

function ClusterRow({ cluster, expanded, onEdit, onDelete }) {
  const [open, setOpen] = useState(!!expanded);

  const name = cluster?.name ?? "-";
  const zones = cluster?.zones ?? [];
  const zoneString = zones.length > 0 ? zones.join(", ") : "-";
  const faasDeployments = cluster?.faas_deployments ?? [];
  const storageDeployements = cluster?.storage_deployments ?? [];

  const faasDeploymentDetails = faasDeployments.map((fd, i) => {
    return (
      <a key={i} className="resource-row-details-link" href={`/faas?faas_id=${fd.faas_id}`}>
        {fd.url}
      </a>
    );
  });

  const storageDeploymentDetails = storageDeployements.map((sd, i) => {
    return (
      <a key={i} className="resource-row-details-link" href={`/storage?storage_id=${sd.storage_id}`}>
        {sd.alias}
      </a>
    );
  });

  return (
    <Accordion expanded={open} onChange={(_, o) => setOpen(o)}>
      <AccordionSummary expandIcon={<ExpandMore />}>
        <div className="resource-row-summary">
          <div className="resource-row-summary-title">{name}</div>
          <div className="resource-row-summary-prop">Zones: {zoneString}</div>
        </div>
      </AccordionSummary>
      <AccordionDetails>
        <Divider />
        <div className="resource-row-details">
          <div className="resource-row-details-title">FaaS Deployments:</div>
          {faasDeploymentDetails}
          <div className="resource-row-details-title">Storage Deployments:</div>
          {storageDeploymentDetails}
        </div>
        <div className="resource-row-details-buttons">
          <Button
            size="small"
            variant="contained"
            color="secondary"
            startIcon={<EditOutlined />}
            className="resource-row-details-buttons-btn"
            onClick={() => onEdit(cluster)}
          >
            Edit
          </Button>
          <Button
            size="small"
            variant="contained"
            color="warning"
            startIcon={<DeleteOutlined />}
            className="resource-row-details-buttons-btn"
            onClick={() => onDelete(cluster)}
          >
            Delete
          </Button>
        </div>
      </AccordionDetails>
    </Accordion>
  );
}

export default function ClustersView({ resources, setResources }) {
  const [addDialogOpen, setAddDialogOpen] = useState(false);
  const [editCluster, setEditCluster] = useState(false);
  const [deleteCluster, setDeleteCluster] = useState(false);

  const query = new URLSearchParams(useLocation().search);
  const queriedClusterId = parseInt(query.get("cluster_id"));

  const clusters = resources?.clusters ?? [];
  const rows = clusters.map((c, i) => (
    <ClusterRow
      cluster={c}
      expanded={c.cluster_id === queriedClusterId}
      onEdit={(c) => setEditCluster(c)}
      onDelete={(c) => setDeleteCluster(c)}
      key={i}
    />
  ));

  return (
    <Page>
      <div className="resource-header">
        <h2>Clusters</h2>

        <Button
          variant="contained"
          color="primary"
          startIcon={<AddOutlined />}
          className="Clusters-header-button"
          onClick={() => setAddDialogOpen(true)}
          size="small"
        >
          Add Cluster
        </Button>
      </div>

      {rows}

      <AddClusterDialog
        open={addDialogOpen}
        onClose={() => setAddDialogOpen(false)}
        resources={resources}
        setResources={setResources}
      />
      <EditClusterDialog
        open={!!editCluster}
        onClose={() => setEditCluster(false)}
        cluster={editCluster}
        resources={resources}
        setResources={setResources}
        key={editCluster.cluster_id}
      />
      <DeleteClusterDialog
        open={!!deleteCluster}
        onClose={() => setDeleteCluster(false)}
        cluster={deleteCluster}
        setResources={setResources}
        key={deleteCluster.cluster_id}
      />
    </Page>
  );
}
