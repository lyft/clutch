import React from "react";
import styled from "@emotion/styled";
import CheckCircleIcon from "@mui/icons-material/CheckCircle";
import ThumbUpIcon from "@mui/icons-material/ThumbUp";
import { alpha, Grid, Theme } from "@mui/material";

const IconContainer = styled(Grid)(({ theme }: { theme: Theme }) => ({
  paddingTop: "4px",
  display: "flex",
  flexDirection: "column",
  justifyContent: "center",
  color: theme.palette.primary[600],
  fontSize: "7rem",
}));

const Icon = styled(ThumbUpIcon)({
  fontSize: "0.5em",
  marginBottom: "10px",
});

const TitleContainer = styled(Grid)(({ theme }: { theme: Theme }) => ({
  color: theme.palette.success[600],
  display: "flex",
  alignItems: "center",
  fontSize: "20px",
  fontWeight: 500,
  textTransform: "capitalize",
}));

const CheckmarkIcon = styled(CheckCircleIcon)({
  marginRight: "8px",
});

const SubtitleContainer = styled.div(({ theme }: { theme: Theme }) => ({
  color: alpha(theme.palette.secondary[900], 0.6),
  fontSize: "12px",
}));

const Confirmation: React.FC<{ action: string }> = ({ action, children }) => (
  <Grid container direction="column" justifyContent="center" alignItems="center">
    <IconContainer item>
      <Icon />
    </IconContainer>
    <TitleContainer item>
      <CheckmarkIcon /> {action} Requested!
    </TitleContainer>
    <SubtitleContainer>{children}</SubtitleContainer>
  </Grid>
);

export default Confirmation;
