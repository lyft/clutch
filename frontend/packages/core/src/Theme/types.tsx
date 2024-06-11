import type { Color } from "@mui/material";

export type ThemeVariant = "light" | "dark" | "system";

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
    labelPrimary: string;
    labelSecondary: string;
  };
  linearTimeline: {
    xAxisStroke: string;
    tooltipBackgroundColor: string;
    tooltipTextColor: string;
    gridBackgroundColor: string;
    gridStroke: string;
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
