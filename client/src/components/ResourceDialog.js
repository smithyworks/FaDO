import * as React from "react";
import { Button, Dialog, DialogActions, DialogTitle } from "@mui/material";

import "./ResourceDialog.css";

export default function ResourceDialog({ open, onClose, onOk, title, children }) {
  return (
    <Dialog
      open={open}
      onClose={onClose}
      aria-labelledby="alert-dialog-title"
      aria-describedby="alert-dialog-description"
    >
      <DialogTitle id="alert-dialog-title">{title}</DialogTitle>
      <div className="ResourceDialog-content">{children}</div>
      <DialogActions>
        <Button variant="contained" color="primary" size="small" onClick={onOk}>
          OK
        </Button>
        <Button variant="outlined" color="secondary" size="small" onClick={onClose} autoFocus>
          Cancel
        </Button>
      </DialogActions>
    </Dialog>
  );
}
