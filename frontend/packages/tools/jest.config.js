module.exports = {
  roots: ["./src"],
  collectCoverageFrom: ["src/**/*.*sx", "!src/**/*.stories.*sx"],
  coverageDirectory: "/tmp",
  coverageReporters: ["text", "cobertura"],
  coverageThreshold: {
    global: {
      branches: 0,
      functions: 0,
      lines: 0,
      statements: 0,
    },
  },
  moduleDirectories: ["node_modules", "src"],
  moduleNameMapper: {
    "\\.(css)$": "identity-obj-proxy",
  },
  setupFilesAfterEnv: ["@clutch-sh/tools/jest.setup.js"],
  verbose: true,
  testEnvironment: "jsdom",
};
