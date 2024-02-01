import * as React from "react";
import styled from "@emotion/styled";
import type { ToggleButtonGroupProps as MuiToggleButtonGroupProps } from "@mui/lab";
import { alpha, ToggleButtonGroup as MuiToggleButtonGroup } from "@mui/material";

export { ToggleButton } from "@mui/material";

export interface ToggleButtonGroupProps
  extends Pick<
    MuiToggleButtonGroupProps,
    "value" | "children" | "size" | "orientation" | "onChange"
  > {
  /** If true, multiple children options can be selected simultaneously. */
  multiple?: boolean;
}

const StyledMuiToggleButtonGroup = styled(MuiToggleButtonGroup)(({ size, theme }) => ({
  border: `1px solid ${alpha(theme.palette.secondary[900], 0.45)}`,
  padding: "8px",
  gap: "8px",
  background: theme.palette.contrastColor,
  ".MuiToggleButton-root": {
    flexDirection: "column",
    justifyContent: "center",
    alignItems: "center",
    padding: size === "small" ? "7px 32px" : "14px 32px",
    fontSize: size === "small" ? "14px" : "16px",
    borderRadius: "4px !important",
    border: "none",
    width: "100%",
    "&.MuiToggleButton-root:hover:not(.Mui-selected)": {
      background: alpha(theme.palette.secondary[900], 0.05),
    },
    "&.MuiToggleButton-root:active:not(.Mui-selected)": {
      background: alpha(theme.palette.secondary[900], 0.18),
    },
    "&.Mui-selected, &.Mui-selected:hover": {
      background: theme.palette.primary[600],
      color: theme.palette.contrastColor,
    },
    "&.Mui-disabled": {
      background: alpha(theme.palette.secondary[900], 0.03),
      color: alpha(theme.palette.secondary[900], 0.48),
    },
    background: theme.palette.contrastColor,
    color: theme.palette.secondary[900],
    textTransform: "none",
  },
}));

// TODO(smonero): add some tests
// TODO(smonero): add another component that is a parent component
// that enforces a default selection
const ToggleButtonGroup = ({
  multiple = false,
  value,
  children,
  size = "medium",
  orientation = "horizontal",
  onChange,
}: ToggleButtonGroupProps) => {
  return (
    <StyledMuiToggleButtonGroup
      exclusive={!multiple}
      value={value}
      size={size}
      orientation={orientation}
      onChange={onChange}
    >
      {children}
    </StyledMuiToggleButtonGroup>
  );
};

export default ToggleButtonGroup;
