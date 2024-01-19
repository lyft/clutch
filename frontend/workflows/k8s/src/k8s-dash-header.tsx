import React from "react";
import { styled } from "@clutch-sh/core";
import { alpha, Theme } from "@mui/material";

const Category = styled("div")(({ theme }: { theme: Theme }) => ({
  fontWeight: 700,
  fontSize: "12px",
  lineHeight: "16px",
  color: alpha(theme.palette.secondary[900], 0.38),
  textTransform: "uppercase",
  paddingBottom: "8px",
}));

const HeaderText = styled("div")(({ theme }: { theme: Theme }) => ({
  fontWeight: 700,
  fontSize: "26px",
  lineHeight: "32px",
  color: theme.palette.secondary[900],
}));

const K8sDashHeader = () => (
  <div>
    <Category>Kubernetes/</Category>
    <HeaderText>Kubernetes Resource Dashboard</HeaderText>
  </div>
);

export default K8sDashHeader;
