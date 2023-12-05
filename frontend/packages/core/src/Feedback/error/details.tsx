import React from "react";
import ChevronRightIcon from "@mui/icons-material/ChevronRight";
import KeyboardArrowDownIcon from "@mui/icons-material/KeyboardArrowDown";
import {
  Accordion as MuiAccordion,
  AccordionDetails as MuiAccordionDetails,
  AccordionSummary as MuiAccordionSummary,
  Button,
  Grid,
  useControlled,
  useTheme,
} from "@mui/material";

import type { ClutchError } from "../../Network/errors";
import { isClutchErrorDetails } from "../../Network/errors";
import { grpcCodeToText } from "../../Network/grpc";
import styled from "../../styled";

import ErrorDetailsDialog from "./dialog";

const ERROR_DETAILS_RENDER_MAX = 4;

const ErrorDetailDivider = styled("div")(() => {
  const theme = useTheme();
  return {
    background: `linear-gradient(to right, ${theme.palette.error[600]} 8px, ${theme.palette.error[200]}44 0%)`,
    height: "1px",
    width: "100%",
  };
});

const Accordion = styled(MuiAccordion)({
  "&.MuiAccordion-root.Mui-expanded": {
    margin: "0px",
  },
  ":before": {
    height: "0",
  },
});

const AccordionSummary = styled(MuiAccordionSummary)<{ $expanded: boolean }>(
  () => {
    const theme = useTheme();
    return {
      background: `linear-gradient(to right, ${theme.palette.error[600]} 8px, ${theme.palette.error[100]} 0%)`,
      color: theme.palette.secondary[900],
      fontSize: "14px",
      fontWeight: 400,
      padding: "12px 16px 12px 24px",
      minHeight: "fit-content",
      "& .MuiAccordionSummary-content": {
        margin: "0",
        alignItems: "center",
      },
      "&.MuiAccordionSummary-root.Mui-expanded": {
        minHeight: "unset",
      },
    };
  },
  props => ({
    borderBottomLeftRadius: props.$expanded ? "0" : "8px",
    borderBottomRightRadius: props.$expanded ? "0" : "8px",
  })
);

const AccordionDetails = styled(MuiAccordionDetails)(() => {
  const theme = useTheme();
  return {
    background: `linear-gradient(to right, ${theme.palette.error[600]} 8px, ${theme.palette.common.white} 0%)`,
    padding: "0",
    paddingLeft: "8px",
    borderBottomLeftRadius: "8px",
    borderBottomRightRadius: "8px",
    display: "flex",
    flexDirection: "column",
  };
});

const ListItem = styled("li")(() => {
  const theme = useTheme();
  return {
    "::marker": {
      color: `${theme.palette.secondary[700]}66`,
    },
    padding: "2px 0",
  };
});

const ErrorDetailContainer = styled("div")(() => {
  const theme = useTheme();
  return {
    width: "100%",
    border: `1px solid ${theme.palette.secondary[200]}`,
    padding: "16px 16px 16px 24px",
    borderBottomRightRadius: "8px",
    borderTop: "unset",
  };
});

const ErrorDetailText = styled("div")(() => {
  const theme = useTheme();
  return {
    color: `${theme.palette.secondary[700]}66`,
    fontSize: "14px",
    lineHeight: "24px",
  };
});

const DialogButton = styled(Button)(() => {
  const theme = useTheme();
  return {
    color: theme.palette.primary[600],
    fontWeight: 700,
    fontSize: "14px",
    padding: "9px 32px",
  };
});

interface ErrorDetailsProps {
  error: ClutchError;
}

const ErrorDetails = ({ error }: ErrorDetailsProps) => {
  const theme = useTheme();
  const [detailsOpen, setDetailsOpen] = React.useState(false);
  const [expanded, setExpanded] = useControlled({
    controlled: undefined,
    default: false,
    name: "Error Accordion",
    state: "expanded",
  });

  React.useEffect(() => {
    setDetailsOpen(false);
  }, [error]);

  const hasWrappedErrorDetails =
    error.details.filter(detail => isClutchErrorDetails(detail)).length > 0;

  const summaryIconStyle = { marginRight: "8px" };

  return (
    <>
      <ErrorDetailDivider />
      <Accordion elevation={0} expanded={expanded}>
        <AccordionSummary
          aria-controls="panel1a-content"
          $expanded={expanded}
          onClick={() => setExpanded(!expanded)}
        >
          {!expanded ? (
            <>
              <ChevronRightIcon style={summaryIconStyle} /> Show more
            </>
          ) : (
            <>
              <KeyboardArrowDownIcon style={summaryIconStyle} /> Show less
            </>
          )}
        </AccordionSummary>
        <AccordionDetails>
          <ErrorDetailContainer>
            {hasWrappedErrorDetails && (
              <div>
                <ErrorDetailText style={{ color: theme.palette.secondary[900] }}>
                  The following errors were encountered:
                </ErrorDetailText>
                <ul style={{ paddingLeft: "16px", margin: "4px 0" }}>
                  {error.details.map(detail => {
                    // Only render Clutch Error wrapped details errors here
                    if (isClutchErrorDetails(detail)) {
                      const renderItems = detail.wrapped.slice(0, ERROR_DETAILS_RENDER_MAX);
                      const remainingItems = detail.wrapped.length - ERROR_DETAILS_RENDER_MAX;
                      return (
                        <>
                          {renderItems.map((wrapped, idx) => {
                            // TODO: This color should be colored according to status code
                            const color = theme.palette.secondary[600];
                            return (
                              // eslint-disable-next-line react/no-array-index-key
                              <ListItem key={`${idx}-${wrapped.message}`}>
                                <ErrorDetailText>
                                  <span style={{ fontWeight: 500, color }}>
                                    {grpcCodeToText(wrapped.code)}&nbsp;
                                  </span>
                                  {wrapped.message}
                                </ErrorDetailText>
                              </ListItem>
                            );
                          })}
                          {remainingItems > 0 && (
                            <ErrorDetailText style={{ margin: "2px 0" }}>
                              and {remainingItems} more...
                            </ErrorDetailText>
                          )}
                        </>
                      );
                    }
                    return null;
                  })}
                </ul>
              </div>
            )}
            <Grid container justifyContent="flex-end">
              <DialogButton onClick={() => setDetailsOpen(true)}>More Details</DialogButton>
            </Grid>
          </ErrorDetailContainer>
        </AccordionDetails>
      </Accordion>
      <ErrorDetailsDialog error={error} open={detailsOpen} onClose={() => setDetailsOpen(false)} />
    </>
  );
};

export default ErrorDetails;
