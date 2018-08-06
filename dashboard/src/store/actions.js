import * as api from '../api'
import * as types from './mutation-types'

export const getAll = ({ commit }) => {
	api.listServices(services => {
		commit(types.RECEIVE_ALL, {
			services
		})
	})
}

export const sendMessage = ({ commit }, payload) => {
	api.createMessage(payload, message => {
		commit(types.RECEIVE_MESSAGE, {
			message
		})
	})
}

export const switchThread = ({ commit }, payload) => {
	commit(types.SWITCH_THREAD, payload)
}
