module.exports = function override(config) {
  const loaders = config.module.rules.find(rule => rule.oneOf).oneOf;
  loaders.splice(
    2,
    0,
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
    },
    {
      test: /\.js$/,
      enforce: "pre",
      use: ["source-map-loader"],
    }
  );

  return { ...config, ignoreWarnings: [/Failed to parse source map/] };
};
