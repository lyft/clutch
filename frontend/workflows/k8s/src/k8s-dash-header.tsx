import React from "react";
import styled from "@emotion/styled";

const Category = styled.div({
  fontWeight: 700,
  fontSize: "12px",
  lineHeight: "16px",
  color: "rgba(13, 16, 48, 0.38)",
  textTransform: "uppercase",
  paddingBottom: "8px",
});

const HeaderText = styled.div({
  fontWeight: 700,
  fontSize: "26px",
  lineHeight: "32px",
  color: "#0D1030",
});

const K8sDashHeader = () => (
  <div>
    <Category>Kubernetes/</Category>
    <HeaderText>Kubernetes Resource Dashboard</HeaderText>
  </div>
);

export default K8sDashHeader;
