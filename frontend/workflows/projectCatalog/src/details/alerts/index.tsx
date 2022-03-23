import React from "react";
import { Button, ClutchError, Grid } from "@clutch-sh/core";

import AlertEventIcon from "../../assets/AlertEvent";
import type { BaseProjectCardProps } from "../card";
import ProjectCard, { LastEvent, StyledLink, StyledRow } from "../card";

import OnCallRow from "./onCallRow";
import SummaryRow from "./summaryRow";
import type { ProjectAlerts } from "./types";

interface ProjectAlertsProps {
  data: ProjectAlerts;
  error?: ClutchError | undefined;
  loading?: boolean;
}

const ProjectAlertsCard = ({ data, loading, error }: ProjectAlertsProps) => {
  const titleData: BaseProjectCardProps = {
    text: data?.title ?? "Alerts",
    icon: <AlertEventIcon />,
    endAdornment: <LastEvent time={data?.lastAlert} />,
  };

  return (
    <ProjectCard loading={loading} error={error} {...titleData}>
      {data?.summary && (
        <StyledRow container item direction="row" justify="space-evenly">
          <SummaryRow {...data.summary} />
        </StyledRow>
      )}
      {data?.onCall && (
        <StyledRow container item direction="column" spacing={1}>
          <OnCallRow {...data.onCall} />
        </StyledRow>
      )}
      {data?.create && (
        <Grid container item direction="column" alignItems="flex-end">
          <Grid item xs={6}>
            <StyledLink href={data.create.url}>
              <Button text={data.create.text} />
            </StyledLink>
          </Grid>
        </Grid>
      )}
    </ProjectCard>
  );
};

export default ProjectAlertsCard;
