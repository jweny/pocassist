const path = require("path");
const CracoAntDesignPlugin = require("craco-antd");
const TerserPlugin = require('terser-webpack-plugin')

// 生产环境去掉map文件
if (process.env.NODE_ENV === 'production') {
  process.env.GENERATE_SOURCEMAP = 'false'
}

module.exports = {
  webpack: {
    alias: {
      "@": path.resolve("src")
    },
  },
  devServer: {
    compress: true,
    port: 3333,
    proxy: {
      // 配置跨域
      "/v1": {
        target: "http://127.0.0.1:1231/api",
        ws: false,
        changOrigin: true, // 允许跨域
        pathRewrite: {
          "^/v1": "" // 请求的时候使用这个api就可以
        }
      }
    }
  },
  plugins: [
    {
      plugin: CracoAntDesignPlugin,
      options: {
        postcssLoaderOptions: {
          plugins: [
            // require("postcss-pxtorem")({ rootValue: 192, propList: ["*"] })
          ]
        }
      }
    }
  ]
};
