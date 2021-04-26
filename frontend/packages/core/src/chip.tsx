import * as React from "react";
import styled from "@emotion/styled";
import type { ChipProps as MuiChipProps } from "@material-ui/core";
import { Chip as MuiChip } from "@material-ui/core";

const CHIP_VARIANTS = [
  "error",
  "warn",
  "attention",
  "neutral",
  "active",
  "pending",
  "success",
] as const;
export interface ChipProps extends Pick<MuiChipProps, "label"> {
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
    borderColor: "##D87313",
  },
  attention: {
    background: "#E2E2E6",
    label: "#0D1030",
    borderColor: "##0D103061",
  },
  neutral: {
    background: "#F8F8F9",
    label: "#0D1030",
    borderColor: "#0D10301A",
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

const StyledChip = styled(MuiChip)(
  {
    height: "32px",
    cursor: "inherit",
    borderStyle: "solid",
    borderWidth: "1px",
    ".MuiChip-label": {
      fontSize: "12px",
      lineHeight: "20px",
      padding: "6px 12px",
    },
  },
  props => ({
    background: CHIP_COLOR_MAP[props["data-variant"]].background,
    color: CHIP_COLOR_MAP[props["data-variant"]].label,
    borderColor: CHIP_COLOR_MAP[props["data-variant"]].borderColor,
  })
);

const Chip = ({ variant, ...props }: ChipProps) => <StyledChip data-variant={variant} {...props} />;

export { Chip, CHIP_VARIANTS };
