import * as React from "react";
import styled from "@emotion/styled";
import type { RadioProps as MuiRadioProps } from "@material-ui/core";
import { Radio as MuiRadio } from "@material-ui/core";

const StyledRadio = styled(MuiRadio)(
  {
    ".MuiIconButton-label": {
      height: "24px",
      width: "24px",
      borderRadius: "100px",
      boxSizing: "border-box",
    },
  },
  props => ({
    "&:hover > .MuiIconButton-label > div": {
      border: props.checked ? "1px solid #283CD2" : "1px solid #2E45DC",
    },
  })
);

const Icon = styled.div(
  {
    height: "1.5rem",
    width: "1.5rem",
    border: "1px solid rgba(13, 16, 48, 0.38)",
    borderRadius: "6.25rem",
    boxSizing: "border-box",
  },
  props => ({
    border: props["data-disabled"] ? "1px solid #DFE2E4" : "1px solid rgba(13, 16, 48, 0.38)",
  })
);

const SelectedIcon = styled.div({
  height: "1.5rem",
  width: "1.5rem",
  background: "#2E45DC",
  border: "1px solid #283CD2",
  borderRadius: "6.25rem",
  boxSizing: "border-box",
});

const SelectedCenter = styled.div({
  height: "0.75rem",
  width: "0.75rem",
  background: "#FFFFFF",
  borderRadius: "6.25rem",
  boxSizing: "border-box",
  margin: "0.313rem 0.313rem",
});

export interface RadioProps extends Pick<MuiRadioProps, "disabled" | "value"> {
  selected?: boolean;
}

const Radio: React.FC<RadioProps> = ({ selected, disabled, value }) => {
  return (
    <StyledRadio
      checked={selected}
      icon={<Icon data-disabled={disabled} />}
      checkedIcon={
        <SelectedIcon>
          <SelectedCenter />
        </SelectedIcon>
      }
      color="primary"
      value={value}
      disabled={disabled}
    />
  );
};

export default Radio;
