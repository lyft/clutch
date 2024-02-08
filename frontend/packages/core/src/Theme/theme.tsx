import React from "react";
import type { Theme as MuiTheme } from "@mui/material";
import {
  createTheme as createMuiTheme,
  CssBaseline,
  StyledEngineProvider,
  ThemeProvider as MuiThemeProvider,
} from "@mui/material";

import { clutchChartColors, clutchColors, THEME_VARIANTS } from "./colors";
import palette from "./palette";
import type { ClutchColors, ThemeVariant } from "./types";

declare module "@emotion/react" {
  export interface Theme extends MuiTheme {}
}

declare module "@mui/material/styles" {
  interface Theme {
    colors: ClutchColors;
    chartColors: string[];
  }
  interface ThemeOptions {
    colors?: ClutchColors;
    chartColors: string[];
  }
  interface Palette {
    contrastColor: string;
    headerGradient: string;
    brandColor: string;
  }
}

// Create a Material UI theme is propagated to all children.
const createTheme = (variant: ThemeVariant): MuiTheme => {
  return createMuiTheme({
    colors: clutchColors(variant),
    chartColors: clutchChartColors,
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
      MuiLink: {
        styleOverrides: {
          underlineAlways: {
            "&:not(:hover)": {
              textDecoration: "none",
            },
          },
        },
      },
    },
  });
};

interface ThemeProps {
  variant?: ThemeVariant;
  children: React.ReactNode;
}

const ThemeProvider = ({ children, variant = THEME_VARIANTS.light }: ThemeProps) => (
  <StyledEngineProvider injectFirst>
    <MuiThemeProvider theme={createTheme(variant)}>
      <CssBaseline />
      {children}
    </MuiThemeProvider>
  </StyledEngineProvider>
);

export default ThemeProvider;
