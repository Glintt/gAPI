import * as actions from './actions';
import * as getters from './getters';
import * as mutations from './mutations';

const state = {
    users: [],
    user: null,
    alert: {
        message: '',
        classType: '',
        showing: false
    }
};

const namespaced = true;


export default {
    namespaced,
    state,
    actions,
    getters,
    mutations
};