const baseConfig = require("@clutch-sh/tools/jest.config");

module.exports = {
  ...baseConfig,
  coverageThreshold: {
    ...baseConfig.coverageThreshold,
    ".": {
      statements: "47",
    },
  },
  moduleNameMapper: {
    ...baseConfig.moduleNameMapper,
    "react-markdown": "<rootDir>/node_modules/react-markdown/react-markdown.min.js",
  }
};
