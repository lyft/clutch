// Need a function that can return colors even outside of the theme where folks can reference them through imports
// ideally this would be something like a function per color where we return based on
// the current mode.
import { fade } from '@material-ui/core/styles/colorManipulator';
import type { ComponentState } from "./types";
import ThemeProvider from "./theme";
import { STATE_OPACITY } from './colors';

// Return the appropriate color for a specified state and color.
const stateColor = (state: ComponentState, color: string) => {
  return fade(color, STATE_OPACITY[state])
};

export { stateColor, ThemeProvider };