module.exports = {
  stories: [
    "../packages/**/*.stories.@(tsx|jsx)",
  ],
  typescript: {
    reactDocgen: 'react-docgen-typescript',
    reactDocgenTypescriptOptions: {
      compilerOptions: {
        allowSyntheticDefaultImports: false,
        esModuleInterop: false,
      },
    }
  },
  addons: [
    "@storybook/addon-essentials",
    "@storybook/addon-links",
  ],
}