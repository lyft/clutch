import React from "react";
import { useTheme as useMuiTheme } from "@mui/material";
import type { Theme as MuiTheme } from "@mui/material/styles";
import get from "lodash/get";
import isEmpty from "lodash/isEmpty";

import { useUserPreferences } from "../Contexts";
import { ThemeProvider } from "../Theme";
import type { ClutchColors } from "../Theme/types";
import { THEME_VARIANTS } from "../Theme/types";

declare module "@mui/material/styles" {
  interface Theme {
    colors: ClutchColors;
  }
  interface ThemeOptions {
    colors?: ClutchColors;
  }
  interface Palette {
    contrastColor: string;
    headerGradient: string;
    brandColor: string;
  }
}

const useTheme = () => useMuiTheme() as MuiTheme;

const Theme: React.FC = ({ children }) => {
  const { preferences } = useUserPreferences();

  const prefersDarkMode =
    window.matchMedia && window.matchMedia("(prefers-color-scheme: dark)").matches;

  const themeVariant = get(preferences, "theme");

  const variant = isEmpty(themeVariant)
    ? prefersDarkMode
      ? THEME_VARIANTS.dark
      : THEME_VARIANTS.light
    : themeVariant;

  return <ThemeProvider variant={variant}>{children}</ThemeProvider>;
};

export { Theme, useTheme };
