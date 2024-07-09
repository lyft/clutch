import type { ComponentState, ThemeVariant } from "./types";
import VARIANTS from "./variants";

export const brandColor = "#02acbe";

export const STATE_OPACITY: { [key in ComponentState]: number } = {
  hover: 0.5,
  focused: 0.1,
  pressed: 0.15,
  selected: 0.12,
  disabled: 0.5,
};

const clutchColors = (variant: ThemeVariant) => VARIANTS[variant];

export { clutchColors };
