import { createMuiTheme, ThemeOptions } from "@material-ui/core";
import { useTheme as useMuiTheme } from "@material-ui/core/styles";
import type { PaletteOptions } from "@material-ui/core/styles/createPalette";
import {} from "styled-components";

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

const getTheme = () => {
  return createMuiTheme({
    palette: lightPalette(),
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

export { getTheme, useTheme };
