import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import { Card, Grid, IconButton, Typography } from "@clutch-sh/core";
import styled from "@emotion/styled";
import CloseIcon from "@material-ui/icons/Close";
import KeyboardArrowRightIcon from "@material-ui/icons/KeyboardArrowRight";

import LanguageIcon from "../helpers/language-icon";

const StyledCard = styled(Card)({
  display: "flex",
  flexDirection: "column",
  width: "384px",
  height: "214px",
  overflow: "hidden",
  padding: "13px 17px 13px 36px",
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
      <Grid container justify="flex-end">
        <Grid item className="showOnHover">
          <IconButton size="small" variant="neutral" onClick={remove}>
            <CloseIcon />
          </IconButton>
        </Grid>
      </Grid>
      <Grid container wrap="nowrap">
        <Grid container item direction="row" alignItems="flex-end">
          <Grid item xs={1}>
            <LanguageIcon language={(project?.languages || [])[0]} />
          </Grid>
          <Grid item xs={11}>
            <Typography variant="caption2">{project?.name}</Typography>
          </Grid>
        </Grid>
      </Grid>
      <Grid
        container
        style={{ marginTop: "16px", paddingRight: "16px", flex: "1", overflow: "hidden" }}
        zeroMinWidth
      >
        <Typography variant="body2">{project.data?.description}</Typography>
      </Grid>
      <Grid container justify="flex-end">
        <Grid item className="showOnHover">
          <KeyboardArrowRightIcon />
        </Grid>
      </Grid>
    </StyledCard>
  );
};

export default ProjectCard;
