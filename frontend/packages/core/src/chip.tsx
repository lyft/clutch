import * as React from "react";
import type { ChipProps as MuiChipProps } from "@material-ui/core";
import { Chip as MuiChip } from "@material-ui/core";

import styled from "./styled";

const CHIP_VARIANTS = [
  "error",
  "warn",
  "attention",
  "neutral",
  "active",
  "pending",
  "success",
] as const;

export interface ChipProps extends Pick<MuiChipProps, "label" | "size" | "icon"> {
  /**
   * Variant of chip.
   *
   * Types of variants sorted from worst to best:
   *  * error:     failed
   *  * warn:      warning
   *  * attention: not started
   *  * neutral:   default state
   *  * active:    active/running
   *  * pending:   low risk
   *  * success:   completion/succeeded
   */
  variant: typeof CHIP_VARIANTS[number];
  filled?: boolean;
}

const CHIP_COLOR_MAP = {
  error: {
    background: "#F9EAE7",
    label: "#C2302E",
    borderColor: "#C2302E",
  },
  warn: {
    background: "#FEF8E8",
    label: "#D87313",
    borderColor: "#D87313",
  },
  attention: {
    background: "#E2E2E6",
    label: "#0D1030",
    borderColor: "##0D103061",
  },
  neutral: {
    background: "#F8F8F9",
    label: "#0D1030",
    borderColor: "rgba(13, 16, 48, 0.1)",
  },
  active: {
    background: "#EBEDFA",
    label: "#3548D4",
    borderColor: "#3548D4",
  },
  pending: {
    background: "#FFFEE8",
    label: "#B09027",
    borderColor: "#B09027",
  },
  success: {
    background: "#E9F6EC",
    label: "#40A05A",
    borderColor: "#40A05A",
  },
};

const StyledChip = styled(MuiChip)<{
  $filled: ChipProps["filled"];
  $variant: ChipProps["variant"];
  size: ChipProps["size"];
}>(
  {
    cursor: "inherit",
    borderStyle: "solid",
    borderWidth: "1px",
    ".MuiChip-label": {
      fontSize: "14px",
      fontWeight: 400,
      lineHeight: "20px",
      padding: "7px 12px",
    },
  },
  props => ({
    height: props.size === "small" ? "24px" : "32px",
    background: props.$filled
      ? CHIP_COLOR_MAP[props.$variant].borderColor
      : CHIP_COLOR_MAP[props.$variant].background,
    color: props.$filled ? "#FFFFFF" : CHIP_COLOR_MAP[props.$variant].label,
    borderColor: CHIP_COLOR_MAP[props.$variant].borderColor,
  })
);

const Chip = ({ variant, filled = false, size = "medium", ...props }: ChipProps) => (
  <StyledChip $variant={variant} $filled={filled} size={size} {...props} />
);

export { Chip, CHIP_VARIANTS };
