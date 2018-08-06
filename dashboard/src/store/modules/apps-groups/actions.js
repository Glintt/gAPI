const serviceDiscoveryAPI = require("@/api/service-discovery");

export const fetchGroups = ({ commit }) => {
  serviceDiscoveryAPI.listAppsGroups(response => {
    commit("updateGroups", response.body);
  });
};

export const updateGroup = ({ commit }, group) => {
  serviceDiscoveryAPI.updateAppsGroup(group, response => {
    commit("updateGroup", response.body);
  });
};

export const deleteGroup = ({ commit }, group) => {
  serviceDiscoveryAPI.deleteAppsGroup(group.Id, response => {
    if (response.status === 200) commit("groupDeleted", group);
  });
};

export const associateServiceToAppGroup = ({ commit }, payload) => {
  // TODO
  serviceDiscoveryAPI.addServiceToAppsGroup(
    payload.GroupId,
    payload.ServiceId,
    response => {
      // if (response.status === 200) commit('groupDeleted', group)
    }
  );
};

export const deassociateServiceFromAppGroup = ({ commit }, payload) => {
  // TODO
  serviceDiscoveryAPI.deassociateServiceFromAppsGroup(
    payload.GroupId,
    payload.ServiceId,
    response => {
      // if (response.status === 200) commit('groupDeleted', group)
    }
  );
};

export const listUngroupedApplications = ({ commit }) => {
  // TODO
  serviceDiscoveryAPI.listUngroupedApps(response => {
    if (response.status === 200) commit("ungroupedApplications", response.body);
  });
};

export const findPossibleMatches = ({ commit }, sg) => {
  // TODO
  serviceDiscoveryAPI.findPossibleMatches(sg.Name, response => {
    if (response.status === 200) commit("possibleMatches", response.body);
  });
};
