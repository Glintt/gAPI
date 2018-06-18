
const serviceDiscoveryAPI = require('@/api/service-discovery')

export const fetchGroups = ({ commit }) => {
    serviceDiscoveryAPI.listServiceGroups((response) => {
        commit('updateGroups', response.body)
    })
}

export const updateGroup = ({commit}, group) => {
    serviceDiscoveryAPI.updateServiceGroup(group, (response) => {
        commit('updateGroup', response.body)
    })
}