import * as React from "react";
import type { ChipProps as MuiChipProps, Theme } from "@mui/material";
import { alpha, Chip as MuiChip } from "@mui/material";

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

export interface ChipProps
  extends Pick<MuiChipProps, "clickable" | "onClick" | "label" | "size" | "icon"> {
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
  onDelete?: (e: any) => void;
  deleteIcon?: any;
}

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
  props => ({ theme }: { theme: Theme }) => {
    const CHIP_COLOR_MAP = {
      error: {
        background: theme.palette.error[50],
        label: theme.palette.error[600],
        borderColor: theme.palette.error[600],
      },
      warn: {
        background: theme.palette.warning[50],
        label: theme.palette.warning[600],
        borderColor: theme.palette.warning[600],
      },
      attention: {
        background: theme.palette.secondary[200],
        label: theme.palette.secondary[900],
        borderColor: alpha(theme.palette.secondary[900], 0.6),
      },
      neutral: {
        background: theme.palette.secondary[50],
        label: theme.palette.secondary[900],
        borderColor: alpha(theme.palette.secondary[300], 0.6),
      },
      active: {
        background: theme.palette.primary[200],
        label: theme.palette.primary[600],
        borderColor: theme.palette.primary[600],
      },
      pending: {
        background: theme.palette.warning[50],
        label: theme.palette.warning[500],
        borderColor: theme.palette.warning[500],
      },
      success: {
        background: theme.palette.success[50],
        label: theme.palette.success[500],
        borderColor: theme.palette.success[500],
      },
    };
    return {
      height: props.size === "small" ? "24px" : "32px",
      background: props.$filled
        ? CHIP_COLOR_MAP[props.$variant].borderColor
        : CHIP_COLOR_MAP[props.$variant].background,
      color: props.$filled ? theme.palette.contrastColor : CHIP_COLOR_MAP[props.$variant].label,
      borderColor: CHIP_COLOR_MAP[props.$variant].borderColor,
    };
  }
);

const Chip = ({
  variant,
  filled = false,
  size = "medium",
  onDelete,
  deleteIcon,
  ...props
}: ChipProps) => (
  <StyledChip
    $variant={variant}
    $filled={filled}
    size={size}
    onDelete={onDelete}
    deleteIcon={deleteIcon}
    {...props}
  />
);

export { Chip, CHIP_VARIANTS };
