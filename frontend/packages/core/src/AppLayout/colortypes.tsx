import type { Color } from "@material-ui/core";

export interface StrokeColor {
  primary: string;
  secondary: string;
  tertiary: string;
  inverse: string;
}

export interface StateColor {
  active: string;
  default: string;
}

export interface BackgroundColor {
  primary: string;
  secondary: string;
  header: string;
}

export interface TypographyColor {
  primary: string;
  secondary: string;
  tertiary: string;
  inverse: string;
}

export interface ClutchPalette {
  neutral: Color;
  blue: Color;
  green: Color;
  amber: Color;
  red: Color;
}
// TODO: choose a better name for this
export interface ClutchColorChoices {
  stroke: StrokeColor;
  background: BackgroundColor;
  typography: TypographyColor; // Typography and Icons
  states: {
    active: string;
    default: string;
    hoverAlpha: number;
    focusedAlpha: number;
    pressedAlpha: number;
    selectedAlpha: number;
    disabledAlpha: number;
  };
}
