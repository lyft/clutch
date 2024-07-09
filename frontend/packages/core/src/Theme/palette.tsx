import { alpha, TypeText } from "@mui/material/styles";

import { brandColor } from "./colors";
import type { ClutchColors, PaletteOptions, ThemeVariant } from "./types";
import VARIANTS from "./variants";

const lightText: Partial<TypeText> = {
  primary: VARIANTS.LIGHT.colors.neutral[900],
  secondary: alpha(VARIANTS.LIGHT.colors.neutral[900], 0.65),
  // tertiary
  // inverse
};

const darkText: Partial<TypeText> = {
  primary: alpha(VARIANTS.DARK.colors.neutral[900], 0.9),
  secondary: alpha(VARIANTS.DARK.colors.neutral[900], 0.75),
  // tertiary
  // inverse
};

const palette = (variant: ThemeVariant): PaletteOptions => {
  const isLightMode = variant === "LIGHT";
  const color = (isLightMode ? VARIANTS.LIGHT.colors : VARIANTS.DARK.colors) as ClutchColors;

  return {
    type: variant,
    mode: (variant === "LIGHT" ? "light" : "dark") as PaletteOptions["mode"],
    brandColor,
    primary: color.blue,
    secondary: color.neutral,
    error: color.red,
    warning: color.amber,
    info: color.blue,
    success: color.green,
    grey: color.neutral,
    background: {
      default: color.blue[50],
      paper: isLightMode ? "#fff" : "#1c1e3c",
    },
    text: isLightMode ? lightText : darkText,
    contrastColor: isLightMode ? "#ffffff" : "#000000", // Either black or white depending on theme
    headerGradient: isLightMode
      ? "linear-gradient(90deg, #38106b 4.58%, #131c5f 89.31%)"
      : "#0D1030",
  };
};

export default palette;
