import React from "react";
import { ThemeProvider as EmotionThemeProvider } from "@emotion/react";
import { createMuiTheme, CssBaseline, MuiThemeProvider, ThemeOptions } from "@material-ui/core";
import { useTheme as useMuiTheme } from "@material-ui/core/styles";
import type { PaletteOptions } from "@material-ui/core/styles/createPalette";
import { StylesProvider } from "@material-ui/styles";
import { ThemeProvider as StyledThemeProvider } from "styled-components";

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
}

declare module "styled-components" {
  export interface ClutchTheme extends ThemeOptions {
    palette: ClutchPalette;
  }
}

const WHITE = "#ffffff";
const GRAY = "#D7DADB";
const TEAL = "#02acbe";
const RED = "#EF474D";
const NAVY = "#2D3F50";

const lightPalette = (): ClutchPalette => {
  return {
    accent: {
      main: TEAL,
    },
    destructive: {
      main: RED,
    },
    primary: {
      main: WHITE,
    },
    secondary: {
      main: NAVY,
    },
    background: {
      default: WHITE,
      paper: WHITE,
    },
    text: {
      primary: NAVY,
      secondary: GRAY,
    },
  };
};

const lightTheme = () => {
  return createMuiTheme({
    palette: lightPalette(),
    transitions: {
      // https://material-ui.com/getting-started/faq/#how-can-i-disable-transitions-globally
      create: () => "none"
    },
    props: {
      MuiButtonBase: {
        // https://material-ui.com/getting-started/faq/#how-can-i-disable-the-ripple-effect-globally
        disableRipple: true,
      },
    },
    overrides: {
      MuiTypography: {
        colorPrimary: {
          color: NAVY,
        },
        colorSecondary: {
          color: GRAY,
        },
      },
      MuiLink: {
        root: {
          color: TEAL,
        },
      },
    },
  });
};

const useTheme = () => {
  return useMuiTheme() as ClutchTheme;
};

interface ThemeProps {
  variant?: "light";
}

const Theme: React.FC<ThemeProps> = ({ children }) => {
  const theme = lightTheme;
  return (
    <MuiThemeProvider theme={theme()}>
      <EmotionThemeProvider theme={theme()}>
        <StyledThemeProvider theme={theme()}>
          <CssBaseline />
          <StylesProvider injectFirst>{children}</StylesProvider>
        </StyledThemeProvider>
      </EmotionThemeProvider>
    </MuiThemeProvider>
  );
};

export { Theme, useTheme };
