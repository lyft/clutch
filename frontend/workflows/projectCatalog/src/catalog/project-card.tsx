import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import { Card, Grid, IconButton, Typography } from "@clutch-sh/core";
import styled from "@emotion/styled";
import CloseIcon from "@material-ui/icons/Close";

const StyledCard = styled(Card)({
  width: "384px",
  height: "214px",
  overflow: "hidden",
  padding: "13px 17px 64px 36px",
  ":hover": {
    backgroundColor: "#F5F6FD",
    ".remove": {
      visibility: "visible",
    },
  },
  ":active": {
    backgroundColor: "#D7DAF6",
  },
  ".remove": {
    visibility: "hidden",
  },
});

interface ProjectCardProps {
  project: IClutch.core.project.v1.IProject;
  onRemove: () => void;
}

const ProjectCard = ({ project, onRemove }: ProjectCardProps) => {
  const remove = event => {
    event.stopPropagation();
    onRemove();
  };

  return (
    <StyledCard>
      <Grid container justify="flex-end">
        <Grid item className="remove">
          <IconButton size="small" variant="neutral" onClick={remove}>
            <CloseIcon />
          </IconButton>
        </Grid>
      </Grid>
      <Grid container wrap="nowrap">
        <Grid container item direction="row" alignItems="flex-end">
          <Grid item xs={11}>
            <Typography variant="caption2">{project?.name}</Typography>
          </Grid>
        </Grid>
      </Grid>
      <Grid
        container
        style={{ marginTop: "16px", paddingRight: "16px", height: "100%", overflow: "hidden" }}
        zeroMinWidth
      >
        <Typography variant="body2">{project.data?.description}</Typography>
      </Grid>
    </StyledCard>
  );
};

export default ProjectCard;
