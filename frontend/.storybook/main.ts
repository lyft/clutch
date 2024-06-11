import { StorybookConfig } from "@storybook/react-webpack5";

const config: StorybookConfig = {
  stories: ["../packages/**/*.stories.@(tsx|jsx)"],
  typescript: {
    reactDocgen: "react-docgen-typescript",
    reactDocgenTypescriptOptions: {
      compilerOptions: {
        allowSyntheticDefaultImports: false,
        esModuleInterop: false,
      },
    },
  },
  framework: {
    name: "@storybook/react-webpack5",
    options: { fastRefresh: true },
  },
  webpackFinal: async (config, { configType }) => {
    config?.module?.rules?.push({
      test: /\.(ts|tsx)$/,
      use: [
        {
          loader: require.resolve("esbuild-loader"),
          options: {
            target: "esnext",
          },
        },
      ],
    });

    config?.resolve?.extensions?.push(".ts", ".tsx");

    return config;
  },
  babel: async (options) => ({
    ...options,
    plugins: [
      "@babel/plugin-proposal-optional-chaining",
      "@babel/plugin-proposal-nullish-coalescing-operator",
      "@babel/plugin-transform-runtime",
      ["@babel/plugin-proposal-class-properties", { loose: true }],
    ],
  }),
  addons: [
    "@storybook/addon-essentials",
    "@storybook/addon-links",
    "@storybook/addon-a11y",
    "@storybook/addon-themes"
  ],
};

export default config;
