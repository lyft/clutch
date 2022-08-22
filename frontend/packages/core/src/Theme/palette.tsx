import type { PaletteOptions as MuiPaletteOptions } from "@mui/material/styles";
import { alpha, TypeText } from "@mui/material/styles";

import { DARK_COLORS, LIGHT_COLORS, STATE_OPACITY } from "./colors";
import type { ClutchColors, ThemeVariant } from "./types";

interface PaletteOptions extends MuiPaletteOptions {
  type: ThemeVariant;
}

const lightText: Partial<TypeText> = {
  primary: LIGHT_COLORS.neutral[900],
  secondary: alpha(LIGHT_COLORS.neutral[900], 0.65),
  // tertiary
  // inverse
};

const darkText: Partial<TypeText> = {
  primary: alpha(DARK_COLORS.neutral[900], 0.9),
  secondary: alpha(DARK_COLORS.neutral[900], 0.75),
  // tertiary
  // inverse
};

const palette = (variant: ThemeVariant): PaletteOptions => {
  const isLightMode = variant === "light";
  const color = (isLightMode ? LIGHT_COLORS : DARK_COLORS) as ClutchColors;
  const inverseColor = (isLightMode ? DARK_COLORS : LIGHT_COLORS) as ClutchColors;

  // TODO: add all clutch colors to "common colors"
  return {
    type: variant,
    primary: color.blue,
    secondary: color.neutral,
    error: color.red,
    warning: color.amber,
    info: color.blue,
    success: color.green,
    grey: color.neutral,
    background: {
      default: inverseColor.neutral[900],
      // secondary
    },
    text: isLightMode ? lightText : darkText,
    action: {
      active: color.blue[600],
      activatedOpacity: STATE_OPACITY.pressed,
      hover: color.blue[600],
      hoverOpacity: STATE_OPACITY.hover,
      selected: color.blue[600],
      selectedOpacity: STATE_OPACITY.selected,
      focus: color.blue[600],
      focusOpacity: STATE_OPACITY.focused,
      disabled: color.neutral[900],
      disabledOpacity: STATE_OPACITY.disabled,
    },
  };
};

export default palette;
