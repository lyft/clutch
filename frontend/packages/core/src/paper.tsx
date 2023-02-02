import * as React from "react";
import styled from "@emotion/styled";
import type { PaperProps as MuiPaperProps } from "@mui/material";
import { Paper as MuiPaper } from "@mui/material";

export interface PaperProps extends Pick<MuiPaperProps, "children"> {}

const StyledPaper = styled(MuiPaper)({
  boxShadow: "0px 4px 6px rgba(53, 72, 212, 0.2)",
  border: "1px solid rgba(13, 16, 48, 0.1)",
  background: "#FFFFFF",
  padding: "16px",
  minWidth: "inherit",
  minHeight: "inherit",
});

const Paper = ({ children, ...props }: PaperProps) => (
  <StyledPaper {...props}>{children}</StyledPaper>
);

export default Paper;
