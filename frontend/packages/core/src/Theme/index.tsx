// Need a function that can return colors even outside of the theme where folks can reference them through imports
// ideally this would be something like a function per color where we return based on
// the current mode.
import { fade } from "@material-ui/core/styles/colorManipulator";

import { STATE_OPACITY } from "./colors";
import ThemeProvider from "./theme";
import type { ComponentState } from "./types";

// Return the appropriate color for a specified state and color.
const stateColor = (state: ComponentState, color: string) => {
  return fade(color, STATE_OPACITY[state]);
};

// TODO: Use the ThemeProvider once the Theme component can be replaced
export { stateColor, ThemeProvider };
