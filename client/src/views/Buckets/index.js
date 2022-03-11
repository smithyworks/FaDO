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
  Typography,
  Checkbox,
  Collapse,
} from "@mui/material";
import { AddOutlined, ExpandMore, DeleteOutlined, EditOutlined, WarningAmberOutlined } from "@mui/icons-material";

import Page from "../../components/Page";
import api from "../../api";

import "./index.css";
import { useLocation } from "react-router";
import ResourceDialog from "../../components/ResourceDialog";

function AddBucketDialog({ open, onClose, resources, setResources }) {
  const [name, setName] = useState("");
  const [zonesStr, setZonesStr] = useState("");
  const [storageId, setStorageId] = useState("");
  const [replicaCount, setReplicaCount] = useState(0);
  const [manual, setManual] = useState(false);
  const [replicaLocations, setReplicaLocations] = useState(new Set());

  const storage_deployments = resources?.storage_deployments ?? [];
  const menuItems = storage_deployments.map((sd, i) => (
    <MenuItem key={i} value={sd.storage_id}>
      {sd.alias} (Cluster: {sd?.cluster?.name})
    </MenuItem>
  ));

  replicaLocations.delete(storageId);
  const replicaOptions = storage_deployments.filter((sd) => sd.storage_id !== storageId);
  const replicaCheckboxes = replicaOptions.map((sd, i) => (
    <Typography key={i}>
      <Checkbox
        checked={replicaLocations.has(sd.storage_id)}
        onChange={(e) => {
          if (e.target.checked) {
            replicaLocations.add(sd.storage_id);
          } else {
            replicaLocations.delete(sd.storage_id);
          }
          setReplicaLocations(new Set(replicaLocations));
        }}
        style={{ margin: "-5px -5px 0 -10px" }}
      />{" "}
      {sd.alias}
    </Typography>
  ));

  function onOk() {
    const data = {
      bucket: {
        storage_id: storageId,
        name,
        replication_overridden: !!manual,
      },
      target_replica_count: parseInt(replicaCount),
      zones:
        zonesStr.trim() === ""
          ? []
          : zonesStr
              .trim()
              .split(",")
              .map((s) => s.trim()),
    };
    if (!!manual) data.replica_storage_ids = [...replicaLocations];

    api
      .addBucket(data)
      .then((resources) => setResources(resources))
      .catch((err) => console.log(err));

    setName("");
    setZonesStr("");
    setStorageId("");
    setReplicaCount(0);
    setManual(false);
    setReplicaLocations(new Set());

    onClose();
  }

  return (
    <ResourceDialog title="Add a new bucket." open={open} onClose={onClose} onOk={onOk}>
      <FormControl fullWidth variant="standard">
        <InputLabel id="add-bucket-select-label">Storage Deployment</InputLabel>
        <Select
          labelId="add-bucket-select-label"
          id="add-bucket-select"
          value={storageId}
          label="Storage Deployment"
          onChange={(e) => setStorageId(e.target.value)}
          fullWidth
        >
          {menuItems}
        </Select>
      </FormControl>
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
        label="Allowed Zones (Comma Separated List)"
        value={zonesStr}
        onChange={(e) => setZonesStr(e.target.value)}
        fullWidth
        margin="normal"
        variant="standard"
      />
      <TextField
        id="outlined-number"
        label="Target Replica Count"
        type="number"
        InputProps={{ inputProps: { min: 0, max: 10 } }}
        value={replicaCount}
        onChange={(e) => setReplicaCount(e.target.value)}
        fullWidth
        margin="normal"
        variant="standard"
      />
      <Typography>
        <Checkbox
          checked={manual}
          onChange={(e) => setManual(e.target.checked)}
          style={{ margin: "-5px -5px 0 -10px" }}
        />{" "}
        Manually define replica locations.
      </Typography>
      <Collapse in={manual}>
        <Typography>Select locations for replication:</Typography>
        {replicaCheckboxes}
      </Collapse>
    </ResourceDialog>
  );
}

function EditBucketDialog({ open, onClose, bucket, resources, setResources }) {
  const [name, setName] = useState(bucket?.name ?? "-");
  const [zonesStr, setZonesStr] = useState(bucket?.allowed_zones?.join(", ") ?? "");
  const [storageId, setStorageId] = useState(bucket?.storage_id ?? 0);
  const [replicaCount, setReplicaCount] = useState(bucket?.target_replica_count ?? 0);
  const [manual, setManual] = useState(bucket?.replication_overridden);
  const [replicaLocations, setReplicaLocations] = useState(
    new Set(bucket?.replica_storage_deployments?.map((rsd) => rsd.storage_id) ?? [])
  );

  if (!bucket) return null;

  const storage_deployments = resources?.storage_deployments ?? [];
  const menuItems = storage_deployments.map((sd, i) => (
    <MenuItem key={i} value={sd.storage_id}>
      {sd.alias} (Cluster: {sd?.cluster?.name})
    </MenuItem>
  ));

  const replicaOptions = storage_deployments.filter((sd) => sd.storage_id !== bucket.storage_id);
  const replicaCheckboxes = replicaOptions.map((sd, i) => (
    <Typography key={i}>
      <Checkbox
        checked={replicaLocations.has(sd.storage_id)}
        onChange={(e) => {
          if (e.target.checked) {
            replicaLocations.add(sd.storage_id);
          } else {
            replicaLocations.delete(sd.storage_id);
          }
          setReplicaLocations(new Set(replicaLocations));
        }}
        style={{ margin: "-5px -5px 0 -10px" }}
      />{" "}
      {sd.alias}
    </Typography>
  ));

  function onOk() {
    const data = {
      bucket: {
        bucket_id: bucket?.bucket_id,
        storage_id: storageId,
        name,
      },
      target_replica_count: parseInt(replicaCount),
      zones: zonesStr.trim() === "" ? [] : zonesStr.split(",").map((s) => s.trim()),
    };
    if (!!manual) data.replica_storage_ids = [...replicaLocations];

    api
      .editBucket(data)
      .then((resources) => setResources(resources))
      .catch((err) => console.log(err));

    setName("");
    setZonesStr("");
    setStorageId("");
    setReplicaCount(0);
    setManual(false);
    setReplicaLocations(new Set());

    onClose();
  }

  return (
    <ResourceDialog title="Edit bucket." open={open} onClose={onClose} onOk={onOk}>
      <FormControl fullWidth variant="standard">
        <InputLabel id="add-bucket-select-label">Storage Deployment</InputLabel>
        <Select
          labelId="add-bucket-select-label"
          id="add-bucket-select"
          value={storageId}
          label="Storage Deployment"
          fullWidth
          disabled
        >
          {menuItems}
        </Select>
      </FormControl>
      <TextField size="small" label="Name" value={name} fullWidth margin="normal" variant="standard" disabled />
      <TextField
        size="small"
        label="Allowed Zones (Comma Separated List)"
        value={zonesStr}
        onChange={(e) => setZonesStr(e.target.value)}
        fullWidth
        margin="normal"
        variant="standard"
      />
      <TextField
        id="outlined-number"
        label="Target Replica Count"
        type="number"
        InputProps={{ inputProps: { min: 0, max: 10, value: replicaCount } }}
        value={replicaCount}
        onChange={(e) => setReplicaCount(e.target.value)}
        fullWidth
        margin="normal"
        variant="standard"
      />
      <Typography>
        <Checkbox
          checked={manual}
          onChange={(e) => setManual(e.target.checked)}
          style={{ margin: "-5px -5px 0 -10px" }}
        />{" "}
        Manually define replica locations.
      </Typography>
      <Collapse in={manual}>
        <Typography>Select locations for replication:</Typography>
        {replicaCheckboxes}
      </Collapse>
    </ResourceDialog>
  );
}

function DeleteBucketDialog({ open, onClose, bucket, resources, setResources }) {
  if (!bucket) return null;

  function onOk() {
    api
      .deleteBucket(bucket)
      .then((resources) => setResources(resources))
      .catch((err) => console.log(err));

    onClose();
  }

  const bucketName = <Typography style={{ padding: "0 10px" }}>{bucket?.name}</Typography>;
  const objectNames = [];
  bucket?.objects?.forEach((o, k) => {
    objectNames.push(
      <Typography style={{ padding: "0 10px" }} key={k}>
        {bucket?.name}/{o?.name}
      </Typography>
    );
  });

  return (
    <ResourceDialog title="The following resources will be deleted:" open={open} onClose={onClose} onOk={onOk}>
      <Typography variant="subtitle1" style={{ marginTop: 10 }}>
        Buckets
      </Typography>
      {bucketName}
      <Typography variant="subtitle1" style={{ marginTop: 10 }}>
        Objects
      </Typography>
      {objectNames}
      <Typography style={{ marginTop: 10 }}>
        <WarningAmberOutlined style={{ margin: "0 5px -5px 0" }} />
        The bucket and its objects will be permanently deleted from storage.
      </Typography>
    </ResourceDialog>
  );
}

function BucketRow({ bucket, expanded, onEdit, onDelete }) {
  const [open, setOpen] = useState(expanded);

  const name = bucket?.name ?? "-";
  const allowed_zones = bucket?.allowed_zones ?? [];
  const zoneString = allowed_zones.length > 0 ? allowed_zones.join(", ") : "-";
  const replica_count = bucket?.target_replica_count ?? 0;

  const origin = bucket?.storage_deployment ?? { name: "alias", storge_id: 0 };

  const replicaLocations = bucket?.replica_storage_deployments ?? [];
  const replicaDetails = replicaLocations.map((sd, i) => {
    return (
      <a key={i} className="resource-row-details-link" href={`/storage?storage_id=${sd.storage_id}`}>
        {sd?.alias ?? "-"}
      </a>
    );
  });

  const objects = bucket?.objects ?? [];
  const objectDetails = objects.map((o, i) => {
    return (
      <a key={i} className="resource-row-details-link" href={`/objects?object_id=${o.object_id}`}>
        {name}/{o.name ?? "-"}
      </a>
    );
  });

  return (
    <Accordion expanded={open} onChange={(_, o) => setOpen(o)}>
      <AccordionSummary expandIcon={<ExpandMore />}>
        <div className="resource-row-summary">
          <div className="resource-row-summary-title">{name}</div>
          <div className="resource-row-summary-prop">Target Replica Count: {replica_count}</div>
          <div className="resource-row-summary-prop">Allowed Zones: {zoneString}</div>
        </div>
      </AccordionSummary>
      <AccordionDetails>
        <div className="resource-row-details">
          <Divider />
          <div className="resource-row-details-title">Origin Storage Deployment:</div>
          <a className="resource-row-details-link" href={`/storage?storage_id=${origin.storage_id}`}>
            {origin.alias}
          </a>
          <div className="resource-row-details-title">
            Replication Storage Deployments:{" "}
            {bucket?.replication_overridden && (
              <Typography
                color="red"
                variant="body2"
                display="inline-block"
                style={{
                  fontSize: "0.75rem",
                  marginLeft: 10,
                  border: "1px solid red",
                  borderRadius: 4,
                  padding: "2px 4px 0px 4px",
                }}
              >
                OVERRIDDEN
              </Typography>
            )}
          </div>
          {replicaDetails}
          <div className="resource-row-details-title">Objects:</div>
          {objectDetails}
          <div className="resource-row-details-buttons">
            <Button
              size="small"
              variant="contained"
              color="secondary"
              startIcon={<EditOutlined />}
              className="resource-row-details-buttons-btn"
              onClick={() => onEdit(bucket)}
            >
              Edit
            </Button>
            <Button
              size="small"
              variant="contained"
              color="warning"
              startIcon={<DeleteOutlined />}
              className="resource-row-details-buttons-btn"
              onClick={() => onDelete(bucket)}
            >
              Delete
            </Button>
          </div>
        </div>
      </AccordionDetails>
    </Accordion>
  );
}

export default function BucketsView({ resources, setResources }) {
  const [addDialogOpen, setAddDialogOpen] = useState(false);
  const [editBucket, setEditBucket] = useState(false);
  const [deleteBucket, setDeleteBucket] = useState(false);

  const query = new URLSearchParams(useLocation().search);
  const queriedBucketId = parseInt(query.get("bucket_id"));

  const buckets = resources?.buckets ?? [];
  const rows = buckets.map((b, i) => (
    <BucketRow
      bucket={b}
      expanded={b.bucket_id === queriedBucketId}
      onEdit={(b) => setEditBucket(b)}
      onDelete={(b) => setDeleteBucket(b)}
      key={i}
    />
  ));

  return (
    <Page>
      <div className="resource-header">
        <h2>Buckets</h2>

        <Button
          variant="contained"
          color="primary"
          startIcon={<AddOutlined />}
          className="Deployments-header-button"
          onClick={() => setAddDialogOpen(true)}
          size="small"
        >
          Add Bucket
        </Button>
      </div>

      {rows}

      <AddBucketDialog
        open={addDialogOpen}
        onClose={() => setAddDialogOpen(false)}
        resources={resources}
        setResources={setResources}
      />
      <EditBucketDialog
        open={!!editBucket}
        onClose={() => setEditBucket(false)}
        bucket={editBucket}
        resources={resources}
        setResources={setResources}
        key={editBucket.bucket_id}
      />
      <DeleteBucketDialog
        open={!!deleteBucket}
        onClose={() => setDeleteBucket(false)}
        bucket={deleteBucket}
        resources={resources}
        setResources={setResources}
        key={deleteBucket.bucket_id}
      />
    </Page>
  );
}
