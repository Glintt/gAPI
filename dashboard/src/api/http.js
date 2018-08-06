import Vue from "vue";

const OAuthAPI = require("@/api/auth");
const API_CONFIG = require("@/configs/urls").config.API;

export function PathToCall(path) {
  return `${API_CONFIG.PROTOCOL}://${API_CONFIG.HOST}:${API_CONFIG.PORT}${
    API_CONFIG.BASE_PATH
  }${path}`;
}

function AddAuthorizationToHeader(config) {
  var token = OAuthAPI.getToken();

  if (config["headers"] === null || config["headers"] === undefined)
    config["headers"] = {};

  if (token !== null && token !== undefined)
    config["headers"] = Object.assign(config["headers"], {
      Authorization: "Bearer " + token
    });

  return config;
}

export function GET(url, params) {
  params = AddAuthorizationToHeader(params);

  return Vue.http.get(url, params);
}

export function DELETE(url, params) {
  params = AddAuthorizationToHeader(params);

  return Vue.http.delete(url, params);
}

export function PUT(url, body, params) {
  params = AddAuthorizationToHeader(params);

  return Vue.http.put(url, body, params);
}

export function POST(url, body, params) {
  params = AddAuthorizationToHeader(params);

  return Vue.http.post(url, body, params);
}

export function handleError(response, cb) {
  if (response.status === 401) {
    window.location.href = "/#/login";
    OAuthAPI.clearAccessToken();
  }

  cb(response);
}
