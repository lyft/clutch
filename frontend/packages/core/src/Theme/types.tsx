import type { Color } from "@mui/material";

export type ThemeVariant = "light" | "dark";

export enum THEME_VARIANTS {
  light = "light",
  dark = "dark",
}

export interface StrokeColor {
  primary: string;
  secondary: string;
  tertiary: string;
  inverse: string;
}

export interface ChartColors {
  common: {
    data: string[];
  };
  pie: {
    label: string;
  };
}

export interface ClutchColors {
  neutral: Color;
  blue: Color;
  green: Color;
  amber: Color;
  red: Color;
  charts: ChartColors;
}

export type ComponentState = "hover" | "focused" | "pressed" | "selected" | "disabled";
