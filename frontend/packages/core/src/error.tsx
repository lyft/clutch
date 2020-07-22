import React from "react";
import {
  Collapse as MuiCollapse,
  ExpansionPanel,
  ExpansionPanelDetails,
  ExpansionPanelSummary,
  IconButton,
  Snackbar,
  Typography,
} from "@material-ui/core";
import ExpandMoreIcon from "@material-ui/icons/ExpandMore";
import RefreshIcon from "@material-ui/icons/Refresh";
import { Alert, AlertTitle } from "@material-ui/lab";
import styled from "styled-components";

const PANEL_MESSAGE_BREAKPOINT = 150;

interface ErrorProps {
  message: string;
  retry?: () => void;
}

const Error: React.FC<ErrorProps> = ({ message, retry }) => {
  const action =
    retry !== undefined ? (
      <IconButton aria-label="retry" color="inherit" size="small" onClick={() => retry()}>
        <RefreshIcon />
      </IconButton>
    ) : null;
  return (
    <Alert severity="error" action={action}>
      {message}
    </Alert>
  );
};

const Collapse = styled(MuiCollapse)`
  margin-top: 10px;
  width: 45%;
`;

const ErrorText = styled(Typography)`
  color: rgb(97, 26, 21);
  font-size: 0.875rem;
`;

const ErrorPanel = styled(ExpansionPanel)`
  background-color: inherit;
  padding: 0px;
  width: 100%;
`;

const CompressedError = ({ title, message }) => {
  const [open, setOpen] = React.useState(message !== "");
  const [errorMsg, setErrorMsg] = React.useState("");

  React.useEffect(() => {
    if (message !== "") {
      setErrorMsg(message);
    }
    setOpen(message !== "");
  }, [message]);

  return (
    <Collapse in={open}>
      <Alert severity="error">
        <AlertTitle>{title || "Error"}</AlertTitle>
        {(errorMsg?.length || 0) > PANEL_MESSAGE_BREAKPOINT ? (
          <ErrorPanel elevation={0}>
            <ExpansionPanelSummary
              style={{ padding: "0px" }}
              expandIcon={<ExpandMoreIcon />}
              aria-controls="panel1a-content"
            >
              <ErrorText>{errorMsg.slice(0, PANEL_MESSAGE_BREAKPOINT)}</ErrorText>
            </ExpansionPanelSummary>
            <ExpansionPanelDetails style={{ padding: "0px" }}>
              <ErrorText>{errorMsg.slice(PANEL_MESSAGE_BREAKPOINT)}</ErrorText>
            </ExpansionPanelDetails>
          </ErrorPanel>
        ) : (
          errorMsg
        )}
      </Alert>
    </Collapse>
  );
};

const Warning = ({ message }) => {
  const [open, setOpen] = React.useState(true);

  return (
    <Snackbar
      open={open}
      autoHideDuration={6000}
      onClose={() => setOpen(false)}
      anchorOrigin={{ vertical: "bottom", horizontal: "right" }}
    >
      <Alert elevation={6} variant="filled" onClose={() => setOpen(false)} severity="warning">
        {message}
      </Alert>
    </Snackbar>
  );
};

export { CompressedError, Error, Warning };
