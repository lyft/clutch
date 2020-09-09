module.exports = {
  collectCoverageFrom: [ "src/*.*sx"],
  coverageDirectory: "/tmp",
  coverageReporters: ["text", "cobertura"],
  coverageThreshold: {
    global: {
      statements: 0,
    },
  },
  moduleDirectories: ["node_modules", "src"],
  moduleNameMapper: {
    "\\.(css)$": "identity-obj-proxy",
  },
  setupFiles: ["jest-canvas-mock"],
  setupFilesAfterEnv: ["@clutch-sh/tools/jest.setup.js"],
  verbose: true,
};
