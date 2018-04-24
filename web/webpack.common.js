const HtmlWebpackPlugin = require('html-webpack-plugin')
const path = require('path')

module.exports = {
  entry: [
    './src/index.js'
  ],
  devtool: 'none',
  module: {
    rules: [
      {
        test: /\.(js|jsx)$/,
        exclude: /node_modules/,
        use: ['babel-loader']
      },
      {
        test: /\.(svg|png|jpg|gif)$/,
        use: [
          {
            loader: 'file-loader',
            options: {
              outputPath: 'assets'
            }
          }
        ]
      }
    ]
  },
  resolve: {
    extensions: ['*', '.js', '.jsx']
  },
  output: {
    path: path.join(__dirname, '/dist'),
    publicPath: '/',
    filename: 'assets/bundle.js'
  },
  performance: {
    hints: false
  },
  plugins: [
    new HtmlWebpackPlugin({
      title: 'Pivotal Ignition',
      template: 'src/index.html',
      hash: true,
      filename: 'index.html'
    })
  ]
}
