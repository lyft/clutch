/** @format */

import { dirname, join } from "path";
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
    name: getAbsolutePath("@storybook/react-webpack5"),
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
    getAbsolutePath("@storybook/addon-essentials"),
    getAbsolutePath("@storybook/addon-links"),
    getAbsolutePath("@storybook/addon-a11y"),
    getAbsolutePath("@storybook/preset-create-react-app"),
  ],

  docs: {},
};

export default config;

function getAbsolutePath(value: string): any {
  return dirname(require.resolve(join(value, "package.json")));
}
