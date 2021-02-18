const baseConfig = require("@clutch-sh/tools/.eslintrc.js");
const path = require("path");

module.exports = {
  ...baseConfig,
  overrides: [
    ...baseConfig.overrides,
    {
      files: ["**/*.test.*"],
      rules: {
        "import/no-extraneous-dependencies": [
          "error",
          {
            devDependencies: ["*.config.js"],
            packageDir: [__dirname, path.join(__dirname, "../tools")],
          },
        ],
      },
    },
    {
      files: ["**/*.stories.*"],
      rules: {
        "import/no-extraneous-dependencies": ["off"],
      },
    },
  ],
};
