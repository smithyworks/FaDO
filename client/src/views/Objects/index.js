import React, { useState } from "react";
import {
  Button,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  Divider,
  Input,
  InputLabel,
  FormControl,
  Select,
  MenuItem,
  Typography,
} from "@mui/material";
import { AddOutlined, DeleteOutlined, DownloadOutlined, ExpandMore, WarningAmberOutlined } from "@mui/icons-material";

import Page from "../../components/Page";
import ResourceDialog from "../../components/ResourceDialog";

import "./index.css";
import api from "../../api";
import { useLocation } from "react-router";

function AddObjectDialog({ open, onClose, resources, setResources }) {
  const [bucketId, setBucketId] = useState("");
  const [file, setFile] = useState(null);

  const buckets = resources?.buckets ?? [];
  const menuItems = buckets.map((b, i) => (
    <MenuItem key={i} value={b.bucket_id}>
      {b.name}
    </MenuItem>
  ));

  function onOk() {
    const data = new FormData();
    data.append("bucket", `${bucketId}`);
    data.append("file", file, file.name);

    api
      .addObject(data)
      .then((resources) => setResources(resources))
      .catch((err) => console.log(err));

    setBucketId("");
    setFile(null);

    onClose();
  }

  return (
    <ResourceDialog title="Add a new object." open={open} onClose={onClose} onOk={onOk}>
      <FormControl fullWidth variant="standard">
        <InputLabel id="add-object-select-label">Bucket</InputLabel>
        <Select
          labelId="add-object-select-label"
          id="add-object-select"
          value={bucketId}
          label="Bucket"
          onChange={(e) => setBucketId(e.target.value)}
          fullWidth
          margin="normal"
        >
          {menuItems}
        </Select>
      </FormControl>

      <InputLabel style={{ paddingTop: 16 }}>Upload a file:</InputLabel>
      <Input name="Upload a file" type="file" fullWidth onChange={(e) => setFile(e.target.files[0])} />
    </ResourceDialog>
  );
}

function DeleteObjectDialog({ open, onClose, object, setResources }) {
  if (!object) return null;

  function onOk() {
    api
      .deleteObject(object)
      .then((resources) => setResources(resources))
      .catch((err) => console.log(err));

    onClose();
  }

  const objectName = (
    <Typography style={{ padding: "0 10px" }}>
      {object?.bucket?.name}/{object?.name}
    </Typography>
  );

  return (
    <ResourceDialog title="The following resources will be deleted:" open={open} onClose={onClose} onOk={onOk}>
      <Typography variant="subtitle1" style={{ marginTop: 10 }}>
        Objects
      </Typography>
      {objectName}
      <Typography style={{ marginTop: 10 }}>
        <WarningAmberOutlined style={{ margin: "0 5px -5px 0" }} />
        The object will be permanently deleted from storage.
      </Typography>
    </ResourceDialog>
  );
}

function ObjectRow({ object, expanded, onDelete }) {
  const [open, setOpen] = useState(expanded);

  const name = object?.name ?? "-";

  const bucketName = object?.bucket?.name ?? "-";
  const bucket = object?.bucket ?? { name: "-", bucket_id: 0 };

  return (
    <Accordion expanded={open} onChange={(_, o) => setOpen(o)}>
      <AccordionSummary expandIcon={<ExpandMore />}>
        <div className="resource-row-summary">
          <div className="resource-row-summary-title">
            {bucketName}/{name}
          </div>
        </div>
      </AccordionSummary>
      <AccordionDetails>
        <Divider />
        <div className="resource-row-details">
          <div className="resource-row-details-title">Bucket:</div>
          <a className="resource-row-details-link" href={`/buckets?bucket_id=${bucket.bucket_id}`}>
            {bucket.name}
          </a>
        </div>
        <div className="resource-row-details-buttons">
          <Button
            size="small"
            variant="contained"
            color="success"
            startIcon={<DownloadOutlined />}
            className="resource-row-details-buttons-btn"
            href={`/api/objects?path=${bucketName}/${name}`}
            style={{ height: 30.75, position: "relative", bottom: -5 }}
          >
            Download
          </Button>
          <Button
            size="small"
            variant="contained"
            color="warning"
            startIcon={<DeleteOutlined />}
            className="resource-row-details-buttons-btn"
            onClick={() => onDelete(object)}
          >
            Delete
          </Button>
        </div>
      </AccordionDetails>
    </Accordion>
  );
}

export default function ObjectsView({ resources, setResources }) {
  const [addDialogOpen, setAddDialogOpen] = useState(false);
  const [deleteObject, setDeleteObject] = useState(false);

  const query = new URLSearchParams(useLocation().search);
  const queriedObjectId = parseInt(query.get("object_id"));

  const objects = resources?.objects ?? [];

  const rows = objects.map((o, i) => (
    <ObjectRow
      object={o}
      expanded={o.object_id === queriedObjectId}
      onDelete={(o) => {
        setDeleteObject(o);
      }}
      key={i}
    />
  ));

  return (
    <Page>
      <div className="resource-header">
        <h2>Objects</h2>

        <Button
          variant="contained"
          color="primary"
          startIcon={<AddOutlined />}
          className="Deployments-header-button"
          size="small"
          onClick={() => setAddDialogOpen(true)}
        >
          Add Object
        </Button>
      </div>

      {rows}

      <AddObjectDialog
        open={addDialogOpen}
        onClose={() => setAddDialogOpen(false)}
        resources={resources}
        setResources={setResources}
      />
      <DeleteObjectDialog
        open={!!deleteObject}
        onClose={() => setDeleteObject(false)}
        object={deleteObject}
        resources={resources}
        setResources={setResources}
        key={deleteObject.object_id}
      />
    </Page>
  );
}
