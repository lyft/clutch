import React from "react";
import { Snackbar } from "@material-ui/core";

import type { AlertProps } from "./alert";
import { Alert } from "./alert";

export interface WarningProps extends AlertProps {
  message: any;
  onClose?: () => void;
}

const Warning: React.FC<WarningProps> = ({ message, onClose, severity = "warning", ...props }) => {
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
      <Alert elevation={6} variant="filled" onClose={onDismiss} severity={severity} {...props}>
        {message}
      </Alert>
    </Snackbar>
  );
};

export default Warning;
