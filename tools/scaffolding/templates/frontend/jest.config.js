const baseConfig = require("@clutch-sh/tools/jest.config");

module.exports = {
  ...baseConfig,
  coverageThreshold: {
    ...baseConfig.coverageThreshold,
  },
};
