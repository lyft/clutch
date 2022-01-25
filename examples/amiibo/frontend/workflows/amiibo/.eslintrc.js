const baseConfig = require("@clutch/tools/.eslintrc");


module.exports = {
  ...baseConfig,
  overrides: [
    ...baseConfig.overrides,
    {
      files: ["**/*.test.*"],
      rules: {
        "import/no-extraneous-dependencies": ["off"],
      }
    }
  ]
};