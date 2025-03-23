module.exports = {
  parser: '@typescript-eslint/parser',
  parserOptions: {
    project: 'tsconfig.json',
    sourceType: 'module',
  },
  plugins: ['@typescript-eslint/eslint-plugin'],
  extends: [
    'plugin:@typescript-eslint/recommended',
    'prettier'
  ],
  rules: {
    quotes: ['error', 'double'],
    '@typescript-eslint/quotes': ['error', 'double']
  }
};