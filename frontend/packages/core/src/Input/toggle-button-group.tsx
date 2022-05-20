import * as React from "react";
import styled from "@emotion/styled";
import {
  ToggleButtonGroup as MuiToggleButtonGroup,
  ToggleButtonGroupProps as MuiToggleButtonGroupProps,
} from "@material-ui/lab";

export interface ToggleButtonGroupProps
  extends Pick<
    MuiToggleButtonGroupProps,
    "value" | "children" | "size" | "orientation" | "onChange"
  > {
  /** If true, multiple children options can be selected simultaneously. */
  multiple?: boolean;
}

const StyledMuiToggleButtonGroup = styled(MuiToggleButtonGroup)({
  ".MuiToggleButton-root": {
    "&.Mui-selected": {
      color: "#3548D4",
      backgroundColor: "rgba(53, 72, 212, 0.12)",
    },
    textTransform: "none",
  },
  padding: "6px 16px",
});

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
