import React from "react";
import MuiOpenInNewIcon from "@mui/icons-material/OpenInNew";
import RefreshIcon from "@mui/icons-material/Refresh";
import { alpha, IconButton, Theme } from "@mui/material";

import { Link } from "../../link";
import type { ClutchError } from "../../Network/errors";
import { isHelpDetails } from "../../Network/errors";
import styled from "../../styled";
import { Alert } from "../alert";

import ErrorDetails from "./details";

const ErrorSummaryContainer = styled("div")({
  width: "100%",
  display: "flex",
  flexDirection: "column",
});

const ErrorSummaryMessage = styled("div")({
  lineHeight: "24px",
  margin: "4px 0",
  flex: "1",
});

const ErrorSummaryLink = styled(Link)(({ theme }: { theme: Theme }) => ({
  fontSize: "14px",
  fontWeight: 400,
  color: alpha(theme.palette.secondary[900], 0.6),
  display: "flex",
  alignItems: "center",
}));

const ErrorAlert = styled(Alert)(props =>
  props["data-detailed"]
    ? {
        borderBottomLeftRadius: "unset",
        borderBottomRightRadius: "unset",
      }
    : {}
);

const OpenInNewIcon = styled(MuiOpenInNewIcon)({
  margin: "3px 8px 3px 0",
});

export interface ErrorProps {
  subject: ClutchError;
  onRetry?: () => void;
  children?: React.ReactChild | React.ReactChildren;
}

const Error = ({ subject: error, onRetry, children }: ErrorProps) => {
  const action =
    onRetry !== undefined ? (
      <IconButton aria-label="retry" color="inherit" size="small" onClick={() => onRetry()}>
        <RefreshIcon />
      </IconButton>
    ) : null;

  if (error?.details === undefined) {
    return (
      <Alert severity="error" title={error.status?.text} action={action}>
        {error.message}
        {children && children}
      </Alert>
    );
  }

  let links = [];
  const hasDetails =
    error.details?.filter(detail => {
      if (isHelpDetails(detail)) {
        links = detail?.links || [];
        return false;
      }
      return true;
    }).length > 0;

  return (
    <div>
      <ErrorAlert
        severity="error"
        title={error.status?.text}
        data-detailed={hasDetails}
        action={action}
      >
        <ErrorSummaryContainer>
          <ErrorSummaryMessage>{error.message}</ErrorSummaryMessage>
          {links.map(link => (
            <ErrorSummaryLink key={link.link} href={link.link}>
              <OpenInNewIcon fontSize="small" />
              {link.description}
            </ErrorSummaryLink>
          ))}
        </ErrorSummaryContainer>
      </ErrorAlert>
      {hasDetails && <ErrorDetails error={error} />}
    </div>
  );
};

export default Error;
