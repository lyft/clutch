import React from "react";
import { ThemeProvider as EmotionThemeProvider } from "@emotion/react";
import type { Theme as MuiTheme } from "@mui/material";
import {
  adaptV4Theme,
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
  interface Theme extends MuiTheme {}
}

// Create a Material UI theme is propagated to all children.
const createTheme = (variant: ThemeVariant): MuiTheme =>
  // return createMuiTheme(
  //   adaptV4Theme({
  //     // inject in custom colors
  //     ...clutchColors(variant),
  //     palette: palette(variant),
  //     transitions: {
  //       // https://material-ui.com/getting-started/faq/#how-can-i-disable-transitions-globally
  //       create: () => "none",
  //     },
  //     props: {
  //       MuiButtonBase: {
  //         // https://material-ui.com/getting-started/faq/#how-can-i-disable-the-ripple-effect-globally
  //         disableRipple: true,
  //       },
  //     },
  //     overrides: {
  //       MuiAccordion: {
  //         root: {
  //           "&$expanded": {
  //             // remove the additional margin rule when expanded so the original margin is used.
  //             margin: null,
  //           },
  //         },
  //       },
  //     },
  //   })
  // );
  createMuiTheme({
    palette: palette(variant),
    transitions: {
      create: () => "none",
    },
    components: {
      MuiButtonBase: {
        // https://material-ui.com/getting-started/faq/#how-can-i-disable-the-ripple-effect-globally
        defaultProps: {
          disableRipple: true,
        },
      },
      MuiCssBaseline: {
        styleOverrides: {
          body: {
            // Default as MUI changed fontSize to 1rem
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
      // MuiTypography: {
      //   styleOverrides: {
      //     root: {
      //       colorPrimary: {
      //         color: NAVY,
      //       },
      //       colorSecondary: {
      //         color: GRAY,
      //       },
      //     },
      //   },
      // },
    },
  });

interface ThemeProps {
  variant?: "light" | "dark";
  children: React.ReactNode;
}

const ThemeProvider = ({ children, variant = "light" }: ThemeProps) => (
  <StyledEngineProvider injectFirst>
    <EmotionThemeProvider theme={createTheme(variant)}>
      <CssBaseline />
      <StylesProvider injectFirst>{children}</StylesProvider>
    </EmotionThemeProvider>
  </StyledEngineProvider>
);

// Note that ThemeProvider can't be used until the Theme component can be replaced.
export default ThemeProvider;
