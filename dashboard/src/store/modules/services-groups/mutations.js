export const updateGroups = (state, obj) => {
  state.groups = obj;
};

export const updateGroup = (state, obj) => {
  state.groups.forEach(element => {
    if (element.Id === obj.Id) {
      element = obj;
    }
  });
};

export const groupDeleted = (state, g) => {
  let newGroups = state.groups;
  console.log(newGroups);

  newGroups = newGroups.filter(obj => {
    return obj.Id !== g.Id;
  });
  console.log(newGroups);

  state.groups = newGroups;
};
