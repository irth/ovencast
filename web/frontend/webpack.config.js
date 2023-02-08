const HtmlWebpackPlugin = require("html-webpack-plugin");
const CopyPlugin = require("copy-webpack-plugin");
const path = require("path");

const dev = (conf) => {
  console.log("USING DEV CONFIG");
  return {
    ...conf,
    devtool: "inline-source-map",
  };
};

const prod = {
  entry: "./src/index.js",

  module: {
    rules: [
      {
        test: /\.js$/i,
        exclude: /node_modules/,
        use: { loader: "babel-loader" },
      },
      {
        test: /\.css$/i,
        use: ["style-loader", "css-loader"],
      },
    ],
  },
  plugins: [
    new HtmlWebpackPlugin({
      template: "src/index.html",
    }),
    new CopyPlugin({
      patterns: [{ from: "src/licenses", to: "licenses" }],
    }),
  ],
};

module.exports = process.env.NODE_ENV == "production" ? prod : dev(prod);
