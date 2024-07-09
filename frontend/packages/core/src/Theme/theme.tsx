import React from "react";
import type { Theme as MuiTheme } from "@mui/material";
import {
  createTheme as createMuiTheme,
  CssBaseline,
  StyledEngineProvider,
  ThemeProvider as MuiThemeProvider,
} from "@mui/material";
import { defaultsDeep } from "lodash";

import { clutchColors } from "./colors";
import palette from "./palette";
import type {
  ChartColors,
  ClutchColors,
  CustomClutchTheme,
  ThemeOverrides,
  ThemeVariant,
} from "./types";

declare module "@emotion/react" {
  export interface Theme extends MuiTheme {}
}

declare module "@mui/material/styles" {
  interface Theme {
    colors: ClutchColors;
    chartColors: ChartColors;
  }
  interface ThemeOptions {
    colors?: ClutchColors;
    chartColors?: ChartColors;
  }
  interface Palette {
    contrastColor: string;
    headerGradient: string;
    brandColor: string;
  }
}

const defaultTheme = (variant: ThemeVariant): MuiTheme => {
  const { colors, chartColors } = clutchColors(variant);
  return createMuiTheme({
    colors,
    chartColors,
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

const generateTheme = (variant: ThemeVariant, theme?: CustomClutchTheme): MuiTheme => {
  const baseTheme = defaultTheme(variant);
  if (!theme) {
    return baseTheme;
  }
  const { colors, chartColors, palette: basePalette, ...rest } = baseTheme;

  return createMuiTheme({
    colors: defaultsDeep(theme.colors ?? {}, colors),
    chartColors: defaultsDeep(theme.chartColors ?? {}, chartColors),
    palette: defaultsDeep(theme?.palette ?? {}, basePalette),
    ...defaultsDeep(theme?.themeOptions ?? {}, rest),
  });
};

interface ThemeProps {
  variant?: ThemeVariant;
  children: React.ReactNode;
  overrides?: ThemeOverrides;
}

const ThemeProvider = ({ children, variant, overrides = {} }: ThemeProps) => (
  <StyledEngineProvider injectFirst>
    <MuiThemeProvider theme={generateTheme(variant, overrides[variant])}>
      <CssBaseline />
      {children}
    </MuiThemeProvider>
  </StyledEngineProvider>
);

export default ThemeProvider;
