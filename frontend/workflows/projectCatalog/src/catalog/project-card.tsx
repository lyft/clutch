import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import { Card, Grid, IconButton, Typography } from "@clutch-sh/core";
import styled from "@emotion/styled";
import CloseIcon from "@material-ui/icons/Close";
import KeyboardArrowRightIcon from "@material-ui/icons/KeyboardArrowRight";

const StyledCard = styled(Card)({
  display: "flex",
  flexDirection: "column",
  width: "350px",
  height: "216px",
  overflow: "hidden",
  padding: "22px 19px 27px 29px",
  ":hover": {
    cursor: "pointer",
    backgroundColor: "#F5F6FD",
    ".showOnHover": {
      visibility: "visible",
    },
  },
  ":active": {
    backgroundColor: "#D7DAF6",
  },
  ".showOnHover": {
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
      <Grid container wrap="nowrap">
        <Grid container item direction="row" alignItems="center">
          <Grid item xs={10}>
            <Typography variant="h6" color="secondary">
              {project?.name?.toUpperCase()}
            </Typography>
          </Grid>
          <Grid container item className="showOnHover" justify="flex-end" xs={2}>
            <IconButton size="small" variant="neutral" onClick={remove}>
              <CloseIcon color="secondary" />
            </IconButton>
          </Grid>
        </Grid>
      </Grid>
      <Grid
        container
        style={{ marginTop: "8px", paddingRight: "16px", flex: "1", overflow: "hidden" }}
        zeroMinWidth
      >
        <Typography variant="body2" color="rgba(13, 16, 48, 0.65)">
          {project?.data?.description}
        </Typography>
      </Grid>
      <Grid container item justify="flex-end" alignItems="center">
        <KeyboardArrowRightIcon color="disabled" />
      </Grid>
    </StyledCard>
  );
};

export default ProjectCard;
