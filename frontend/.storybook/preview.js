import { DecoratorHelpers } from "@storybook/addon-themes";
import { withClutchTheme } from "./withClutchTheme.decorator";
import { Global, css } from "@emotion/react";
import { clutchColors } from "../packages/core/src/Theme/colors";

const { pluckThemeFromContext } = DecoratorHelpers;

const withGlobalStyles = (Story, context) => {
  const selectedTheme = pluckThemeFromContext(context);
  const GlobalStyles = css`
    body {
      background: ${clutchColors(selectedTheme).neutral["50"]};
    }
  `;

  return (
    <>
      <Global styles={GlobalStyles} />
      <Story {...context} />
    </>
  );
};

export const decorators = [
  withClutchTheme({
    themes: {
      light: "light",
      dark: "dark",
    },
    defaultTheme: "light",
  }),
  withGlobalStyles,
];

export const parameters = {
  backgrounds: {
    default: "light",
    values: [
      {
        name: "light",
        value: clutchColors("light").neutral["50"],
      },
      {
        name: "dark",
        value: clutchColors("dark").neutral["50"],
      },
    ],
  },
};

const preview = {
  decorators,
  parameters,
};

export default preview;
