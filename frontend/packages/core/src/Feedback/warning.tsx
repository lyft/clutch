import React from "react";
import { Snackbar } from "@material-ui/core";
import type { AlertProps as MuiAlertProps } from "@material-ui/lab";

import { Alert } from "./alert";

export interface WarningProps extends Pick<MuiAlertProps, "severity" | "action" | "onClose"> {
  message: any;
  title?: React.ReactNode;
  onClose?: () => void;
}

const Warning: React.FC<WarningProps> = ({ message, title, onClose, severity = "warning" }) => {
  const [open, setOpen] = React.useState(true);

  const onDismiss = () => {
    if (onClose !== undefined) {
      onClose();
    }
    setOpen(false);
  };

  return (
    <Snackbar
      open={open}
      autoHideDuration={6000}
      onExiting={onDismiss}
      onClose={() => setOpen(false)}
      anchorOrigin={{ vertical: "bottom", horizontal: "right" }}
    >
      <Alert elevation={6} variant="filled" onClose={onDismiss} title={title} severity={severity}>
        {message}
      </Alert>
    </Snackbar>
  );
};

export default Warning;
