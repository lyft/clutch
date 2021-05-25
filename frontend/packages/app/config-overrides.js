module.exports = function override(config) {
  const loaders = config.module.rules[1].oneOf;
  loaders.splice(2, 0, {
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
  });
  return config;
};
