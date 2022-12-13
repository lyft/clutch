import * as React from "react";
import styled from "@emotion/styled";
import type { ToggleButtonGroupProps as MuiToggleButtonGroupProps } from "@mui/lab";
import { ToggleButtonGroup as MuiToggleButtonGroup } from "@mui/material";

export { ToggleButton } from "@mui/material";

export interface ToggleButtonGroupProps
  extends Pick<
    MuiToggleButtonGroupProps,
    "value" | "children" | "size" | "orientation" | "onChange"
  > {
  /** If true, multiple children options can be selected simultaneously. */
  multiple?: boolean;
}

const StyledMuiToggleButtonGroup = styled(MuiToggleButtonGroup)(({ size }) => ({
  border: "1px solid rgba(13, 16, 48, 0.45)",
  padding: "8px",
  gap: "8px",
  background: "#FFFFFF",
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
      background: "#0D10300D",
    },
    "&.MuiToggleButton-root:active:not(.Mui-selected)": {
      background: "#0D10302E",
    },
    "&.Mui-selected": {
      background: "#3548D4",
      color: "#FFFFFF",
    },
    "&.Mui-disabled": {
      background: "rgba(13, 16, 48, 0.03)",
      color: "rgba(13, 16, 48, 0.48)",
    },
    background: "#FFFFFF",
    color: "#0D1030",
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
