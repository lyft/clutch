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
    height: "24px",
    width: "24px",
    border: "1px solid rgba(13, 16, 48, 0.38)",
    borderRadius: "100px",
    boxSizing: "border-box",
  },
  props => ({
    border: props["data-disabled"] ? "1px solid #DFE2E4" : "1px solid rgba(13, 16, 48, 0.38)",
  })
);

const SelectedIcon = styled.div({
  height: "24px",
  width: "24px",
  background: "#2E45DC",
  border: "1px solid #283CD2",
  borderRadius: "100px",
  boxSizing: "border-box",
});

const SelectedCenter = styled.div({
  height: "12px",
  width: "12px",
  background: "#FFFFFF",
  borderRadius: "100px",
  boxSizing: "border-box",
  margin: "5px 5px",
});

export interface RadioProps
  extends Pick<MuiRadioProps, "disabled" | "name" | "onChange" | "required" | "value"> {
  selected?: boolean;
}

const Radio: React.FC<RadioProps> = ({ selected, disabled, name, onChange, required, value }) => {
  return (
    <StyledRadio
      checked={selected}
      checkedIcon={
        <SelectedIcon>
          <SelectedCenter />
        </SelectedIcon>
      }
      color="primary"
      icon={<Icon data-disabled={disabled} />}
      disabled={disabled}
      name={name}
      onChange={onChange}
      required={required}
      value={value}
    />
  );
};

export default Radio;
