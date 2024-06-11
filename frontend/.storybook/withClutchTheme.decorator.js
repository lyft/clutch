import { DecoratorHelpers } from "@storybook/addon-themes";
import { ThemeProvider } from "../packages/core/src/Theme";

const {
  initializeThemeState,
  pluckThemeFromContext,
  useThemeParameters,
} = DecoratorHelpers;

export const withClutchTheme = ({ themes, defaultTheme }) => {
  initializeThemeState(Object.keys(themes), defaultTheme);

  return (Story, context) => {
    const selectedTheme = pluckThemeFromContext(context);
    const { themeOverride } = useThemeParameters();

    const selected = themeOverride || selectedTheme || defaultTheme;

    return (
      <ThemeProvider variant={themes[selected]}>
        <Story {...context} />
      </ThemeProvider>
    );
  };
};
