import React from "react";
import { Button, Grid } from "@clutch-sh/core";

import AlertEventIcon from "../../assets/AlertEvent";
import ProjectCard, { LastEvent, StyledLink, StyledRow } from "../card";

import OnCallRow from "./onCallRow";
import SummaryRow from "./summaryRow";
import type { ProjectAlerts } from "./types";

const ProjectAlertsCard = ({
  create,
  lastAlert,
  onCall,
  summary,
  title = "Alerts",
}: ProjectAlerts) => {
  const titleData = {
    text: title,
    icon: <AlertEventIcon />,
    endAdornment: <LastEvent time={lastAlert} />,
  };

  return (
    <ProjectCard {...titleData}>
      {summary && (
        <StyledRow container item direction="row" justify="space-evenly">
          <SummaryRow {...summary} />
        </StyledRow>
      )}
      {onCall && (
        <StyledRow container item direction="column" spacing={1}>
          <OnCallRow {...onCall} />
        </StyledRow>
      )}
      {create && (
        <Grid container item direction="column" alignItems="flex-end">
          <Grid item xs={6}>
            <StyledLink href={create.url}>
              <Button text={create.text} />
            </StyledLink>
          </Grid>
        </Grid>
      )}
    </ProjectCard>
  );
};

export default ProjectAlertsCard;
