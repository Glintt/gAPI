const userPermissionsAPI = require("@/api/user_permissions");

export const get = ({ commit }, payload) => {
  userPermissionsAPI.get(payload, response => {
    commit("updatePermissions", response.body);
  });
};

export const update = ({ commit }, { newPermissions, user }) => {
  userPermissionsAPI.update(user, newPermissions, response => {
    if (response.status !== 201) {
      console.log("Update failed");
    } else {
      commit("updatePermissions", newPermissions);
      console.log("Update success");
    }
  });
};
