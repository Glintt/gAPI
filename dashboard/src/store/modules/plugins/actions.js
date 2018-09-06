const api = require("@/api/plugins");

export const updatePlugins = ({ commit }) => {
  var pluginsPromise = new Promise(resolve => {
    api.all(response => {
      commit("plugins", response.body);
      resolve();
    });
  });

  var activePromise = new Promise(resolve => {
    api.active(response => {
      commit("active", response.body);
      resolve();
    });
  });

  Promise.all([pluginsPromise, activePromise]);
};
