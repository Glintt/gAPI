import Vue from "vue";
import Vuex from "vuex";
var auth = require('./auth')
Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    isLoggedIn : false,
    loggedInUser: null
  },
  mutations: {
    loggedIn (state) {
      state.isLoggedIn = true
    },
    loggedOut (state) {
      state.isLoggedIn = false
    },
    loggedInUserUpdate (state,user) {
      state.loggedInUser = user
    }
  },
  getters: {
    isLoggedIn: state => {
      return state.isLoggedIn
    },
    loggedInUser: state => {
      return state.loggedInUser
    }
  },
  actions: {
    loggedInUserUpdate: ({commit}, user) => {
      commit("loggedInUserUpdate", user)
    },
    login : ({ commit }) => {
      commit("loggedIn")
    },
    logout : ({ commit }) => {
      commit("loggedOut")
    }
  }
});
