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
      color: "rgba(53, 72, 212, 1)",
      backgroundColor: "rgba(53, 72, 212, 0.12)",
    },
    textTransform: "none",
  },
  padding: "6px 16px",
});

// TODO(smonero): add some tests
/** Note when passing null to the `value` prop, the behavior is:
 * If `multiple` is true, allow a value of null.
 *
 * If `multiple` is false, default to the first child since at least
 * one child must be selected.
 */
const ToggleButtonGroup = ({
  multiple = false,
  value,
  children,
  size = "medium",
  orientation = "horizontal",
  onChange,
}: ToggleButtonGroupProps) => {
  // Need to investigate why this doesn't work (circ reference object, hm????)
  return (
    <StyledMuiToggleButtonGroup
      exclusive={!multiple}
      value={value || (multiple ? value : children?.[0]?.props.value)}
      size={size}
      orientation={orientation}
      onChange={onChange}
    >
      {children}
    </StyledMuiToggleButtonGroup>
  );
};

export default ToggleButtonGroup;
