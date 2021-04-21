import * as React from "react";
import styled from "@emotion/styled";
import type { ChipProps as MuiChipProps } from "@material-ui/core";
import { Chip as MuiChip } from "@material-ui/core";

export interface ChipProps extends Pick<MuiChipProps, "label"> {
  variant: "WORST" | "BAD" | "NEUTRALBAD" | "NEUTRAL" | "NEUTRALGOOD" | "GOOD" | "BEST";
}

const CHIP_BACKGROUND_COLOR_MAP = {
  WORST: "#F9EAE7", // red
  BAD: "#FEF8E8", // orange
  NEUTRALBAD: "#E2E2E6", // greyish
  NEUTRAL: "#F8F8F9", // whiteish
  NEUTRALGOOD: "#EBEDFA", // purplish
  GOOD: "#FFFEE8", // yellow
  BEST: "#E9F6EC", // green
};

const CHIP_LABEL_COLOR_MAP = {
  WORST: "#C2302E", // red
  BAD: "#D87313", // orange
  NEUTRALBAD: "#0D1030", // greyish
  NEUTRAL: "#0D1030", // whiteish
  NEUTRALGOOD: "#3548D4", // purplish
  GOOD: "#B09027", // yellow
  BEST: "#40A05A", // green
};

const StyledChip = styled(MuiChip)(
  {
    height: "20px",
    fontSize: "12px",
    lineHeight: "20px",
    cursor: "inherit",
    margin: "10px",
    borderStyle: "solid",
    borderWidth: "1px",
  },
  props => ({
    background: CHIP_BACKGROUND_COLOR_MAP[props["data-variant"]],
    color: CHIP_LABEL_COLOR_MAP[props["data-variant"]],
  })
);

const Chip = ({ variant, ...props }: ChipProps) => {
  return <StyledChip {...props} data-variant={variant} />;
};

export default Chip;
