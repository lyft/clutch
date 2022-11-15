module.exports = {
  env: {
    browser: true,
    es2021: true,
  },
  extends: [
    "eslint:recommended",
    "plugin:react/recommended",
    "plugin:prettier/recommended",
    "standard-with-typescript",
  ],
  overrides: [],
  parserOptions: {
    ecmaVersion: "latest",
    sourceType: "module",
    project: "./tsconfig.json",
    // Required for VSCode to properly find tsconfig within docs dir.
    // See https://github.com/typescript-eslint/typescript-eslint/issues/251 for more info
    tsconfigRootDir: __dirname,
  },
  plugins: ["react"],
  settings: {
    react: {
      version: "detect",
    },
  },
  rules: {
    // prettier enforces trailing commas
    "comma-dangle": ["error", "only-multiline"],
    // prettier enforces double quotes
    quotes: ["error", "double"],
    // prettier enforces trailing semicolons
    semi: ["error", "always"],
    // prettier enforces indentation
    indent: ["off"],
    // prettier enforces multiline ternary logic
    "multiline-ternary": ["off"],
    // prettier enforces double quotes
    "@typescript-eslint/quotes": ["error", "double"],
    // prettier enforces trailing semicolons
    "@typescript-eslint/semi": ["error", "always"],
    // prettier enforces indentation
    "@typescript-eslint/indent": ["off"],
    // prettier enforces trailing semicolons
    "@typescript-eslint/member-delimiter-style": [
      "error",
      {
        multiline: {
          delimiter: "semi",
          requireLast: true,
        },
        singleline: {
          delimiter: "semi",
          requireLast: false,
        },
      },
    ],
    // prettier enforces trailing commas
    "@typescript-eslint/comma-dangle": ["error", "only-multiline"],
    // see https://eslint.org/docs/latest/rules/space-before-function-paren#anonymous-always-named-never-asyncarrow-always for more info
    "@typescript-eslint/space-before-function-paren": [
      "error",
      {
        anonymous: "always",
        named: "never",
        asyncArrow: "always",
      },
    ],
    // Prevents casting with `as` keyword
    "@typescript-eslint/consistent-type-assertions": ["off"],
  },
};
