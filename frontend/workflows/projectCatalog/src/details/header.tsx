import React from "react";
import { Grid, styled, Typography } from "@clutch-sh/core";
import { capitalize } from "lodash";

interface ProjectHeaderProps {
  name: string;
  routeTitle?: string;
  description?: string;
}

const StyledHeading = styled("div")({
  padding: "8px 0px 8px 0px",
});

const ProjectHeader = ({
  name,
  routeTitle = "Project Catalog",
  description = "",
}: ProjectHeaderProps) => (
  <>
    <Grid container direction="column" style={{ width: "100%", height: "100%" }}>
      <Grid container item direction="row" alignItems="flex-end">
        <Typography variant="body4">{routeTitle}</Typography>&nbsp;/&nbsp;
        <Typography variant="caption2">{name}</Typography>
      </Grid>
      <StyledHeading>
        <Typography variant="h2">{capitalize(name)}</Typography>
      </StyledHeading>
      {description.length && <Typography variant="body2">{description}</Typography>}
    </Grid>
  </>
);

export default ProjectHeader;
