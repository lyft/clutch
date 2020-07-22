module.exports = {
  presets: ["@babel/preset-env", "@babel/preset-react", "@babel/typescript"],
  plugins: [
    "@babel/plugin-proposal-optional-chaining",
    "@babel/plugin-proposal-nullish-coalescing-operator",
    "@babel/plugin-transform-runtime",
  ],
};
