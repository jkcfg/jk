{
  "parser": "@typescript-eslint/parser",
  "plugins": ["jest", "import", "@typescript-eslint"],
  "extends": ["airbnb-base", "plugin:@typescript-eslint/recommended", "plugin:import/typescript"],
  "globals": {
    "V8Worker2": true
  },
  "settings": {
     "import/resolver": {
        "node": {
          "paths": ["."]
        }
     }
  },
  "env": {
    "jest/globals": true
  },
  "rules": {
    "indent": "off",
    "@typescript-eslint/indent": ["error", 2, {"SwitchCase": 0, "CallExpression": {"arguments": "first"}, "FunctionDeclaration": {"parameters": "first"}}],
    "no-undef": "off",
    "@typescript-eslint/no-explicit-any": off,
    "@typescript-eslint/no-angle-bracket-type-assertion": "off",
    "@typescript-eslint/explicit-member-accessibility": "off",
    "no-continue": 0,
    "lines-between-class-members": "off",
    "prefer-const": ["error", {"destructuring": "all"}],
    "import/prefer-default-export": 0,
    "import/no-extraneous-dependencies": "off",
    "import/no-unresolved": [
      2,
      {
        "ignore": [
          "^@jkcfg/std$"
        ]
      }
    ],
    "no-restricted-syntax": ["error", "ForInStatement", "LabeledStatement", "WithStatement"],
    "no-use-before-define": ["error", { "functions": false }],
    "no-useless-constructor": "off"
  }
}
