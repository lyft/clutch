import type { Color } from "@material-ui/core";

export interface StrokeColors {
  primary: string;
  secondary: string;
  tertiary: string;
  inverse: string;
}

export interface StatesColors {
  primaryHover: string;
  primaryFocused: string;
  primaryPressed: string;
  primarySelected: string;
  neutralHover: string;
  neutralFocused: string;
  neutralPressed: string;
  neutralSelected: string;
}

export interface BackgroundColors {
  primary: string;
  secondary: string;
}

export interface TypographyColors {
  primary: string;
  secondary: string;
  tertiary: string;
  inverse: string;
}

export interface ClutchPalette {
  Neutral: Color;
  Blue: Color;
  Green: Color;
  Amber: Color;
  Red: Color;
}
// TODO: choose a better name for this
export interface ClutchColorChoices {
  Stroke: StrokeColors;
  Background: BackgroundColors;
  States: StatesColors;
  Typography: TypographyColors; // Typography and Icons
}