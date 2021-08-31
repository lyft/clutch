module.exports = {
  roots: ["./src"],
  clearMocks: true,
  collectCoverageFrom: ["src/*.*sx"],
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
