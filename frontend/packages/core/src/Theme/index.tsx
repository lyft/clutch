import { alpha } from "@mui/material/styles";

import { STATE_OPACITY } from "./colors";
import ThemeProvider from "./theme";
import type { ComponentState } from "./types";

// Return the appropriate color for a specified state and color.
const stateColor = (state: ComponentState, color: string) => {
  return alpha(color, STATE_OPACITY[state]);
};

// TODO: Use the ThemeProvider once the Theme component can be replaced
export { stateColor, ThemeProvider };
