const path = require('path');
// const  nodeExternals = require('webpack-node-externals');


module.exports = {
    target: 'web',
  entry: './src/serverEntry.js',

  mode: 'development',  
  output: {
    filename: 'bundle.js',
    path: path.resolve(__dirname, 'build'),
    publicPath: '/ui/', // note the trailing slash
  },
  module: {
    rules: [
      {
        test: /\.css$/,
        use: ['style-loader', 'css-loader'],
      },
      {
        test: /\.(js|jsx)$/,
        exclude: /node_modules/,
        use: {
          loader: 'babel-loader',
          options: {
            presets: ['@babel/preset-env', '@babel/preset-react']
          }
        }
      },
      {
          test: /\.svg$/,
          use: [
              {
                  loader: 'file-loader',
                  options: {
                      name: '[name].[ext]',
                      outputPath: 'images/',
                  }
              }
          ]
      },
    ]
  }
};