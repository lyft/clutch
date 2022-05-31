import React from "react";
import { ThemeProvider as EmotionThemeProvider } from "@emotion/react";
import type { Theme as MuiTheme } from "@mui/material";
import {
  adaptV4Theme,
  createTheme as createMuiTheme,
  CssBaseline,
  DeprecatedThemeOptions,
  StyledEngineProvider,
  ThemeProvider as MuiThemeProvider,
} from "@mui/material";
import { PaletteOptions, useTheme as useMuiTheme } from "@mui/material/styles";
import { StylesProvider } from "@mui/styles";

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

interface ClutchTheme extends DeprecatedThemeOptions {
  palette: ClutchPalette;
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
  return createMuiTheme(
    adaptV4Theme({
      palette: lightPalette(),
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
        // TODO: Figure out if this needs to be adjusted
        // MuiTypography: {
        //   colorPrimary: {
        //     color: NAVY,
        //   },
        //   colorSecondary: {
        //     color: GRAY,
        //   },
        // },
      },
    })
  );
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
    <StyledEngineProvider injectFirst>
      <MuiThemeProvider theme={theme()}>
        <EmotionThemeProvider theme={theme()}>
          <CssBaseline />
          <StylesProvider injectFirst>{children}</StylesProvider>
        </EmotionThemeProvider>
      </MuiThemeProvider>
    </StyledEngineProvider>
  );
};

export { Theme, useTheme };
