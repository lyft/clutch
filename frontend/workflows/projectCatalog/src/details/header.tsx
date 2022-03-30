import React from "react";
import { Grid, styled, Typography } from "@clutch-sh/core";

interface ProjectHeaderProps {
  name: string;
  routeTitle?: string;
  description?: string;
}

const StyledHeading = styled("div")({
  padding: "8px 0px 8px 0px",
  textTransform: "capitalize",
});

const StyledContainer = styled(Grid)({
  width: "100%",
  height: "100%",
});

const ProjectHeader = ({
  name,
  routeTitle = "Project Catalog",
  description = "",
}: ProjectHeaderProps) => (
  <StyledContainer container direction="column">
    <Grid container item direction="row" alignItems="flex-end">
      <Typography variant="body4">{routeTitle}</Typography>&nbsp;/&nbsp;
      <Typography variant="caption2">{name}</Typography>
    </Grid>
    <StyledHeading>
      <Typography variant="h2">{name}</Typography>
    </StyledHeading>
    {description.length && <Typography variant="body2">{description}</Typography>}
  </StyledContainer>
);

export default ProjectHeader;
