import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import { Card, Grid, IconButton, styled, Typography, useTheme } from "@clutch-sh/core";
import CloseIcon from "@mui/icons-material/Close";
import KeyboardArrowRightIcon from "@mui/icons-material/KeyboardArrowRight";
import { alpha, Theme } from "@mui/material";

const StyledCard = styled(Card)(({ theme }: { theme: Theme }) => ({
  display: "flex",
  flexDirection: "column",
  width: "350px",
  height: "216px",
  overflow: "hidden",
  padding: "22px 19px 27px 29px",
  ":hover": {
    cursor: "pointer",
    backgroundColor: theme.palette.primary[100],
    ".showOnHover": {
      visibility: "visible",
    },
  },
  ":active": {
    backgroundColor: theme.palette.primary[300],
  },
  ".showOnHover": {
    visibility: "hidden",
  },
}));

interface ProjectCardProps {
  project: IClutch.core.project.v1.IProject;
  onRemove: () => void;
}

const ProjectCard = ({ project, onRemove }: ProjectCardProps) => {
  const theme = useTheme();
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
          <Grid container item className="showOnHover" justifyContent="flex-end" xs={2}>
            <IconButton size="small" variant="neutral" onClick={remove}>
              <CloseIcon />
            </IconButton>
          </Grid>
        </Grid>
      </Grid>
      <Grid container item flex={1} marginTop={1} paddingRight={2} overflow="hidden" zeroMinWidth>
        <Typography variant="body2" color={alpha(theme.palette.secondary[900], 0.65)}>
          {project?.data?.description}
        </Typography>
      </Grid>
      <Grid container item justifyContent="flex-end" alignItems="center">
        <KeyboardArrowRightIcon color="disabled" />
      </Grid>
    </StyledCard>
  );
};

export default ProjectCard;
