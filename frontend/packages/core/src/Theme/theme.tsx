import React from "react";
import { ThemeProvider as EmotionThemeProvider } from "@emotion/react";
import { createMuiTheme, CssBaseline, MuiThemeProvider, ThemeOptions } from "@material-ui/core";
import { StylesProvider } from "@material-ui/styles";
import type { Theme as MuiTheme } from "@material-ui/core";
import type { ThemeVariant } from "./types";
import palette from "./palette";
import Typography from "./typography";
import { clutchColors } from "./colors";

// Create a Material UI theme is propegated to all children. 
const createTheme = (variant: ThemeVariant): MuiTheme => {
  return createMuiTheme({
    //inject in custom colors
    ...clutchColors(variant),
    palette: palette(variant),
    typography: Typography,
    transitions: {
      // https://material-ui.com/getting-started/faq/#how-can-i-disable-transitions-globally
      create: () => "none",
    },
    props: {
      MuiButtonBase: {
        // https://material-ui.com/getting-started/faq/#how-can-i-disable-the-ripple-effect-globally
        disableRipple: true,
      },
    },
    overrides: {
      MuiAccordion: {
        root: {
          "&$expanded": {
            // remove the additional margin rule when expanded so the original margin is used.
            margin: null,
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
  <MuiThemeProvider theme={createTheme(variant)}>
    <EmotionThemeProvider theme={createTheme(variant)}>
      <CssBaseline />
      <StylesProvider injectFirst>{children}</StylesProvider>
    </EmotionThemeProvider>
  </MuiThemeProvider>
);

export default ThemeProvider;
