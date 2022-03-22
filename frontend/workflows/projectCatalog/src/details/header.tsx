import React from "react";
import { Grid, styled, Typography } from "@clutch-sh/core";
import { capitalize } from "lodash";

interface ProjectHeaderProps {
  name: string;
  description: string;
}

const TextLink = styled("a")({
  textDecoration: "none",
  color: "unset",
});

const ProjectHeader = ({ name, description = "" }: ProjectHeaderProps) => (
  <>
    <Grid container direction="column" style={{ width: "100%", height: "100%" }}>
      <Grid container item direction="row" alignItems="flex-end">
        <Typography variant="body4">
          <TextLink href="/catalog">Project Catalog</TextLink>&nbsp;/&nbsp;
        </Typography>
        <Typography variant="caption2">{name}</Typography>
      </Grid>
      <div style={{ padding: "8px 0px 8px 0px" }}>
        <Typography variant="h2">{capitalize(name)}</Typography>
      </div>
      {description.length && <Typography variant="body2">{description}</Typography>}
    </Grid>
  </>
);

export default ProjectHeader;
