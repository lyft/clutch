module.exports = {
  roots: ["./src"],
  collectCoverageFrom: [ "src/*.jsx"],
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
  setupFilesAfterEnv: ["@clutch-sh/tools/jest.setup.js"],
  verbose: true,
};
