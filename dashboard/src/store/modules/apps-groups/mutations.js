
export const updateGroups = (state, obj) => {
    state.groups = obj;
};

export const updateGroup = (state, obj) => {
    state.groups.forEach(element => {
        if (element.Id == obj.Id) {
            element = obj
        }
    });
};

export const groupDeleted = (state, g) => {
    let newGroups = state.groups
    
    newGroups = newGroups.filter(obj => {
        return obj.Id !== g.Id;
    })
    
    state.groups = newGroups
}