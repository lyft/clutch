import * as React from "react";
import type { RadioProps as MuiRadioProps, Theme } from "@mui/material";
import { alpha, Radio as MuiRadio } from "@mui/material";

import styled from "../styled";

const StyledRadio = styled(MuiRadio)<{ checked: RadioProps["selected"] }>(
  {
    ".MuiIconButton-label": {
      height: "24px",
      width: "24px",
      borderRadius: "100px",
      boxSizing: "border-box",
    },
  },
  props => ({ theme }: { theme: Theme }) => ({
    "&:hover > .MuiIconButton-label > div": {
      border: props.checked
        ? `1px solid ${theme.palette.primary[700]}`
        : `1px solid ${theme.palette.primary[600]}`,
    },
  })
);

const Icon = styled("div")<{ $disabled?: MuiRadioProps["disabled"] }>(
  ({ theme }: { theme: Theme }) => ({
    height: "24px",
    width: "24px",
    border: `1px solid ${alpha(theme.palette.secondary[900], 0.38)}`,
    borderRadius: "100px",
    boxSizing: "border-box",
  }),
  props => ({ theme }: { theme: Theme }) => ({
    border: props.$disabled
      ? `1px solid ${theme.palette.secondary[200]}`
      : `1px solid ${alpha(theme.palette.secondary[900], 0.38)}`,
  })
);

const SelectedIcon = styled("div")(({ theme }: { theme: Theme }) => ({
  height: "24px",
  width: "24px",
  background: theme.palette.primary[600],
  border: `1px solid ${theme.palette.primary[700]}`,
  borderRadius: "100px",
  boxSizing: "border-box",
}));

const SelectedCenter = styled("div")(({ theme }: { theme: Theme }) => ({
  height: "12px",
  width: "12px",
  background: theme.palette.contrastColor,
  borderRadius: "100px",
  boxSizing: "border-box",
  margin: "5px 5px",
}));

export interface RadioProps
  extends Pick<MuiRadioProps, "disabled" | "name" | "onChange" | "required" | "value"> {
  selected?: boolean;
}

const Radio: React.FC<RadioProps> = ({
  name,
  onChange,
  required,
  value,
  selected = false,
  disabled = false,
}) => {
  return (
    <StyledRadio
      checked={selected}
      checkedIcon={
        <SelectedIcon>
          <SelectedCenter />
        </SelectedIcon>
      }
      color="primary"
      icon={<Icon $disabled={disabled} />}
      disabled={disabled}
      name={name}
      onChange={onChange}
      required={required}
      value={value}
    />
  );
};

export default Radio;
