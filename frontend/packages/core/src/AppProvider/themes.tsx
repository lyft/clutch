import React from "react";
import type { Theme as MuiTheme, ThemeOptions } from "@mui/material";
import { CssBaseline, StyledEngineProvider } from "@mui/material";
import { PaletteOptions, useTheme as useMuiTheme } from "@mui/material/styles";
import { StylesProvider } from "@mui/styles";

import { ThemeProvider } from "../Theme";
import type { ClutchColors } from "../Theme/types";

declare module "@mui/styles/defaultTheme" {
  interface DefaultTheme extends MuiTheme {}
}

interface ClutchPalette extends PaletteOptions {
  accent: {
    main: string;
  };
  destructive: {
    main: string;
  };
}

interface ClutchTheme extends ThemeOptions {
  palette: ClutchPalette;
  colors: ClutchColors;
}

const useTheme = () => {
  return useMuiTheme() as ClutchTheme;
};

interface ThemeProps {
  // disabling temporarily as we figure out theming
  // eslint-disable-next-line react/no-unused-prop-types
  variant?: "light";
}

const Theme: React.FC<ThemeProps> = ({ children }) => {
  return (
    <StyledEngineProvider injectFirst>
      <ThemeProvider variant="light">
        <CssBaseline />
        <StylesProvider>{children}</StylesProvider>
      </ThemeProvider>
    </StyledEngineProvider>
  );
};

export { Theme, useTheme };
