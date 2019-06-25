const userPermissionsAPI = require("@/api/user_permissions");

export const get = ({ commit }, payload) => {
  userPermissionsAPI.get(payload, response => {
    commit("updatePermissions", response.body);
  });
};

export const add = ({ dispatch }, { applicationId, user }) => {
  userPermissionsAPI.add(user, applicationId, response => {
    if (response.status !== 201) {
      console.log("Update failed");
    } else {
      dispatch("get", user);
      console.log("Update success");
    }
  });
};

export const remove = ({ dispatch }, { applicationId, user }) => {
  userPermissionsAPI.remove(user, applicationId, response => {
    if (response.status !== 201) {
      console.log("Delete failed");
    } else {
      dispatch("get", user);
      console.log("Delete success");
    }
  });
};
