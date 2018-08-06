module.exports = {
  presets: ["@vue/app"],
  plugins: [
    [
      "transform-inline-environment-variables",
      {
        include: [
          "NODE_ENV",
          "API_HOST",
          "API_PORT",
          "SOCKET_HOST",
          "SOCKET_PORT",
          "API_PROTOCOL"
        ]
      }
    ]
  ]
};
