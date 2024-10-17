import React from "react";
import { Grid, styled, Typography } from "@clutch-sh/core";

export interface ProjectHeaderProps {
  title?: string;
  description?: string;
}

const StyledContainer = styled(Grid)({
  width: "100%",
  height: "100%",
});

const ProjectHeader = ({ title, description = "" }: ProjectHeaderProps) => (
  <StyledContainer container direction="column">
    <Grid item>
      <Typography variant="h2" textTransform="capitalize">
        {title}
      </Typography>
    </Grid>
    {description.length > 0 && (
      <Grid item>
        <Typography variant="body2">{description}</Typography>
      </Grid>
    )}
  </StyledContainer>
);

export default ProjectHeader;
