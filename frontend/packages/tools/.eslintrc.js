module.exports = {
  parser: "babel-eslint",
  ignorePatterns: [
    "build/",
    "dist/",
    "node_modules/"
  ],
  extends: [
    "airbnb",
    "prettier",
    "prettier/react",
    "plugin:cypress/recommended",
    "plugin:import/typescript",
    "plugin:jest/recommended",
    "plugin:jest/style",
    "plugin:prettier/recommended"
  ],
  plugins: [
    "simple-import-sort"
  ],
  env: {
    "browser": true,
    "es6": true,
  },
  settings: {
    "import/extensions": [
      ".js",
      ".jsx",
      ".ts",
      ".tsx",
    ],
    "import/resolver": {
      "node": {
        "extensions": [".js", ".jsx", ".ts", ".tsx"],
        "paths": [
          "src"
        ]
      }
    },
    "import/parsers": {
      "@typescript-eslint/parser": [".ts", ".tsx"],
    }
  },
  reportUnusedDisableDirectives: true,
  rules: {
    "simple-import-sort/imports": [
      "error",
      {
        // Groups taken from example, plus internal packages since plugin doesn't do import resolution.
        // https://github.com/lydell/eslint-plugin-simple-import-sort/blob/15ad031/examples/.eslintrc.js#L74
        "groups": [
          [
            "^react",
            "^@?\\w"
          ],
          [
            // Internal packages.
            "^(components|workflows)(/.*|$)"
          ],
          [
            "^\\u0000"
          ],
          [
            "^\\.\\.(?!/?$)",
            "^\\.\\./?$"
          ],
          [
            "^\\./(?=.*/)(?!/?$)",
            "^\\.(?!/?$)",
            "^\\./?$"
          ],
          [
            "^.+\\.s?css$"
          ]
        ]
      }
    ],
    "consistent-return": [
      "off"
    ],
    "jest/expect-expect": [
      "error",
      {
        "assertFunctionNames": ["expect", "cy"]
      }

    ],
    "import/extensions": [
      "error",
      "ignorePackages",
      {
        "js": "never",
        "jsx": "never",
        "ts": "never",
        "tsx": "never",
      }
    ],
    "import/no-extraneous-dependencies": [
      "error",
      {
        "devDependencies": ["**/*.config.js"],
      }
    ],
    "no-console": [
      "error"
    ],
    "no-empty": [
      "off"
    ],
    "no-nested-ternary": [
      "off"
    ],
    "react/jsx-filename-extension": [
      1,
      {
        "extensions": [".js", ".jsx", ".ts", ".tsx"]
      }
    ],
    "react-hooks/exhaustive-deps": [
      "off"
    ],
    "react/jsx-no-duplicate-props": [
      1,
      {
        // Required due to Material UI TextField having both inputProps and InputProps props
        "ignoreCase": false,
      }
    ],
    "react/prop-types": [
      "off"
    ],
    "react/jsx-props-no-spreading": [
      "off",
      {
        "exceptions": [
          "Wizard"
        ]
      }
    ],
    "jsx-a11y/label-has-associated-control": [
      "error",
      {
        "required": {
          "some": [
            "nesting",
            "id"
          ]
        }
      }
    ]
  },
  overrides: [
    {
      files: ["**/*.ts", "**/*.tsx"],
      parser: "@typescript-eslint/parser",
      plugins: [
        "@typescript-eslint"
      ],
      rules: {
        "no-undef": "off",
        "no-unused-vars": "off",
        "@typescript-eslint/no-unused-vars": ["error"]
      }
    },
  ]
}
