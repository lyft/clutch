import React from "react";
import {
  Chip,
  Grid,
  IconButton,
  styled,
  Table,
  TableCell,
  TableRow,
  Typography,
} from "@clutch-sh/core";
import { Collapse } from "@material-ui/core";
import KeyboardArrowDownIcon from "@material-ui/icons/KeyboardArrowDown";
import KeyboardArrowUpIcon from "@material-ui/icons/KeyboardArrowUp";

import { EventTime } from "../helpers";

import IncidentRow from "./incidentRow";
import NoAlertRow from "./noAlertRow";
import type { Alert } from "./types";

const StyledTableCell = styled(TableCell)({
  padding: 0,
});

const AlertRow = ({
  project,
  alerts,
  singleProject,
}: {
  project: string;
  alerts: Alert;
  singleProject: boolean;
}) => {
  const [open, setOpen] = React.useState(false);

  if (singleProject) {
    return (
      <>
        {alerts.incidents?.length &&
          alerts.incidents.map(incident => <IncidentRow incident={incident} />)}
      </>
    );
  }

  const { incidents = [] } = alerts;

  if (!incidents.length) {
    return <NoAlertRow project={project} />;
  }

  const ackedAlerts = incidents.filter(inc => inc.status === "ack").length;
  const openAlerts = incidents.filter(inc => inc.status === "open").length;
  return (
    <>
      <TableRow key={project}>
        <Typography variant="subtitle3">{project}</Typography>
        <Grid container spacing={1}>
          {openAlerts > 0 && (
            <Grid item key="openalerts">
              <Chip label={`${openAlerts} Open`} variant="error" />
            </Grid>
          )}
          {ackedAlerts > 0 && (
            <Grid item key="ackalerts">
              <Chip label={`${ackedAlerts} Ack`} variant="warn" />
            </Grid>
          )}
        </Grid>
        <EventTime date={incidents[0].created} />
        <span style={{ display: "flex", justifyContent: "flex-end" }}>
          <IconButton aria-label="expand row" size="small" onClick={() => setOpen(!open)}>
            {open ? <KeyboardArrowUpIcon /> : <KeyboardArrowDownIcon />}
          </IconButton>
        </span>
      </TableRow>
      <>
        <StyledTableCell colSpan={4}>
          <Collapse in={open} timeout="auto" unmountOnExit>
            <Table columns={["", "", "", "", ""]}>
              {incidents.map(incident => (
                <IncidentRow incident={incident} />
              ))}
            </Table>
          </Collapse>
        </StyledTableCell>
      </>
    </>
  );
};

export default AlertRow;
