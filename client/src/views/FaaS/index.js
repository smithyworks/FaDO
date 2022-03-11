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
} from "@mui/material";
import { AddOutlined, ExpandMore, DeleteOutlined, EditOutlined } from "@mui/icons-material";

import Page from "../../components/Page";
import api from "../../api";

import "./index.css";
import { useLocation } from "react-router";
import ResourceDialog from "../../components/ResourceDialog";

function AddFaaSDialog({ open, onClose, resources, setResources }) {
  const [url, setURL] = useState("");
  const [clusterId, setClusterId] = useState("");

  const clusters = resources?.clusters ?? [];
  const menuItems = clusters.map((c, i) => (
    <MenuItem key={i} value={c.cluster_id}>
      {c.name}
    </MenuItem>
  ));

  function onOk() {
    const data = { faas_deployment: { cluster_id: clusterId, url } };

    api
      .addFaaS(data)
      .then((resources) => setResources(resources))
      .catch((err) => console.log(err));

    setURL("");
    setClusterId("");

    onClose();
  }

  return (
    <ResourceDialog title="Add a new FaaS deployment." open={open} onClose={onClose} onOk={onOk}>
      <FormControl fullWidth variant="standard">
        <InputLabel id="add-faas-select-label">Cluster</InputLabel>
        <Select
          labelId="add-faas-select-label"
          id="add-faas-select"
          value={clusterId}
          label="Cluster"
          onChange={(e) => setClusterId(e.target.value)}
          fullWidth
        >
          {menuItems}
        </Select>
      </FormControl>
      <TextField
        size="small"
        label="URL"
        value={url}
        onChange={(e) => setURL(e.target.value)}
        fullWidth
        margin="normal"
        variant="standard"
      />
    </ResourceDialog>
  );
}

function EditFaaSDialog({ open, onClose, faas, resources, setResources }) {
  const [url, setURL] = useState(faas?.url ?? "-");
  const [clusterId, setClusterId] = useState(faas?.cluster_id ?? "");

  if (!faas) return null;

  const clusters = resources?.clusters ?? [];
  const menuItems = clusters.map((c, i) => (
    <MenuItem key={i} value={c.cluster_id}>
      {c.name}
    </MenuItem>
  ));

  function onOk() {
    const data = { faas_deployment: { faas_id: faas?.faas_id ?? 0, cluster_id: clusterId, url } };

    api
      .editFaaS(data)
      .then((resources) => setResources(resources))
      .catch((err) => console.log(err));

    setURL("");
    setClusterId("");

    onClose();
  }

  return (
    <ResourceDialog title="Edit FaaS deployment." open={open} onClose={onClose} onOk={onOk}>
      <FormControl fullWidth variant="standard">
        <InputLabel id="add-faas-select-label">Cluster</InputLabel>
        <Select
          labelId="add-faas-select-label"
          id="add-faas-select"
          value={clusterId}
          label="Cluster"
          onChange={(e) => setClusterId(e.target.value)}
          fullWidth
          disabled
        >
          {menuItems}
        </Select>
      </FormControl>
      <TextField
        size="small"
        label="URL"
        value={url}
        onChange={(e) => setURL(e.target.value)}
        fullWidth
        margin="normal"
        variant="standard"
      />
    </ResourceDialog>
  );
}

function DeleteFaaSDialog({ open, onClose, faas, setResources }) {
  if (!faas) return null;

  function onOk() {
    api
      .deleteFaaS(faas)
      .then((resources) => setResources(resources))
      .catch((err) => console.log(err));

    onClose();
  }

  const faasName = <Typography style={{ padding: "0 10px" }}>{faas?.url}</Typography>;

  return (
    <ResourceDialog title="The following resources will be deleted:" open={open} onClose={onClose} onOk={onOk}>
      <Typography variant="subtitle1">FaaS Deployments:</Typography>
      {faasName}
    </ResourceDialog>
  );
}

function FaaSRow({ faas_deployment, expanded, onEdit, onDelete }) {
  const [open, setOpen] = useState(expanded);

  const url = faas_deployment?.url ?? [];
  const zones = faas_deployment?.cluster?.zones ?? [];
  const zoneString = zones.length > 0 ? zones.join(", ") : "-";

  const cluster = faas_deployment?.cluster ?? { name: "-", cluster_id: 0 };

  return (
    <Accordion expanded={open} onChange={(_, o) => setOpen(o)}>
      <AccordionSummary expandIcon={<ExpandMore />}>
        <div className="resource-row-summary">
          <div className="resource-row-summary-title">{url}</div>
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
        </div>
        <div className="resource-row-details-buttons">
          <Button
            size="small"
            variant="contained"
            color="secondary"
            startIcon={<EditOutlined />}
            className="resource-row-details-buttons-btn"
            onClick={() => onEdit(faas_deployment)}
          >
            Edit
          </Button>
          <Button
            size="small"
            variant="contained"
            color="warning"
            startIcon={<DeleteOutlined />}
            className="resource-row-details-buttons-btn"
            onClick={() => onDelete(faas_deployment)}
          >
            Delete
          </Button>
        </div>
      </AccordionDetails>
    </Accordion>
  );
}

export default function FaaSView({ resources, setResources }) {
  const [addDialogOpen, setAddDialogOpen] = useState(false);
  const [editFaaS, setEditFaaS] = useState(false);
  const [deleteFaaS, setDeleteFaaS] = useState(false);

  const query = new URLSearchParams(useLocation().search);
  const queriedFaaSId = parseInt(query.get("faas_id"));

  const faas_deployments = resources?.faas_deployments ?? [];
  const rows = faas_deployments.map((fd, i) => (
    <FaaSRow
      faas_deployment={fd}
      expanded={fd.faas_id === queriedFaaSId}
      onEdit={(f) => setEditFaaS(f)}
      onDelete={(f) => setDeleteFaaS(f)}
      key={i}
    />
  ));

  return (
    <Page>
      <div className="resource-header">
        <h2>FaaS Deployments</h2>

        <Button
          variant="contained"
          color="primary"
          startIcon={<AddOutlined />}
          className="Deployments-header-button"
          onClick={() => setAddDialogOpen(true)}
          size="small"
        >
          Add FaaS Deployment
        </Button>
      </div>

      {rows}

      <AddFaaSDialog
        open={addDialogOpen}
        onClose={() => setAddDialogOpen(false)}
        resources={resources}
        setResources={setResources}
      />
      <EditFaaSDialog
        open={!!editFaaS}
        onClose={() => setEditFaaS(false)}
        faas={editFaaS}
        resources={resources}
        setResources={setResources}
        key={editFaaS.faas_id}
      />
      <DeleteFaaSDialog
        open={!!deleteFaaS}
        onClose={() => setDeleteFaaS(false)}
        faas={deleteFaaS}
        setResources={setResources}
        key={deleteFaaS.faas_id}
      />
    </Page>
  );
}
