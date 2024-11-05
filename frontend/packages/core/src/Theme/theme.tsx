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

type SpacingUnit = "xs" | "sm" | "base" | "md" | "lg" | "xl";

// NOTE: `string & {}` allows `SpacingUnit` to be autocompleted
type SpacingArg = SpacingUnit | number | (string & {});

declare module "@emotion/react" {
  export interface Theme extends MuiTheme {}
}

declare module "@mui/material/styles" {
  interface Spacing {
    (value: SpacingArg): string;
    (topBottom: SpacingArg, rightLeft: SpacingArg): string;
    (top: SpacingArg, rightLeft: SpacingArg, bottom: SpacingArg): string;
    (top: SpacingArg, right: SpacingArg, bottom: SpacingArg, left: SpacingArg): string;
  }

  interface Theme {
    spacing: Spacing;
    clutch: {
      useWorkflowLayout: boolean;
      spacing: Record<SpacingUnit, string>;
      layout: {
        gutter: string;
      };
    };
  }
  interface ThemeOptions {
    clutch: {
      useWorkflowLayout: boolean;
      spacing: Record<SpacingUnit, string>;
      layout: {
        gutter: string;
      };
    };
  }
}

// `8` is the default scaling factor in MUI, so we define and use it
// to allow predictability using a `number` value
// https://v5.mui.com/material-ui/customization/spacing/
const DEFAULT_SCALING_FACTOR = 8;

const CLUTCH_CUSTOM_SPACING: Record<SpacingUnit, string> = {
  xs: "4px",
  sm: "8px",
  base: "16px",
  md: "24px",
  lg: "32px",
  xl: "40px",
};

// Create a Material UI theme is propagated to all children.
const createTheme = (variant: ThemeVariant, useWorkflowLayout: boolean): MuiTheme => {
  return createMuiTheme({
    colors: clutchColors(variant),
    palette: palette(variant),
    spacing: (...args) => {
      const spacingValues = args.map(value => {
        if (typeof value === "number") {
          return `${value * DEFAULT_SCALING_FACTOR}px`;
        }

        if (typeof value === "string" && CLUTCH_CUSTOM_SPACING[value] !== undefined) {
          return CLUTCH_CUSTOM_SPACING[value];
        }

        return value;
      });

      return spacingValues.join(" ");
    },
    clutch: {
      useWorkflowLayout,
      spacing: { ...CLUTCH_CUSTOM_SPACING },
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
