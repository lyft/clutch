import React from "react";
import type { Theme as MuiTheme } from "@mui/material";
import {
  createTheme as createMuiTheme,
  CssBaseline,
  StyledEngineProvider,
  ThemeProvider as MuiThemeProvider,
} from "@mui/material";
import { StylesProvider } from "@mui/styles";

import { clutchColors } from "./colors";
import palette from "./palette";
import type { ThemeVariant } from "./types";

declare module "@mui/styles/defaultTheme" {
  interface DefaultTheme extends MuiTheme {}
}

declare module "@emotion/react" {
  export interface Theme extends MuiTheme {}
}

// Create a Material UI theme is propagated to all children.
const createTheme = (variant: ThemeVariant): MuiTheme => {
  return createMuiTheme({
    colors: clutchColors(variant),
    palette: palette(variant),
    transitions: {
      // https://material-ui.com/getting-started/faq/#how-can-i-disable-transitions-globally
      create: () => "none",
    },
    components: {
      MuiButtonBase: {
        defaultProps: {
          // https://material-ui.com/getting-started/faq/#how-can-i-disable-the-ripple-effect-globally
          disableRipple: true,
        },
      },
      MuiAccordion: {
        styleOverrides: {
          root: {
            "&$expanded": {
              // remove the additional margin rule when expanded so the original margin is used.
              margin: null,
            },
          },
        },
      },
      MuiCssBaseline: {
        styleOverrides: {
          body: {
            fontSize: "0.875rem",
          },
        },
      },
      MuiSelect: {
        styleOverrides: {
          select: {
            fontSize: "0.875rem",
            height: "20px",
          },
        },
      },
    },
  });
};

interface ThemeProps {
  variant?: "light" | "dark";
  children: React.ReactNode;
}

const ThemeProvider = ({ children, variant = "light" }: ThemeProps) => (
  <StyledEngineProvider injectFirst>
    <MuiThemeProvider theme={createTheme(variant)}>
      <CssBaseline />
      <StylesProvider injectFirst>{children}</StylesProvider>
    </MuiThemeProvider>
  </StyledEngineProvider>
);

// Note that ThemeProvider can't be used until the Theme component can be replaced.
export default ThemeProvider;
