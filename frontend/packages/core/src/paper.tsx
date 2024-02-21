import * as React from "react";
import styled from "@emotion/styled";
import type { PaperProps as MuiPaperProps, Theme } from "@mui/material";
import { alpha, Paper as MuiPaper } from "@mui/material";

export interface PaperProps extends Pick<MuiPaperProps, "children" | "className"> {}

const StyledPaper = styled(MuiPaper)(({ theme }: { theme: Theme }) => ({
  boxShadow: `0px 4px 6px ${alpha(theme.palette.primary[600], 0.2)}`,
  border: `1px solid ${alpha(theme.palette.secondary[900], 0.1)}`,
  background: theme.palette.contrastColor,
  padding: "16px",
  minWidth: "inherit",
  minHeight: "inherit",
}));

const Paper = ({ children, className }: PaperProps) => (
  <StyledPaper className={className}>{children}</StyledPaper>
);

export default Paper;
