import React from "react";
import {
  Accordion as MuiAccordion,
  AccordionDetails,
  AccordionSummary,
  Collapse as MuiCollapse,
  IconButton,
  Snackbar,
  Typography,
} from "@material-ui/core";
import ExpandMoreIcon from "@material-ui/icons/ExpandMore";
import RefreshIcon from "@material-ui/icons/Refresh";
import { Alert as MuiAlert, AlertTitle } from "@material-ui/lab";
import styled from "styled-components";

const BREAKPOINT_LENGTH = 100;
const BREAKPOINT_REGEX = /[.\n]/g;

interface ErrorProps {
  message: string;
  retry?: () => void;
}

const Alert = styled(MuiAlert)`
  margin: 5px;
  min-width: fit-content;
  max-width: 45vw;
`;

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

const Accordion = styled(MuiAccordion)`
  background-color: inherit;
  padding: 0px;
`;

export interface CompressedErrorProps {
  title?: string;
  message: string;
}

const findBreakpoint = (message: string): number => {
  // n.b. if no breakpoint this will return -1 and become 0.
  let breakpoint = message.search(BREAKPOINT_REGEX) + 1;
  if (breakpoint === 0) {
    let letterCount = 0;
    message.split(" ").some((word: string): boolean => {
      // n.b. add 1 to account for the space
      const newCount = letterCount + word.length + 1;
      if (newCount > BREAKPOINT_LENGTH) {
        return true;
      }
      letterCount = newCount;
      return false;
    });
    breakpoint = letterCount || BREAKPOINT_LENGTH;
  }
  return breakpoint;
};

const CompressedError: React.FC<CompressedErrorProps> = ({ title, message }) => {
  const [open, setOpen] = React.useState(message !== "");
  const [errorMsg, setErrorMsg] = React.useState("");

  React.useEffect(() => {
    if (message !== "") {
      setErrorMsg(message);
    }
    setOpen(message !== "");
  }, [message]);

  const breakpoint = findBreakpoint(errorMsg);
  return (
    <Collapse in={open}>
      <Alert severity="error">
        <AlertTitle>{title || "Error"}</AlertTitle>
        {(errorMsg?.length || 0) > BREAKPOINT_LENGTH ? (
          <Accordion elevation={0}>
            <AccordionSummary
              style={{ padding: "0px" }}
              expandIcon={<ExpandMoreIcon />}
              aria-controls="panel1a-content"
            >
              <ErrorText>{errorMsg.slice(0, breakpoint)}</ErrorText>
            </AccordionSummary>
            <AccordionDetails style={{ padding: "0px", overflowWrap: "anywhere" }}>
              <ErrorText>{errorMsg.slice(breakpoint)}</ErrorText>
            </AccordionDetails>
          </Accordion>
        ) : (
          errorMsg
        )}
      </Alert>
    </Collapse>
  );
};

interface WarningProps {
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
      onExit={onDismiss}
      onClose={() => setOpen(false)}
      anchorOrigin={{ vertical: "bottom", horizontal: "right" }}
    >
      <Alert elevation={6} variant="filled" onClose={onDismiss} severity="warning">
        {message}
      </Alert>
    </Snackbar>
  );
};

export { CompressedError, Error, Warning };
