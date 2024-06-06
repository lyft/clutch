import { withClutchTheme } from "./withClutchTheme.decorator";

export const decorators = [
  withClutchTheme({
    themes: {
      light: "light",
      dark: "dark",
    },
    defaultTheme: "light",
  }),
];

export const parameters = {
  backgrounds: {
    default: "clutch",
    values: [
      {
        name: "clutch",
        value: "#f9fafe",
      },
      {
        name: "light",
        value: "#ffffff",
      },
    ],
  },
};

const preview = {
  decorators,
  parameters,
};

export default preview;
