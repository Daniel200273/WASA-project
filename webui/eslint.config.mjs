/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * @author ENDERZOMBI102 <enderzombi102.end@gmail.com> 2024
 * @description Quick and dirty `lint-openapi` config to better conform to the Prof's requests and style.
 */
import vue from "eslint-plugin-vue";

// noinspection JSUnusedGlobalSymbols
export default [
  ...vue.configs["flat/recommended"],
  {
    files: ["src/**/*.{js,vue}"],
    rules: {
      "vue/multi-word-component-names": "off",
      "vue/max-attributes-per-line": "off",
      "vue/require-default-prop": "off",
      "vue/singleline-html-element-content-newline": "off",
      "vue/html-indent": ["error", 2], // Use 2 spaces for indentation
      "vue/require-prop-types": "off", // Turn off prop type requirement for now
    },
  },
  {
    ignores: [
      "dist/**",
      "public/**",
      ".yarn/**",
      ".pnp.*",
      "node_modules/**",
      "*.cjs",
      "*.mjs",
    ],
  },
];
