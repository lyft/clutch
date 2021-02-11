const baseConfig = require("@clutch-sh/tools/.eslintrc.js");

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
  ],
};
