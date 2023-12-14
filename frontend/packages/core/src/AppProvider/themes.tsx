import React from "react";
import { CssBaseline, StyledEngineProvider, useTheme as useMuiTheme } from "@mui/material";
import type { Theme as MuiTheme } from "@mui/material/styles";
import { StylesProvider } from "@mui/styles";

import { ThemeProvider } from "../Theme";
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
  }
}

const useTheme = () => useMuiTheme() as MuiTheme;

const Theme: React.FC = ({ children }) => {
  // Uncomment to use dark mode
  /* // Detect system color mode
  const prefersDarkMode =
    window.matchMedia && window.matchMedia("(prefers-color-scheme: dark)").matches; */
  const prefersDarkMode = false;
  return (
    <StyledEngineProvider injectFirst>
      <ThemeProvider variant={prefersDarkMode ? "dark" : "light"}>
        <CssBaseline />
        <StylesProvider>{children}</StylesProvider>
      </ThemeProvider>
    </StyledEngineProvider>
  );
};

export { Theme, useTheme };
