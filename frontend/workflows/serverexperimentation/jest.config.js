const baseConfig = require("@clutch-sh/tools/jest.config.js");

module.exports = {
  ...baseConfig,
  coverageThreshold: {
    ...baseConfig.coverageThreshold,
  },
};
