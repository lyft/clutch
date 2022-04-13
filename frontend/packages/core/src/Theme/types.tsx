import type { Color } from "@material-ui/core";

export type ThemeVariant = "light" | "dark";

export interface StrokeColor {
  primary: string;
  secondary: string;
  tertiary: string;
  inverse: string;
}

export interface ClutchColors {
  neutral: Color;
  blue: Color;
  green: Color;
  amber: Color;
  red: Color;
}

export type ComponentState = "hover" | "focused" | "pressed" | "selected" | "disabled";
