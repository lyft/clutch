const baseConfig = require("@clutch/tools/jest.config");

module.exports = {
  ...baseConfig,
  coverageThreshold: {
    ...baseConfig.coverageThreshold,
  },
};
