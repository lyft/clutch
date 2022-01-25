export interface NeutralColors {
  5: string; // PRIMARY BUTTON DISABLED / STROKE TERTIARY
  10: string;
  12: string; // SECONDARY BUTTON PRESSED / FOCUSED
  15: string;
  20: string;
  25: string;
  30: string; // STROKE SECONDARY / TEXT TERTIARY
  40: string;
  50: string; // TEXT SECONDARY
  60: string;
  70: string;
  80: string;
  90: string;
  100: string;
}

export interface BlueColors {
  10: string; // BG
  20: string; // HOVER
  30: string; // FOCUSED
  40: string; // PRESSED
  50: string;
  60: string;
  70: string;
  80: string;
  90: string;
  100: string;
}

export interface GreenColors {
  10: string;
  20: string;
  30: string;
  40: string;
  50: string;
  60: string;
  70: string;
  80: string;
  90: string;
  100: string;
}

export interface AmberColors {
  10: string;
  20: string;
  30: string;
  40: string;
  50: string;
  60: string;
  70: string;
  80: string;
  90: string;
  100: string;
}

export interface RedColors {
  10: string;
  20: string;
  30: string;
  40: string;
  50: string;
  60: string;
  70: string;
  80: string;
  90: string;
  100: string;
}

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

export interface MainColors {
  primary: string;
  secondary: string;
  tertiary: string;
  interactive: string;
  negative: string;
  positive: string;
  inverse: string;
}

export interface ClutchPalette {
  NEUTRAL: NeutralColors;
  BLUE: BlueColors;
  GREEN: GreenColors;
  AMBER: AmberColors;
  RED: RedColors;
  STROKE: StrokeColors;
  BACKGROUND: BackgroundColors;
  STATES: StatesColors;
  MAIN: MainColors; // Typography and Icons
}
