import * as React from "react";
import styled from "@emotion/styled";
import type { ChipProps as MuiChipProps } from "@material-ui/core";
import { Chip as MuiChip } from "@material-ui/core";

type variants = "error" | "unknown" | "pending" | "running" | "succeeded";
export interface ChipProps extends Pick<MuiChipProps, "label"> {
  variant: variants;
}

const CHIP_COLOR_MAP = {
  error: "#FF8A80",
  unknown: "#FFCC80",
  pending: "#FFF59D",
  running: "#C2C8F2",
  succeeded: "#69F0AE",
};

const StyledChip = styled(MuiChip)({
  height: "20px",
  fontSize: "12px",
  lineHeight: "20px",
  cursor: "inherit",
  margin: "10px",
});

const ClutchChip = ({ variant, ...props }: ChipProps) => {
  return <StyledChip {...props} style={{ background: CHIP_COLOR_MAP[variant] }} />;
};

export default ClutchChip;
