const baseConfig = require("@clutch-sh/tools/.eslintrc");

module.exports = {
  ...baseConfig,
  overrides: [
    ...baseConfig.overrides,
    {
      files: ["**/*.test.*", "**/*.mjs"],
      rules: {
        "import/no-extraneous-dependencies": ["off"],
      },
    },
  ],
};
