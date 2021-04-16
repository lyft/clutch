import * as React from "react";
import styled from "@emotion/styled";
import type { ChipProps as MuiChipProps } from "@material-ui/core";
import { Chip as MuiChip } from "@material-ui/core";

export interface ChipProps extends Pick<MuiChipProps, "label"> {
  variant: "WORST" | "BAD" | "NEUTRAL" | "BLANK" | "GOOD" | "BEST";
}

const CHIP_COLOR_MAP = {
  WORST: "#FF8A80", // red
  BAD: "#FFCC80", // orange
  NEUTRAL: "#FFF59D", // yellow
  BLANK: "#EBEDFB", // whiteish
  GOOD: "#C2C8F2", // purplish
  BEST: "#69F0AE", // green
};

const StyledChip = styled(MuiChip)(
  {
    height: "20px",
    fontSize: "12px",
    lineHeight: "20px",
    cursor: "inherit",
    margin: "10px",
  },
  props => ({
    background: CHIP_COLOR_MAP[props["data-variant"]],
  })
);

const Chip = ({ variant, ...props }: ChipProps) => {
  return <StyledChip {...props} data-variant={variant} />;
};

export default Chip;
