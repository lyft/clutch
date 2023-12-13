import React from "react";
import styled from "@emotion/styled";
import { alpha } from "@mui/material";

const Category = styled("div")(({ theme }) => ({
  fontWeight: 700,
  fontSize: "12px",
  lineHeight: "16px",
  color: alpha(theme.palette.secondary[900], 0.38),
  textTransform: "uppercase",
  paddingBottom: "8px",
}));

const HeaderText = styled("div")(({ theme }) => ({
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
