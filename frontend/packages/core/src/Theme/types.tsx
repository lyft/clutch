import type { Color, Theme as MuiTheme } from "@mui/material";
import type { PaletteOptions as MuiPaletteOptions } from "@mui/material/styles";

import type VARIANTS from "./variants";

export type ThemeVariant = keyof typeof VARIANTS;

export interface StrokeColor {
  primary: string;
  secondary: string;
  tertiary: string;
  inverse: string;
}

export interface ThemeOptions extends MuiTheme {}

export interface PaletteOptions extends MuiPaletteOptions {
  type: ThemeVariant;
  contrastColor: string;
  headerGradient: string;
  brandColor: string;
}

export interface ClutchColors {
  neutral: Color;
  blue: Color;
  green: Color;
  amber: Color;
  red: Color;
}

interface CommonChartColors {
  data: string[];
}

interface PieChartColors {
  labelPrimary: string;
  labelSecondary: string;
}

interface TimelineChartColors {
  xAxisStroke: string;
  tooltipBackgroundColor: string;
  tooltipTextColor: string;
  gridBackgroundColor: string;
  gridStroke: string;
}

export interface ChartColors {
  common: CommonChartColors;
  pie: PieChartColors;
  linearTimeline: TimelineChartColors;
}

export interface ClutchTheme {
  colors: ClutchColors;
  chartColors: ChartColors;
}

type Partial<T> = {
  [P in keyof T]?: T[P];
};

export type CustomClutchColor = Partial<Color>;
export type CustomClutchColors = {
  [P in keyof ClutchColors]?: CustomClutchColor;
};
export type CustomChartColors = {
  [P in keyof ChartColors]?: Partial<ChartColors[P]>;
};
export type CustomPalette = {
  [P in keyof PaletteOptions]?: Partial<PaletteOptions[P]>;
};
export type CustomThemeOptions = {
  [P in keyof ThemeOptions]?: Partial<ThemeOptions[P]>;
};
export interface CustomClutchTheme {
  colors?: CustomClutchColors;
  chartColors?: CustomChartColors;
  palette?: CustomPalette;
  themeOptions?: CustomThemeOptions;
}

export type ThemeOverrides = {
  [P in ThemeVariant]?: CustomClutchTheme;
};

export type ComponentState = "hover" | "focused" | "pressed" | "selected" | "disabled";
