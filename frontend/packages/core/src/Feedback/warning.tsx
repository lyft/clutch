import React from "react";
import styled from "@emotion/styled";
import { Snackbar } from "@material-ui/core";
import { Alert as MuiAlert } from "@material-ui/lab";

const Alert = styled(MuiAlert)({
  margin: "5px",
});

export interface WarningProps {
  message: any;
  onClose?: () => void;
}

const Warning: React.FC<WarningProps> = ({ message, onClose }) => {
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
      <Alert elevation={6} variant="filled" onClose={onDismiss} severity="warning">
        {message}
      </Alert>
    </Snackbar>
  );
};

export default Warning;
