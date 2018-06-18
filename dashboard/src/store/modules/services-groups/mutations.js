
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
