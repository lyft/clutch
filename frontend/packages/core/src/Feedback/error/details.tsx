import React from "react";
import styled from "@emotion/styled";
import {
  Accordion as MuiAccordion,
  AccordionDetails as MuiAccordionDetails,
  AccordionSummary as MuiAccordionSummary,
  Button,
  Grid,
  useControlled,
} from "@material-ui/core";
import ChevronRightIcon from "@material-ui/icons/ChevronRight";
import KeyboardArrowDownIcon from "@material-ui/icons/KeyboardArrowDown";

import type { ClutchError } from "../../Network/errors";
import { isClutchErrorDetails } from "../../Network/errors";
import { grpcCodeToText } from "../../Network/grpc";

import ErrorDetailsDialog from "./dialog";

const ERROR_DETAILS_RENDER_MAX = 4;

const ErrorDetailDivider = styled.div({
  background: "linear-gradient(to right, #DB3615 8px, rgba(219, 54, 21, 0.4) 0%)",
  height: "1px",
  width: "100%",
});

const Accordion = styled(MuiAccordion)({
  ":before": {
    height: "0",
  },
});

const AccordionSummary = styled(MuiAccordionSummary)(
  {
    background: "linear-gradient(to right, #DB3615 8px, #FDE9E7 0%)",
    color: "#0D1030",
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
  },
  props => ({
    borderBottomLeftRadius: props["data-expanded"] ? "0" : "8px",
    borderBottomRightRadius: props["data-expanded"] ? "0" : "8px",
  })
);

const AccordionDetails = styled(MuiAccordionDetails)({
  background: "linear-gradient(to right, #DB3615 8px, #FFFFFF 0%)",
  padding: "0",
  paddingLeft: "8px",
  borderBottomLeftRadius: "8px",
  borderBottomRightRadius: "8px",
  display: "flex",
  flexDirection: "column",
  // boxShadow: "0px 4px 6px rgba(53, 72, 212, 0.2)",
});

const ListItem = styled.li({
  "::marker": {
    color: "rgba(13, 16, 48, 0.6)",
  },
  margin: "2px 0",
});

const ErrorDetailContainer = styled.div({
  width: "100%",
  border: "1px solid #E7E7EA",
  padding: "16px 16px 16px 24px",
  borderBottomRightRadius: "8px",
  borderTop: "unset",
});

const ErrorDetailText = styled.div({
  color: "rgba(13, 16, 48, 0.6)",
  fontSize: "14px",
});

const DialogButton = styled(Button)({
  color: "#3548D4",
  fontWeight: 700,
  fontSize: "14px",
  padding: "9px 32px",
});

interface ErrorDetailsProps {
  error: ClutchError;
}

const ErrorDetails = ({ error }: ErrorDetailsProps) => {
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
          data-expanded={expanded}
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
                <ErrorDetailText style={{ color: "#0D1030" }}>
                  The following errors were encountered:
                </ErrorDetailText>
                <ul>
                  {error.details.map(detail => {
                    // Only render Clutch Error wrapped details errors here
                    if (isClutchErrorDetails(detail)) {
                      const renderItems = detail.wrapped.slice(0, ERROR_DETAILS_RENDER_MAX);
                      const remainingItems = detail.wrapped.length - ERROR_DETAILS_RENDER_MAX;
                      return (
                        <>
                          {renderItems.map((wrapped, idx) => {
                            // TODO: This color should be colored according to status code
                            const color = "#DB3615";
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
                            <ErrorDetailText>and {remainingItems} more...</ErrorDetailText>
                          )}
                        </>
                      );
                    }
                    return null;
                  })}
                </ul>
              </div>
            )}
            <Grid container justify="flex-end">
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
