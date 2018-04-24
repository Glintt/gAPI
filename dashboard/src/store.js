import Vue from "vue";
import Vuex from "vuex";
var auth = require('./auth')
Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    isLoggedIn : false
  },
  mutations: {
    loggedIn (state) {
      state.isLoggedIn = true
    },
    loggedOut (state) {
      state.isLoggedIn = false
    },

  },
  getters: {
    isLoggedIn: state => {
      return state.isLoggedIn
    }
  },
  actions: {
    login : ({ commit }) => {
      commit("loggedIn")
    },
    logout : ({ commit }) => {
      commit("loggedOut")
    }
  }
});
