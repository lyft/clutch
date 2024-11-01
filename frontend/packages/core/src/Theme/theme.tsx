import React from "react";
import type { Theme as MuiTheme } from "@mui/material";
import {
  createTheme as createMuiTheme,
  CssBaseline,
  StyledEngineProvider,
  ThemeProvider as MuiThemeProvider,
} from "@mui/material";

import { clutchColors, THEME_VARIANTS } from "./colors";
import palette from "./palette";
import type { ThemeVariant } from "./types";

declare module "@emotion/react" {
  export interface Theme extends MuiTheme {}
}

declare module "@mui/material/styles" {
  interface Theme {
    clutch: {
      useWorkflowLayout: boolean;
      spacing: {
        none: number;
        xs: number;
        sm: number;
        base: number;
        md: number;
        lg: number;
        xl: number;
      };
      layout: {
        gutter: string;
      };
    };
  }
  interface ThemeOptions {
    clutch: {
      useWorkflowLayout: boolean;
      spacing: {
        none: number;
        xs: number;
        sm: number;
        base: number;
        md: number;
        lg: number;
        xl: number;
      };
      layout: {
        gutter: string;
      };
    };
  }
}

// Create a Material UI theme is propagated to all children.
const createTheme = (variant: ThemeVariant, useWorkflowLayout: boolean): MuiTheme => {
  return createMuiTheme({
    colors: clutchColors(variant),
    palette: palette(variant),
    // `8` is the default scaling factor in MUI, we are setting it again to make it explicit
    // https://v5.mui.com/material-ui/customization/spacing/
    spacing: 8,
    clutch: {
      useWorkflowLayout,
      spacing: {
        none: 0,
        xs: 0.5,
        sm: 1,
        base: 2,
        md: 3,
        lg: 4,
        xl: 5,
      },
      layout: {
        gutter: useWorkflowLayout ? "0px" : "24px",
      },
    },
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
  useWorkflowLayout?: boolean;
  children: React.ReactNode;
}

const ThemeProvider = ({
  children,
  useWorkflowLayout = false,
  variant = THEME_VARIANTS.light,
}: ThemeProps) => (
  <StyledEngineProvider injectFirst>
    <MuiThemeProvider theme={createTheme(variant, useWorkflowLayout)}>
      <CssBaseline />
      {children}
    </MuiThemeProvider>
  </StyledEngineProvider>
);

export default ThemeProvider;
