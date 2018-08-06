const serviceDiscoveryAPI = require("@/api/service-discovery");

export const fetchGroups = ({ commit }) => {
  serviceDiscoveryAPI.listServiceGroups(response => {
    commit("updateGroups", response.body);
  });
};

export const updateGroup = ({ commit }, group) => {
  serviceDiscoveryAPI.updateServiceGroup(group, response => {
    commit("updateGroup", response.body);
  });
};

export const deleteGroup = ({ commit }, group) => {
  serviceDiscoveryAPI.deleteServiceGroup(group.Id, response => {
    if (response.status === 200) commit("groupDeleted", group);
  });
};
