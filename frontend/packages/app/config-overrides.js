const path = require("path");

module.exports = {
  // The Webpack config to use when compiling your react app for development or production.
  webpack(config) {
    config.module.rules.unshift({
      test: /\.(js|mjs|jsx|ts|tsx)$/,
      include: [path.join(__dirname, "src")],
      use: [
        {
          loader: "esbuild-loader",
          options: {
            loader: "jsx",
            target: "esnext",
          },
        },
      ],
    });
    return config;
  },
  // The Jest config to use when running your jest tests - note that the normal rewires do not
  // work here.
  jest(config) {
    return config;
  },
  // The function to use to create a webpack dev server configuration when running the development
  // server with 'npm run start' or 'yarn start'.
  // Example: set the dev server to use a specific certificate in https.
  devServer(configFunction) {
    return configFunction;
  },
  // The paths config to use when compiling your react app for development or production.
  paths(paths) {
    return paths;
  },
};
