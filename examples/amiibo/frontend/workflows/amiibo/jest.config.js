const baseConfig = require("@clutch/tools/jest.config.js");

module.exports = {
  ...baseConfig,
  coverageThreshold: {
    ...baseConfig.coverageThreshold,
  },
};
