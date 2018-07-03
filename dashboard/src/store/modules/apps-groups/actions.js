
const serviceDiscoveryAPI = require('@/api/service-discovery')

export const fetchGroups = ({ commit }) => {
    serviceDiscoveryAPI.listAppsGroups((response) => {
        commit('updateGroups', response.body)
    })
}

export const updateGroup = ({commit}, group) => {
    serviceDiscoveryAPI.updateAppsGroup(group, (response) => {
        commit('updateGroup', response.body)
    })
}

export const deleteGroup = ({commit}, group) => {
    serviceDiscoveryAPI.deleteAppsGroup(group.Id, response => {
        if (response.status == 200) commit('groupDeleted', group)
    })
}