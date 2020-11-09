import React from "react";
import {
  AccordionDetails,
  AccordionSummary,
  Collapse as MuiCollapse,
  ExpansionPanel,
  IconButton,
  Typography,
} from "@material-ui/core";
import ExpandMoreIcon from "@material-ui/icons/ExpandMore";
import RefreshIcon from "@material-ui/icons/Refresh";
import { Alert as MuiAlert, AlertTitle } from "@material-ui/lab";
import styled from "styled-components";

const PANEL_MESSAGE_BREAKPOINT = 150;

export interface ErrorProps {
  message: string;
  onRetry?: () => void;
}

const Alert = styled(MuiAlert)`
  margin: 5px;
`;

const Error: React.FC<ErrorProps> = ({ message, onRetry }) => {
  const action =
    onRetry !== undefined ? (
      <IconButton aria-label="retry" color="inherit" size="small" onClick={() => onRetry()}>
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
            <AccordionSummary
              style={{ padding: "0px" }}
              expandIcon={<ExpandMoreIcon />}
              aria-controls="panel1a-content"
            >
              <ErrorText>{errorMsg.slice(0, PANEL_MESSAGE_BREAKPOINT)}</ErrorText>
            </AccordionSummary>
            <AccordionDetails style={{ padding: "0px" }}>
              <ErrorText>{errorMsg.slice(PANEL_MESSAGE_BREAKPOINT)}</ErrorText>
            </AccordionDetails>
          </ErrorPanel>
        ) : (
          errorMsg
        )}
      </Alert>
    </Collapse>
  );
};

export { CompressedError, Error };
