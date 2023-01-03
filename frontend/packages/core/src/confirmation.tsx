import React from "react";
import CheckCircleIcon from "@mui/icons-material/CheckCircle";
import ThumbUpIcon from "@mui/icons-material/ThumbUp";
import { Grid } from "@mui/material";

import styled from "./styled";

const IconContainer = styled(Grid)({
  paddingTop: "4px",
  display: "flex",
  flexDirection: "column",
  justifyContent: "center",
  color: "#2F67F6",
  fontSize: "7rem",
});

const Icon = styled(ThumbUpIcon)({
  fontSize: "0.5em",
  marginBottom: "10px",
});

const TitleContainer = styled(Grid)({
  color: "#1E942E",
  display: "flex",
  alignItems: "center",
  fontSize: "20px",
  fontWeight: 500,
  textTransform: "capitalize",
});

const CheckmarkIcon = styled(CheckCircleIcon)({
  marginRight: "8px",
});

const SubtitleContainer = styled("div")({
  color: "rgba(13, 16, 48, 0.6)",
  fontSize: "12px",
});

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
