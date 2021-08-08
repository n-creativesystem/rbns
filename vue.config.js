const path = require('path')
module.exports = {
  pages: {
    index: {
      entry: 'frontend/main.js',
      title: "Role Based N Security",
    }
  },
  transpileDependencies: [
    'vuetify'
  ],
  devServer: {
    proxy: {
      '^/api/v1': {
        target: 'http://api-rbac-dev:8080',
        secure: false
      }
    },
    disableHostCheck: true
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
  },
  outputDir: 'static',
  publicPath: './'
  // publicPath: './',
}
