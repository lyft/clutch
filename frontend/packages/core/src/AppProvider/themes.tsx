import React from "react";
import { useTheme as useMuiTheme } from "@mui/material";
import type { Theme as MuiTheme } from "@mui/material/styles";

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
  // Uncomment to use dark mode
  // Detect system color mode
  const prefersDarkMode =
    window.matchMedia && window.matchMedia("(prefers-color-scheme: dark)").matches;

  return (
    <ThemeProvider variant={prefersDarkMode ? THEME_VARIANTS.dark : THEME_VARIANTS.light}>
      {children}
    </ThemeProvider>
  );
};

export { Theme, useTheme };
