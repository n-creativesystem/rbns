const path = require('path')
// const webpack = require('webpack')
const index = process.env.NODE_ENV == 'production' ? 'index.html' : 'index-debug.html',
  publicPath = process.env.STATIC_FILE || undefined
module.exports = {
  productionSourceMap: process.env.NODE_ENV == 'production' ? false : true,
  publicPath: publicPath,
  pages: {
    index: {
      entry: 'frontend/main.js',
      title: "Role Based N Security",
      template: path.join('public', index)
    }
  },
  css: {
    extract: false,
    requireModuleExtension: true,
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
      },
      '^/': {
        target: 'http://localhost:8080',
        secure: false
      }
    },
    disableHostCheck: true,
    port: 9999
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
    // output: {
    //   filename: 'rbns.js',
    //   chunkFilename: 'rbns.js'
    // },
    // plugins: [
    //   new webpack.optimize.LimitChunkCountPlugin({
    //     maxChunks: 1
    //   })
    // ]
  },
  // chainWebpack: config => {
  //   config.plugins.delete('preload-index')
  //   config.plugins.delete('prefetch-index')
  // }
}
