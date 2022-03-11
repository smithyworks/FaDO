import { EditOutlined, ExpandMore, InfoOutlined } from "@mui/icons-material";
import {
  Accordion,
  AccordionDetails,
  AccordionSummary,
  Button,
  Checkbox,
  Collapse,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Divider,
  FormControl,
  InputLabel,
  MenuItem,
  Paper,
  Select,
  TextField,
  Typography,
} from "@mui/material";
import React, { useState } from "react";
import api from "../../api";
import Page from "../../components/Page";
import ResourceDialog from "../../components/ResourceDialog";

function EditDefaultsDialog({ open, onClose, settings, setResources }) {
  const [headerName, setHeaderName] = useState(settings?.match_header ?? "");
  const [policy, setPolicy] = useState(settings?.policy ?? "");

  function onOk() {
    settings.policy = policy.trim();
    settings.match_header = headerName.trim();

    api
      .setLBDefaults(settings)
      .then((res) => setResources(res.data))
      .catch((err) => console.log(err))
      .finally(() => console.log("send defaults request finished"));

    setTimeout(() => window.location.reload(), 1000);
    setHeaderName("");
    setPolicy("");

    onClose();
  }

  return (
    <ResourceDialog title="Edit default settings." open={open} onClose={onClose} onOk={onOk}>
      <FormControl fullWidth variant="standard">
        <InputLabel id="add-default_lb_policy-select-label">Load Balancing Selection Policy</InputLabel>
        <Select
          labelId="add-default_lb_policy-select-label"
          id="add-default_lb_policy-select"
          label="Load Balancing Selection Policy"
          fullWidth
          value={policy}
          onChange={(e) => setPolicy(e.target.value)}
        >
          <MenuItem value="round_robin">round_robin</MenuItem>
          <MenuItem value="random">random</MenuItem>
          <MenuItem value="first">first</MenuItem>
          <MenuItem value="least_conn">least_conn</MenuItem>
          <MenuItem value="ip_hash">ip_hash</MenuItem>
          <MenuItem value="uri_hash">uri_hash</MenuItem>
        </Select>
      </FormControl>
      <TextField
        size="small"
        label="Match Header Name"
        fullWidth
        margin="normal"
        variant="standard"
        value={headerName}
        onChange={(e) => setHeaderName(e.target.value)}
      />
    </ResourceDialog>
  );
}

function OverrideRouteDialog({ open, onClose, data, setResources }) {
  let { bucketName, policy, upstreams, overrides } = data;
  if (!overrides) overrides = {};

  const [overridden, setOverridden] = useState(!!overrides[bucketName]);

  const [urls, setUrls] = useState(upstreams?.join(", ") ?? "");
  const [newPolicy, setPolicy] = useState(policy ?? "");

  function onOk() {
    if (overridden) {
      overrides[bucketName] = {
        policy: newPolicy,
        upstreams: urls.split(",").map((url) => url.trim()),
      };
    } else {
      delete overrides[bucketName];
    }

    api
      .setLBOverrides({ route_overrides: overrides })
      .then((res) => setResources(res.data))
      .catch((err) => console.log(err));

    setTimeout(() => window.location.reload(), 1000);
    setOverridden(false);
    setUrls("");
    setPolicy("");

    onClose();
  }

  return (
    <ResourceDialog title={`Configure route bucket '${bucketName}'.`} open={open} onClose={onClose} onOk={onOk}>
      <Typography>
        <Checkbox
          checked={overridden}
          onChange={(e) => setOverridden(e.target.checked)}
          style={{ margin: "-5px -5px 0 -10px" }}
        />{" "}
        Manually configure route.
      </Typography>
      <Collapse in={overridden}>
        <TextField
          size="small"
          label="Upstream URLs (Comma Separated List)"
          fullWidth
          margin="normal"
          variant="standard"
          value={urls}
          onChange={(e) => setUrls(e.target.value)}
        />
        <FormControl fullWidth variant="standard">
          <InputLabel id="add-default_lb_policy-select-label">Load Balancing Selection Policy</InputLabel>
          <Select
            labelId="add-default_lb_policy-select-label"
            id="add-default_lb_policy-select"
            label="Load Balancing Selection Policy"
            fullWidth
            value={newPolicy}
            onChange={(e) => setPolicy(e.target.value)}
          >
            <MenuItem value="round_robin">round_robin</MenuItem>
            <MenuItem value="random">random</MenuItem>
            <MenuItem value="first">first</MenuItem>
            <MenuItem value="least_conn">least_conn</MenuItem>
            <MenuItem value="ip_hash">ip_hash</MenuItem>
            <MenuItem value="uri_hash">uri_hash</MenuItem>
          </Select>
        </FormControl>
      </Collapse>
    </ResourceDialog>
  );
}

function RouteRow({ route, overrides, settings, onConfigure }) {
  const headerName = settings?.match_header_name ?? "X-FaDO-Bucket";
  const bucketName = route?.bucket_name ?? "-";
  const policy = route?.policy ?? "-";
  const upstreamLines = (route?.upstreams ?? []).map((us, i) => (
    <div className="resource-row-details-item" key={i}>
      {us}
    </div>
  ));

  const overridden = !!overrides[bucketName];

  return (
    <Accordion>
      <AccordionSummary expandIcon={<ExpandMore />}>
        <div className="resource-row-summary">
          <div className="resource-row-summary-title">
            {bucketName}
            {overridden && (
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
                  position: "relative",
                  top: -1,
                }}
              >
                OVERRIDDEN
              </Typography>
            )}
          </div>
        </div>
      </AccordionSummary>
      <AccordionDetails>
        <div className="resource-row-details">
          <Divider />
          <div className="resource-row-details-title">Header Matcher:</div>
          <div className="resource-row-details-item">
            {headerName}: {bucketName}
          </div>
          <div className="resource-row-details-title">Load Balancing Policy:</div>
          <div className="resource-row-details-item">{policy}</div>
          <div className="resource-row-details-title">Upstreams:</div>
          {upstreamLines}

          <div className="resource-row-details-buttons">
            <Button
              size="small"
              variant="contained"
              color="secondary"
              startIcon={<EditOutlined />}
              className="resource-row-details-buttons-btn"
              onClick={() => onConfigure({ bucketName, policy, upstreams: route?.upstreams ?? [], overrides })}
            >
              Configure
            </Button>
          </div>
        </div>
      </AccordionDetails>
    </Accordion>
  );
}

export default function LoadBalancingView({ resources, setResources }) {
  const [infoOpen, setInfoOpen] = useState(false);
  const [editDefaults, setEditDefaults] = useState(false);
  const [configureRoute, setConfigureRoute] = useState(false);

  const caddyConfig = resources?.load_balancer_config ?? {};
  const settings = resources?.load_balancer_settings ?? {};
  const routes = resources?.load_balancer_routes ?? {};
  const routeOverrides = resources?.load_balancer_route_overrides ?? {};
  const host = settings?.host;
  const port = settings?.port;

  const listenURL = `${host}:${port}`;
  const lbConfigPretty = JSON.stringify(caddyConfig, null, 2);

  const rows = Object.values(routes).map((r, i) => (
    <RouteRow
      route={r}
      onConfigure={(d) => setConfigureRoute(d)}
      settings={settings}
      overrides={routeOverrides}
      key={i}
    />
  ));

  return (
    <Page>
      <div className="resource-header">
        <h2>Load Balancing</h2>

        <Button
          variant="contained"
          color="info"
          startIcon={<InfoOutlined />}
          className="Deployments-header-button"
          onClick={() => setInfoOpen(true)}
          size="small"
        >
          View Raw Configuration
        </Button>
      </div>

      <Paper style={{ padding: 16 }}>
        <Typography style={{ padding: "10px 10px 20px 10px", fontSize: "1.1rem", fontWeight: 600 }}>
          Global Configuration
        </Typography>
        <Divider />
        <Typography style={{ padding: "10px" }}>Ingress:</Typography>
        <Typography style={{ padding: "2px 20px" }}>{listenURL}</Typography>
        <Typography style={{ padding: "10px" }}>Default Load Balancing Policy:</Typography>
        <Typography style={{ padding: "2px 20px" }}>{settings?.policy ?? "round_robin"}</Typography>
        <Typography style={{ padding: "10px" }}>Match Header Name:</Typography>
        <Typography style={{ padding: "2px 20px" }}>{settings?.match_header ?? "X-FaDO-Bucket"}</Typography>
        <div className="resource-row-details-buttons">
          <Button
            size="small"
            variant="contained"
            color="secondary"
            startIcon={<EditOutlined />}
            className="resource-row-details-buttons-btn"
            onClick={() => setEditDefaults(true)}
          >
            Edit
          </Button>
        </div>
      </Paper>

      <h2 style={{ margin: "30px 0 20px 0" }}>Bucket Routes</h2>

      {rows}

      <Dialog open={infoOpen} onClose={() => setInfoOpen(false)}>
        <DialogTitle>Load Balancing Server Configuration</DialogTitle>
        <DialogContent style={{ minWidth: 600 }}>
          <pre>{lbConfigPretty}</pre>
        </DialogContent>
        <DialogActions>
          <Button size="small" variant="outlined" color="secondary" onClick={() => setInfoOpen(false)}>
            Close
          </Button>
        </DialogActions>
      </Dialog>

      <EditDefaultsDialog
        open={editDefaults}
        settings={settings}
        setResources={setResources}
        onClose={() => setEditDefaults(false)}
        key={JSON.stringify(settings)}
      />

      <OverrideRouteDialog
        open={!!configureRoute}
        data={configureRoute}
        setResources={setResources}
        onClose={() => setConfigureRoute(false)}
        key={configureRoute?.bucketName}
      />
    </Page>
  );
}
