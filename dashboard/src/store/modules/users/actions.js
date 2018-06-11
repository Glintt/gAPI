
const usersAPI = require('@/api/users')

export const updateList = ({ commit }, payload) => {
    if (payload == undefined) payload = {
        query: '',
        page: 1
    }
    usersAPI.find(payload.query, payload.page, (response) => {
        commit('usersListUpdated', response.body)
    })
}

export const changeUser = ({ commit }, user) => {
    commit('changeUser', user)
}

export const emptyUser = ({commit}) => {
    commit('changeUser', {
        Username: '',
        Password: '',
        Email: '',
        IsAdmin: false
    })
}

export const createUser = ({ commit }, user) => {
    usersAPI.create(user, (response) => {
        if (response.status != 201) {
            commit('newAlert', {msg:'Error creating user', classType: 'danger'})
        }else {
            commit('newAlert', {msg:'User created successfuly', classType: 'success'})
        }
    })
}

export const updateUser = ({ commit, state }) => {
    usersAPI.update(state.user, (response) => {
        if (response.status != 200) {
            commit('newAlert', {msg:'Error updating user', classType: 'danger'})
        }else {
            commit('newAlert', {msg:'User updated successfuly', classType: 'success'})
        }
    })
}

export const closeAlert = ({commit}) => {
    commit('closeAlert')
}