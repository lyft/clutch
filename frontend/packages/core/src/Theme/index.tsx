import React from "react";
import type { Theme as MuiTheme } from "@mui/material";
import { useTheme as useMuiTheme } from "@mui/material";
import { alpha } from "@mui/material/styles";
import get from "lodash/get";
import isEmpty from "lodash/isEmpty";

import { useUserPreferences } from "../Contexts";

import { STATE_OPACITY } from "./colors";
import ThemeProvider from "./theme";
import type { ComponentState, ThemeOverrides } from "./types";

interface ThemeProps {
  children: React.ReactNode;
  overrides?: ThemeOverrides;
}

// Return the appropriate color for a specified state and color.
const stateColor = (state: ComponentState, color: string) => {
  return alpha(color, STATE_OPACITY[state]);
};

const useTheme = () => useMuiTheme() as MuiTheme;

const Theme = ({ children, overrides }: ThemeProps) => {
  const { preferences } = useUserPreferences();

  const prefersDarkMode =
    window.matchMedia && window.matchMedia("(prefers-color-scheme: dark)").matches;

  const themeVariant = get(preferences, "theme");

  const variant = isEmpty(themeVariant) ? (prefersDarkMode ? "DARK" : "LIGHT") : themeVariant;

  return (
    <ThemeProvider variant={variant} overrides={overrides}>
      {children}
    </ThemeProvider>
  );
};

export { Theme, useTheme, stateColor, ThemeProvider };
