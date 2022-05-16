import * as React from "react";
import styled from "@emotion/styled";
import {
  ToggleButton as MuiToggleButton,
  ToggleButtonGroup as MuiToggleButtonGroup,
} from "@material-ui/lab";

export interface ToggleButtonGroupProps {
  /** A boolean to decide if the user can select multiple options
   * simultaneously. If true, users will be able to select multiple options,
   * whereas if false, an exclusive choice will be enforced.
   */
  multiple?: boolean;
  /** The current value that is selected in the group. Usually this will be
   * a useState value.
   */
  currentValue: string;
  /** The onChange function that gets called when a user changes selection */
  onChange: (event: React.ChangeEvent<{}>, value: string) => void;
  /** An array of strings where each string will be one toggle button */
  toggleButtonValues: string[];
}

const StyledMuiToggleButtonGroup = styled(MuiToggleButtonGroup)({
  ".MuiToggleButton-root": {
    "&.Mui-selected": {
      color: "#3548D4",
      backgroundColor: "#3548D41F",
      inset: "1px 1px 1px 1px",
    },
    textTransform: "none",
  },
  padding: "6px 16px 16px 6px",
});

// TODO(smonero): add some tests
const ToggleButtonGroup = ({
  multiple = false,
  currentValue,
  onChange,
  toggleButtonValues,
}: ToggleButtonGroupProps) => (
  <StyledMuiToggleButtonGroup exclusive={!multiple} value={currentValue} onChange={onChange}>
    {toggleButtonValues.map(toggValue => (
      <MuiToggleButton value={toggValue}>{toggValue}</MuiToggleButton>
    ))}
  </StyledMuiToggleButtonGroup>
);

export default ToggleButtonGroup;
