module.exports = {
  stories: [
    "../packages/**/*.stories.@(tsx|jsx)",
  ],
  typescript: {
    reactDocgen: 'none',
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