const path = require('path')
const webpack = require('webpack')
module.exports = {
  pages: {
    index: {
      entry: 'frontend/main.js',
      title: "Role Based N Security",
    }
  },
  css: {
    extract: false
  },
  transpileDependencies: [
    'vuetify'
  ],
  devServer: {
    proxy: {
      '^/api/v1': {
        target: 'http://localhost:8080',
        secure: false
      },
      'settings.json': {
        target: 'http://localhost:8080',
        secure: false
      }
    },
    disableHostCheck: true,
    port: 8081
  },
  configureWebpack: {
    resolve: {
      extensions: ['.webpack.js', '.web.js', '.js', '.vue'],
      alias: {
        '@': path.resolve(__dirname, 'frontend'),
        '@assets': path.resolve(__dirname, 'frontend', "assets"),
        '@plugins': path.resolve(__dirname, 'frontend', 'plugins'),
        '@page': path.resolve(__dirname, 'frontend', 'components', 'pages'),
        '@tpl': path.resolve(__dirname, 'frontend', 'components', 'templates'),
        '@org': path.resolve(__dirname, 'frontend', 'components', 'organisms'),
        '@mixin': path.resolve(__dirname, 'frontend', 'mixins'),
      }
    },
    output: {
      filename: 'rbns.js',
      chunkFilename: 'rbns.js'
    },
    plugins: [
      new webpack.optimize.LimitChunkCountPlugin({
        maxChunks: 1
      })
    ]
  },
  // outputDir: 'static',
  // publicPath: './static'
}
