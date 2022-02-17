import React from "react";
import { Snackbar } from "@material-ui/core";

import type { AlertProps } from "./alert";
import { Alert } from "./alert";

export interface ToastProps extends AlertProps {
  duration?: number;
  onClose?: () => void;
}

const Toast: React.FC<ToastProps> = ({
  onClose,
  severity = "warning",
  duration = 6000,
  ...props
}) => {
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
      autoHideDuration={duration}
      onExiting={onDismiss}
      onClose={(_, reason) => {
        // This way it will not auto close when clicking in the window, will instead wait on the timeout or onClose
        if (reason !== "clickaway") {
          setOpen(false);
        }
      }}
      anchorOrigin={{ vertical: "bottom", horizontal: "right" }}
    >
      <Alert
        elevation={6}
        variant="filled"
        onClose={onClose ? onDismiss : null}
        severity={severity}
        {...props}
      />
    </Snackbar>
  );
};

export default Toast;
