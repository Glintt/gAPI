export const config = {
  API: {
    PROTOCOL: process.env.API_PROTOCOL || "http",
    HOST: process.env.API_HOST || "localhost",
    PORT: process.env.API_PORT || "8084",
    SOCKET_HOST: process.env.SOCKET_HOST || "localhost",
    SOCKET_PORT: process.env.SOCKET_PORT || "5000",
    BASE_PATH: "",
    SERVICE_DISCOVERY_BASEPATH: "/service-discovery",
    ANALYTICS_BASEPATH: "/analytics",

    getHeader: function(name) {
      var req = new XMLHttpRequest();
      req.open("GET", document.location, false);
      req.send(null);

      return req.getResponseHeader(name);
    },

    getApiBaseUrl: function() {
      let apiBaseUrl = this.getHeader("Api-Base");

      if (apiBaseUrl == undefined || apiBaseUrl == null) {
        apiBaseUrl = `${this.PROTOCOL}://${this.HOST}:${this.PORT}/${
          this.BASE_PATH
        }`;
      }
      return apiBaseUrl;
    },

    getSocketBaseUrl: function() {
      var socketBaseUrl = this.getHeader("Socket-Base");

      if (socketBaseUrl == undefined || socketBaseUrl == null) {
        socketBaseUrl = `${this.PROTOCOL}://${this.SOCKET_HOST}:${
          this.SOCKET_PORT
        }/`;
      }
      return socketBaseUrl;
    }
  }
};
