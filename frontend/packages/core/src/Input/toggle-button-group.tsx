import * as React from "react";
import styled from "@emotion/styled";
import {
  ToggleButton as MuiToggleButton,
  ToggleButtonGroup as MuiToggleButtonGroup,
} from "@material-ui/lab";

export interface ToggleButtonGroupProps {
  exclusive?: boolean;
  currentValue: string;
  onChange: (event: React.ChangeEvent<{}>, value: string) => void;
  toggleButtonValues: string[];
}

const StyledMuiToggleButtonGroup = styled(MuiToggleButtonGroup)({
  ".MuiToggleButton-root": {
    "&.Mui-selected": {
      color: "#3548D4",
      backgroundColor: "#3548D41F",
    },
  },
  padding: "6px 16px 16px 6px",
});

// TODO(smonero): add some tests
const ToggleButtonGroup = ({
  exclusive = true,
  currentValue,
  onChange,
  toggleButtonValues,
}: ToggleButtonGroupProps) => (
  <StyledMuiToggleButtonGroup exclusive={exclusive} value={currentValue} onChange={onChange}>
    {toggleButtonValues.map(toggValue => (
      <MuiToggleButton value={toggValue}>{toggValue}</MuiToggleButton>
    ))}
  </StyledMuiToggleButtonGroup>
);

export default ToggleButtonGroup;
