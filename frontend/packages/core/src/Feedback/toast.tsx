import React from "react";
import type { SnackbarCloseReason, SnackbarProps } from "@material-ui/core";
import { Snackbar } from "@material-ui/core";

import type { AlertProps } from "./alert";
import { Alert } from "./alert";

export interface ToastProps
  extends Pick<SnackbarProps, "anchorOrigin" | "autoHideDuration">,
    Pick<AlertProps, "title" | "severity"> {
  onClose?: () => void;
}

const Toast: React.FC<ToastProps> = ({
  onClose,
  title,
  severity = "warning",
  autoHideDuration = 6000,
  anchorOrigin = { vertical: "bottom", horizontal: "right" },
  children,
}) => {
  const [open, setOpen] = React.useState(true);

  const onDismiss = () => {
    if (open && onClose !== undefined) {
      onClose();
    }
    setOpen(false);
  };

  return (
    <Snackbar
      open={open}
      autoHideDuration={autoHideDuration}
      anchorOrigin={anchorOrigin}
      onClose={(_, reason: SnackbarCloseReason) => {
        // This way it will not auto close when clicking in the window, will instead wait on the timeout or onClose
        if (reason !== "clickaway") {
          onDismiss();
        }
      }}
    >
      <Alert
        elevation={6}
        variant="filled"
        onClose={onClose ? onDismiss : null}
        severity={severity}
        title={title}
      >
        {children}
      </Alert>
    </Snackbar>
  );
};

export default Toast;
