import React from "react";
import styled from "@emotion/styled";
import { IconButton } from "@material-ui/core";
import MuiOpenInNewIcon from "@material-ui/icons/OpenInNew";
import RefreshIcon from "@material-ui/icons/Refresh";

import { StyledLink } from "../../link";
import type { ClutchError } from "../../Network/errors";
import { isHelpDetails } from "../../Network/errors";
import { Alert } from "../alert";

import ErrorDetails from "./details";

const ErrorSummaryContainer = styled.div({
  width: "100%",
  display: "flex",
  flexDirection: "column",
});

const ErrorSummaryMessage = styled.div({
  height: "24px",
  flex: "1",
});

const ErrorSummaryLink = styled(StyledLink)({
  fontSize: "14px",
  fontWeight: 400,
  color: "rgb(13,16,48,0.6)",
  display: "flex",
  alignItems: "center",
});

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
}

const Error = ({ subject: error, onRetry }: ErrorProps) => {
  const action =
    onRetry !== undefined ? (
      <IconButton aria-label="retry" color="inherit" size="small" onClick={() => onRetry()}>
        <RefreshIcon />
      </IconButton>
    ) : null;

  if (error?.details === undefined) {
    return (
      <Alert severity="error" title={error.status.text} action={action}>
        {error.message}
      </Alert>
    );
  }

  let links = [];
  const hasDetails =
    error.details.filter(detail => {
      if (isHelpDetails(detail)) {
        links = detail?.links || [];
        return false;
      }
      return true;
    }).length > 0;

  return (
    <>
      <ErrorAlert
        severity="error"
        title={error.status.text}
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
    </>
  );
};

export default Error;
