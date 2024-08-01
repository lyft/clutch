module.exports = function override(config) {
  const { resolve } = config;
  resolve.fallback = { fs: false, path: false, crypto: false };
  const loaders = config.module.rules.find((rule) => rule.oneOf).oneOf;
  loaders.splice(
    2,
    0,
    {
      test: /.*@babel\/runtime\/helpers\/esm\/.*m?js$/,
      resolve: {
        fullySpecified: false,
      },
      use: [
        {
          loader: "esbuild-loader",
          options: {
            loader: "jsx",
            target: "esnext",
          },
        },
      ],
    },
    {
      test: /\.(js|mjs|jsx|ts|tsx)$/,
      exclude: /.*node_modules.*/,
      use: [
        {
          loader: "esbuild-loader",
          options: {
            loader: "jsx",
            target: "esnext",
          },
        },
      ],
    }
  );

  return config;
};
