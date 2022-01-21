// COLORS FROM FIGMA
// NOTE THAT THOSE WITH GRADIENTS ONLY HAVE THEIR BASES

interface NeutralColors {
  5: string;
  10: string;
  12: string;
  15: string;
  20: string;
  25: string;
  30: string;
  40: string;
  50: string;
  60: string;
  70: string;
  80: string;
  90: string;
  100: string;
}

interface BlueColors {
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

interface GreenColors {
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

interface AmberColors {
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

interface RedColors {
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

interface StrokeColors {
  primary: string;
  secondary: string;
  tertiary: string;
  inverse: string;
}
interface StatesColors {
  primaryHover: string;
  primaryFocused: string;
  primaryPressed: string;
  primarySelected: string;
  neutralHover: string;
  neutralFocused: string;
  neutralPressed: string;
  neutralSelected: string;
}

interface BackgroundColors {
  primary: string;
  secondary: string;
}

interface MainColors {
  primary: string;
  secondary: string;
  tertiary: string;
  interactive: string;
  negative: string;
  positive: string;
  inverse: string;
}

interface ClutchPalette {
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

export const LIGHT_MODE: ClutchPalette = {
  NEUTRAL: {
    5: "#F8F8F9", // PRIMARY BUTTON DISABLED / STROKE TERTIARY
    10: "#E2E2E6",
    12: "#E2E2E6", // SECONDARY BUTTON PRESSED / FOCUSED
    15: "#DBDBE0",
    20: "#CFCFD6",
    25: "#C2C3CB",
    30: "#A3A4B0", // STROKE SECONDARY / TEXT TERTIARY
    40: "#868797",
    50: "#6E7083", // TEXT SECONDARY
    60: "#56586E",
    70: "#494C64",
    80: "#31344F",
    90: "#1E213E",
    100: "#0D1030",
  },
  BLUE: {
    10: "#F9F9FE",
    20: "#F5F6FD",
    30: "#EBEDFB",
    40: "#D7DAF6",
    50: "#C2C8F2",
    60: "#727FE1",
    70: "#3548D4",
    80: "#1629B9",
    90: "#0A1CA6",
    100: "#011082",
  },
  GREEN: {
    10: "#E5FCE8",
    20: "#CBF7CF",
    30: "#ACF2B2",
    40: "#6CD47A",
    50: "#4AB958",
    60: "#32A140",
    70: "#1C872A",
    80: "#106E1D",
    90: "#086515",
    100: "#02590E",
  },
  AMBER: {
    10: "#FFFBEB",
    20: "#FEF3C7",
    30: "#FDE68A",
    40: "#FCD34D",
    50: "#FBBF24",
    60: "#F59E0B",
    70: "#D97706",
    80: "#B45309",
    90: "#92400E",
    100: "#78350F",
  },
  RED: {
    10: "#FDE9E7",
    20: "#FBDAD6",
    30: "#F5BFBA",
    40: "#F4AAA3",
    50: "#F38D84",
    60: "#E95F52",
    70: "#DB3615",
    80: "#C3290B",
    90: "#A11E05",
    100: "#811500",
  },
  STROKE: {
    primary: "#0D1030",
    secondary: "#0D103061",
    tertiary: "#0D10301A",
    inverse: "#FFFFFF",
  },
  BACKGROUND: {
    primary: "#FFFFFF",
    secondary: "#F9FAFE",
  },
  STATES: {
    primaryHover: "#3548D41A",
    primaryFocused: "#3548D40D",
    primaryPressed: "#3548D426",
    primarySelected: "#3548D433",
    neutralHover: "#0D10301A",
    neutralFocused: "#0D10300D",
    neutralPressed: "#0D103021",
    neutralSelected: "#0D103026",
  },
  MAIN: {
    primary: "#0D1030",
    secondary: "#0D103099",
    tertiary: "#0D103061",
    interactive: "#3548D4",
    negative: "#DB3615",
    positive: "#1C872A",
    inverse: "#FFFFFF",
  },
};

export const DARK_MODE: ClutchPalette = {
  NEUTRAL: {
    5: "#272946", // PRIMARY BUTTON DISABLED / STROKE TERTIARY
    10: "#33344F",
    12: "#373953", // SECONDARY BUTTON PRESSED / FOCUSED
    15: "#3E4059",
    20: "#494B63",
    25: "#55566D",
    30: "#606176", // STROKE SECONDARY / TEXT TERTIARY
    40: "#77788A",
    50: "#8D8E9D", // TEXT SECONDARY
    60: "#A4A5B1",
    70: "#B0B0BB",
    80: "#BBBBC4",
    90: "#D2D2D8",
    100: "#E8E8EB",
  },
  BLUE: {
    10: "#050653", // BG
    20: "#060869", // HOVER
    30: "#13117C", // FOCUSED
    40: "#161DBC", // PRESSED
    50: "#2A4FF6",
    60: "#4281F6",
    70: "#5AABF6",
    80: "#8CC4F8",
    90: "#C2E1FE",
    100: "#DCECFB",
  },
  GREEN: {
    10: "#00580B",
    20: "#0D7620",
    30: "#1C872A",
    40: "#299935",
    50: "#32A73E",
    60: "#54B45B",
    70: "#73C178",
    80: "#9CD29E",
    90: "#C3E4C4",
    100: "#E6F4E7",
  },
  AMBER: {
    10: "#6B442A",
    20: "#D26803",
    30: "#D97706",
    40: "#DF850A",
    50: "#E3900E",
    60: "#E69F2A",
    70: "#EAB04E",
    80: "#EFC67F",
    90: "#F6DCB1",
    100: "#FBF1E0",
  },
  RED: {
    10: "#6C1404",
    20: "#DB3615",
    30: "#E93D19",
    40: "#F7441F",
    50: "#FF4A23",
    60: "#FF6843",
    70: "#FF8464",
    80: "#FFA790",
    90: "#FFCABC",
    100: "#FBE8E7",
  },
  STROKE: {
    primary: "#FFFFFFDE",
    secondary: "#FFFFFF73",
    tertiary: "#FFFFFF40",
    inverse: "#0D1030",
  },
  BACKGROUND: {
    primary: "#FFFFFF0F",
    secondary: "#0D1030",
  },
  STATES: {
    primaryHover: "#5AABF60D",
    primaryFocused: "#5AABF61A",
    primaryPressed: "#5AABF626",
    primarySelected: "#5AABF633",
    neutralHover: "#FFFFFF12",
    neutralFocused: "#FFFFFF1F",
    neutralPressed: "#FFFFFF2E",
    neutralSelected: "#FFFFFF40",
  },
  MAIN: {
    primary: "#FFFFFFE5",
    secondary: "#FFFFFFBF",
    tertiary: "#FFFFFF61",
    interactive: "#3548D4",
    negative: "#DB3615",
    positive: "#1C872A",
    inverse: "#0D1030",
  },
};

// EXTRA COLORS (TODO) ////////////////////////////////////////////////////////
