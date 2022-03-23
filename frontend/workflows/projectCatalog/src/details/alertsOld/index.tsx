import React from "react";
import type { ClutchError } from "@clutch-sh/core";
import { styled, Table, TableRow, Toast, Typography } from "@clutch-sh/core";
import { Card } from "@clutch-sh/project-selector";

import AlertEventIcon from "../../assets/AlertEvent";
import { DefaultSummaryTitle, EventTime, setMilliseconds } from "../helpers";

import AlertRow from "./alertRow";
import type { AlertInfo } from "./types";

export interface ProjectAlertsProps {
  alerts: AlertInfo;
  loading?: boolean;
  error?: ClutchError | undefined;
  singleProject?: boolean;
}

const StyledCard = styled(Card)({
  width: "100%",
  height: "100%",
});

const ProjectAlertsCard = ({
  alerts,
  error,
  loading = false,
  singleProject = true,
}: ProjectAlertsProps) => {
  const [projects, setProjects] = React.useState<string[]>([]);

  React.useEffect(() => {
    if (alerts && alerts.projectAlerts) {
      setProjects(Object.keys(alerts.projectAlerts));
    }
  }, [alerts]);

  return (
    <>
      {error && <Toast>Failed to fetch Deploys</Toast>}
      <StyledCard
        avatar={<AlertEventIcon />}
        title="Alerts"
        error={error}
        isLoading={loading}
        summary={[
          {
            title:
              alerts?.lastAlert > 0 ? (
                <Typography variant="subtitle2">
                  <EventTime date={setMilliseconds(alerts.lastAlert)} />
                </Typography>
              ) : (
                <DefaultSummaryTitle />
              ),
            subheader: "Last Alert",
          },
          {
            title: alerts?.open ? (
              <Typography variant="subtitle2" color="#DB3615">
                {alerts.open}
              </Typography>
            ) : (
              <DefaultSummaryTitle />
            ),
            subheader: "Open",
          },
          {
            title: alerts?.acknowledged ? (
              <Typography variant="subtitle2" color="#D87313">
                {alerts.acknowledged}
              </Typography>
            ) : (
              <DefaultSummaryTitle />
            ),
            subheader: "Acknowledged",
          },
        ]}
      >
        <Table columns={["", "", "", ""]}>
          {projects.length ? (
            projects.map(pkey => (
              <AlertRow
                alerts={alerts.projectAlerts[pkey]}
                project={pkey}
                singleProject={singleProject}
              />
            ))
          ) : (
            <TableRow>
              <div>No alerts found for selected project(s)</div>
            </TableRow>
          )}
        </Table>
      </StyledCard>
    </>
  );
};

export default ProjectAlertsCard;
