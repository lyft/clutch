import React from "react";
import styled from "@emotion/styled";
import {
  Accordion as MuiAccordion,
  AccordionDetails,
  AccordionSummary,
  Collapse as MuiCollapse,
  IconButton,
  Typography,
} from "@material-ui/core";
import ExpandMoreIcon from "@material-ui/icons/ExpandMore";
import RefreshIcon from "@material-ui/icons/Refresh";

import { Alert } from "./alert";

const BREAKPOINT_LENGTH = 100;
const BREAKPOINT_REGEX = /[.\n]/g;

export interface ErrorProps {
  message: string;
  onRetry?: () => void;
}

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

const ErrorText = styled(Typography)({
  color: "rgba(13, 16, 48, 0.6)",
  fontSize: "14px",
});

const Accordion = styled(MuiAccordion)({
  backgroundColor: "inherit",
  margin: "0",
  padding: "0px",
  ":before": {
    height: "0",
  },
  "&.Mui-expanded": {
    margin: "0",
    minHeight: "fit-content",
  },
  "& .MuiAccordionSummary-root": {
    padding: "0",
    minHeight: "fit-content",
  },
});

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

export interface CompressedErrorProps {
  title?: string;
  message: string;
}

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
    <MuiCollapse in={open}>
      <Alert severity="error" title={title || "Error"}>
        {(errorMsg?.length || 0) > BREAKPOINT_LENGTH ? (
          <Accordion elevation={0}>
            <AccordionSummary expandIcon={<ExpandMoreIcon />} aria-controls="panel1a-content">
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
    </MuiCollapse>
  );
};

export { CompressedError, Error };
