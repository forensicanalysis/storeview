// webpack.config.js

const VuetifyLoaderPlugin = require('vuetify-loader/lib/plugin');

module.exports = {
  plugins: [
    new VuetifyLoaderPlugin()
  ],
  rules: [
    {
      test: /\.s(c|a)ss$/,
      use: [
        'vue-style-loader',
        'css-loader',
        {
          loader: 'sass-loader',
          // Requires sass-loader@^7.0.0
          /* options: {
            implementation: require('sass'),
            fiber: require('fibers'),
            indentedSyntax: true // optional
          }, */
          // Requires sass-loader@^8.0.0
          options: {
            implementation: require('sass'),
            sassOptions: {
              fiber: require('fibers'),
              indentedSyntax: true // optional
            },
          },
        },
      ],
    },
  ],
};
