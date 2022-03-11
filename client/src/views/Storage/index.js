import React, { useState } from "react";
import {
  Button,
  TextField,
  MenuItem,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  Divider,
  Select,
  InputLabel,
  FormControl,
  FormGroup,
  FormControlLabel,
  Checkbox,
  Typography,
} from "@mui/material";
import { AddOutlined, ExpandMore, DeleteOutlined, ExitToAppOutlined } from "@mui/icons-material";

import Page from "../../components/Page";
import api from "../../api";

import "./index.css";
import { useLocation } from "react-router";
import ResourceDialog from "../../components/ResourceDialog";

function AddStorageDialog({ open, onClose, storage, resources, setResources }) {
  const [alias, setAlias] = useState("");
  const [endpoint, setEndpoint] = useState("");
  const [accessKey, setAccessKey] = useState("");
  const [secretKey, setSecretKey] = useState("");
  const [useSSL, setUseSSL] = useState("");
  const [clusterId, setClusterId] = useState("");
  const [managementUrl, setManagementUrl] = useState("");

  const clusters = resources?.clusters ?? [];
  const menuItems = clusters.map((c, i) => (
    <MenuItem key={i} value={c.cluster_id}>
      {c.name}
    </MenuItem>
  ));

  function onOk() {
    const data = {
      storage_deployment: {
        cluster_id: clusterId,
        alias,
        management_url: managementUrl,
        endpoint,
        access_key: accessKey,
        secret_key: secretKey,
        use_ssl: useSSL,
      },
    };

    api
      .addStorage(data)
      .then((resources) => setResources(resources))
      .catch((err) => console.log(err));

    setClusterId("");
    setAlias("");
    setManagementUrl("");
    setEndpoint("");
    setAccessKey("");
    setSecretKey("");
    setUseSSL(false);

    onClose();
  }

  return (
    <ResourceDialog title="Add a storage deployment." open={open} onClose={onClose} onOk={onOk}>
      <FormControl fullWidth variant="standard">
        <InputLabel id="add-storage-select-label">Cluster</InputLabel>
        <Select
          labelId="add-storage-select-label"
          id="add-storage-select"
          value={clusterId}
          onChange={(e) => setClusterId(e.target.value)}
          label="Cluster"
          fullWidth
        >
          {menuItems}
        </Select>
      </FormControl>
      <TextField
        size="small"
        label="Alias"
        value={alias}
        onChange={(e) => setAlias(e.target.value)}
        fullWidth
        margin="normal"
        variant="standard"
      />
      <TextField
        size="small"
        label="Management URL"
        value={managementUrl}
        onChange={(e) => setManagementUrl(e.target.value)}
        fullWidth
        margin="normal"
        variant="standard"
      />
      <TextField
        size="small"
        label="Endpoint"
        value={endpoint}
        onChange={(e) => setEndpoint(e.target.value)}
        fullWidth
        margin="normal"
        variant="standard"
      />
      <TextField
        size="small"
        label="Access Key"
        value={accessKey}
        onChange={(e) => setAccessKey(e.target.value)}
        fullWidth
        margin="normal"
        variant="standard"
      />
      <TextField
        size="small"
        label="Secret Key"
        value={secretKey}
        onChange={(e) => setSecretKey(e.target.value)}
        fullWidth
        margin="normal"
        variant="standard"
      />
      <FormGroup>
        <FormControlLabel
          control={<Checkbox />}
          label="Use SSL"
          value={!!useSSL}
          onChange={(e) => setUseSSL(!!e.target.checked)}
        />
      </FormGroup>
    </ResourceDialog>
  );
}

function DeleteStorageDialog({ open, onClose, storage, setResources }) {
  const [permanent, setPermanent] = useState(false);

  if (!storage) return null;

  function onOk() {
    api
      .deleteStorage(storage, permanent)
      .then((resources) => setResources(resources))
      .catch((err) => console.log(err));

    onClose();
  }

  const bucketNames = [],
    objectNames = [];
  const storageName = <Typography style={{ padding: "0 10px" }}>{storage?.alias}</Typography>;
  storage?.buckets?.forEach((b, j) => {
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

  return (
    <ResourceDialog title="The following resources will be deleted:" open={open} onClose={onClose} onOk={onOk}>
      <Typography variant="subtitle1" style={{ marginTop: 10 }}>
        Storage Deployments
      </Typography>
      {storageName}
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

function StorageRow({ storage_deployment, expanded, onDelete }) {
  const [open, setOpen] = useState(expanded);

  const alias = storage_deployment?.alias ?? "-";
  const endpoint = storage_deployment?.endpoint ?? [];
  const zones = storage_deployment?.cluster?.zones ?? [];
  const zoneString = zones.length > 0 ? zones.join(", ") : "-";

  const cluster = storage_deployment?.cluster ?? { name: "-", cluster_id: 0 };

  const buckets = storage_deployment?.buckets ?? [];
  const bucketDetails = buckets.map((b, i) => {
    return (
      <a key={i} className="resource-row-details-link" href={`/buckets?bucket_id=${b.bucket_id}`}>
        {b.name}
      </a>
    );
  });

  const replicas = storage_deployment?.replica_buckets ?? [];
  const replicaDetails = replicas.map((r, i) => {
    return (
      <a key={i} className="resource-row-details-link" href={`/buckets?bucket_id=${r.bucket_id}`}>
        {r.name}
      </a>
    );
  });

  return (
    <Accordion expanded={open} onChange={(_, o) => setOpen(o)}>
      <AccordionSummary expandIcon={<ExpandMore />}>
        <div className="resource-row-summary">
          <div className="resource-row-summary-title">{alias}</div>
          <div className="resource-row-summary-prop">{endpoint}</div>
          <div className="resource-row-summary-prop">Zones: {zoneString}</div>
        </div>
      </AccordionSummary>
      <AccordionDetails>
        <Divider />
        <div className="resource-row-details">
          <div className="resource-row-details-title">Cluster:</div>
          <a className="resource-row-details-link" href={`/clusters?cluster_id=${cluster.cluster_id}`}>
            {cluster.name}
          </a>
          <div className="resource-row-details-title">Master Buckets:</div>
          {bucketDetails}
          <div className="resource-row-details-title">Replica Buckets:</div>
          {replicaDetails}
        </div>
        <div className="resource-row-details-buttons">
          {storage_deployment?.management_url && storage_deployment?.management_url !== "" && (
            <Button
              size="small"
              variant="contained"
              color="info"
              startIcon={<ExitToAppOutlined />}
              className="resource-row-details-buttons-btn"
              style={{ height: 30.75, position: "relative", bottom: -5 }}
              href={storage_deployment?.management_url}
              target="_blank"
            >
              Go To Console
            </Button>
          )}
          <Button
            size="small"
            variant="contained"
            color="warning"
            startIcon={<DeleteOutlined />}
            className="resource-row-details-buttons-btn"
            onClick={() => onDelete(storage_deployment)}
          >
            Delete
          </Button>
        </div>
      </AccordionDetails>
    </Accordion>
  );
}

export default function StorageView({ resources, setResources }) {
  const [addDialogOpen, setAddDialogOpen] = useState(false);
  const [deleteStorage, setDeleteStorage] = useState(false);

  const query = new URLSearchParams(useLocation().search);
  const queriedStorageId = parseInt(query.get("storage_id"));

  const storage_deployments = resources?.storage_deployments ?? [];
  const rows = storage_deployments.map((sd, i) => (
    <StorageRow
      storage_deployment={sd}
      expanded={sd.storage_id === queriedStorageId}
      onDelete={(s) => setDeleteStorage(s)}
      key={i}
    />
  ));

  return (
    <Page>
      <div className="resource-header">
        <h2>Storage Deployments</h2>

        <Button
          variant="contained"
          color="primary"
          startIcon={<AddOutlined />}
          className="Deployments-header-button"
          onClick={() => setAddDialogOpen(true)}
          size="small"
        >
          Add Storage Deployment
        </Button>
      </div>

      {rows}

      <AddStorageDialog
        open={addDialogOpen}
        onClose={() => setAddDialogOpen(false)}
        resources={resources}
        setResources={setResources}
      />
      <DeleteStorageDialog
        open={!!deleteStorage}
        onClose={() => setDeleteStorage(false)}
        storage={deleteStorage}
        resources={resources}
        setResources={setResources}
        key={deleteStorage.storage_id}
      />
    </Page>
  );
}
