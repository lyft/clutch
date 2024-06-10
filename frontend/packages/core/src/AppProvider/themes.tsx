import React from "react";
import { useTheme as useMuiTheme } from "@mui/material";
import type { Theme as MuiTheme } from "@mui/material/styles";

import { useUserPreferences } from "../Contexts";
import { ThemeProvider } from "../Theme";
import { THEME_VARIANTS } from "../Theme/colors";
import type { ClutchColors } from "../Theme/types";

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
  // Detect system color mode
  const prefersDarkMode =
    window.matchMedia && window.matchMedia("(prefers-color-scheme: dark)").matches;
  const themeVariant = preferences.themeMode;

  return <ThemeProvider variant={themeVariant}>{children}</ThemeProvider>;
};

export { Theme, useTheme };
