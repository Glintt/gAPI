export const usersListUpdated = (state, obj) => {
  state.users = obj;
};

export const changeUser = (state, user) => {
  state.user = user;
};

export const updateUser = (state, user) => {
  state.user = user;
};

export const closeAlert = state => {
  state.alert.message = "";
  state.alert.classType = "";
  state.alert.showing = false;
};

export const newAlert = (state, alert) => {
  state.alert.message = alert.msg;
  state.alert.classType = alert.classType;
  state.alert.showing = true;
};
