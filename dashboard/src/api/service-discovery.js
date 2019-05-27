const HTTP = require("@/api/http");
const APIConfig = require("@/configs/urls").config.API;
const ServiceDiscoveryBaseURL = APIConfig.SERVICE_DISCOVERY_BASEPATH;

const Endpoints = {
  list: "/services",
  get: "/endpoint",
  store: "/services",
  delete: "/services/<service_id>",
  manage: "/services/manage",
  manage_types: "/services/manage/types",
  update: "/services/<service_id>",
  add_to_group: "/service-groups/<group_id>/services",
  deassociate_from_group: "/service-groups/<group_id>/services/<service_id>",
  list_groups: "/service-groups",
  store_group: "/service-groups",
  update_group: "/service-groups/<group_id>",
  remove_group: "/service-groups/<group_id>",

  store_application_group: "/apps-groups",
  get_application_group_by_id: "/apps-groups/<group_id>",
  list_application_group: "/apps-groups?page=-1",
  update_application_group: "/apps-groups/<group_id>",
  delete_application_group: "/apps-groups/<group_id>",
  associate_to_application_group: "/apps-groups/<group_id>/<service_id>",
  deassociate_from_application_group: "/apps-groups/<group_id>/<service_id>",
  application_group_for_service: "/apps-groups/search/<service_id>",
  list_ungrouped_applications: "/apps-groups/ungrouped",
  find_possible_group_matches: "/apps-groups/matches?group_name=<group_name>"
};

export const CustomManagementActions = ["logs"];

export function listServices(page, searchQuery, cb) {
  HTTP.GET(
    HTTP.PathToCall(
      ServiceDiscoveryBaseURL +
        Endpoints.list +
        "?page=" +
        page +
        "&q=" +
        searchQuery
    ),
    {}
  ).then(
    response => {
      cb(response);
    },
    response => {
      HTTP.handleError(response, cb);
    }
  );
}

export function getServices(serviceEndpoint, cb) {
  HTTP.GET(HTTP.PathToCall(ServiceDiscoveryBaseURL + Endpoints.get), {
    params: {
      uri: serviceEndpoint
    }
  }).then(
    response => {
      cb(response);
    },
    response => {
      HTTP.handleError(response, cb);
    }
  );
}

export function storeService(service, cb) {
  HTTP.POST(
    HTTP.PathToCall(ServiceDiscoveryBaseURL + Endpoints.store),
    service,
    {}
  ).then(
    response => {
      cb(response);
    },
    response => {
      HTTP.handleError(response, cb);
    }
  );
}

export function storeServiceGroup(group, cb) {
  HTTP.POST(
    HTTP.PathToCall(ServiceDiscoveryBaseURL + Endpoints.store_group),
    group,
    {}
  ).then(
    response => {
      cb(response);
    },
    response => {
      HTTP.handleError(response, cb);
    }
  );
}

export function addServiceToServiceGroup(groupId, serviceId, cb) {
  let obj = {
    service_id: serviceId
  };
  HTTP.POST(
    HTTP.PathToCall(
      ServiceDiscoveryBaseURL +
        Endpoints.add_to_group.replace("<group_id>", groupId)
    ),
    obj,
    {}
  ).then(
    response => {
      cb(response);
    },
    response => {
      HTTP.handleError(response, cb);
    }
  );
}

export function deassociateServiceFromServiceGroup(groupId, serviceId, cb) {
  HTTP.DELETE(
    HTTP.PathToCall(
      ServiceDiscoveryBaseURL +
        Endpoints.deassociate_from_group
          .replace("<group_id>", groupId)
          .replace("<service_id>", serviceId)
    ),
    {}
  ).then(
    response => {
      cb(response);
    },
    response => {
      HTTP.handleError(response, cb);
    }
  );
}

export function listServiceGroups(cb) {
  HTTP.GET(
    HTTP.PathToCall(ServiceDiscoveryBaseURL + Endpoints.list_groups),
    {}
  ).then(
    response => {
      cb(response);
    },
    response => {
      HTTP.handleError(response, cb);
    }
  );
}

export function updateServiceGroup(group, cb) {
  HTTP.PUT(
    HTTP.PathToCall(
      ServiceDiscoveryBaseURL +
        Endpoints.update_group.replace("<group_id>", group.Id)
    ),
    group,
    {}
  ).then(
    response => {
      cb(response);
    },
    response => {
      HTTP.handleError(response, cb);
    }
  );
}
export function deleteServiceGroup(groupId, cb) {
  HTTP.DELETE(
    HTTP.PathToCall(
      ServiceDiscoveryBaseURL +
        Endpoints.remove_group.replace("<group_id>", groupId)
    ),
    {}
  ).then(
    response => {
      cb(response);
    },
    response => {
      HTTP.handleError(response, cb);
    }
  );
}

export function deleteService(serviceId, cb) {
  HTTP.DELETE(
    HTTP.PathToCall(
      ServiceDiscoveryBaseURL +
        Endpoints.delete.replace("<service_id>", serviceId)
    ),
    {}
  ).then(
    response => {
      cb(response);
    },
    response => {
      HTTP.handleError(response, cb);
    }
  );
}

export function updateService(service, cb) {
  HTTP.PUT(
    HTTP.PathToCall(
      ServiceDiscoveryBaseURL +
        Endpoints.update.replace("<service_id>", service.Id)
    ),
    service,
    {}
  ).then(
    response => {
      cb(response);
    },
    response => {
      HTTP.handleError(response, cb);
    }
  );
}

export function manageService(service, action, cb) {
  HTTP.POST(
    HTTP.PathToCall(
      ServiceDiscoveryBaseURL +
        Endpoints.manage +
        "?service=" +
        service +
        "&action=" +
        action
    ),
    {},
    {}
  ).then(
    response => {
      cb(response);
    },
    response => {
      HTTP.handleError(response, cb);
    }
  );
}
export function manageServiceTypes(cb) {
  HTTP.GET(
    HTTP.PathToCall(ServiceDiscoveryBaseURL + Endpoints.manage_types),
    {},
    {}
  ).then(
    response => {
      cb(response);
    },
    response => {
      HTTP.handleError(response, cb);
    }
  );
}

export function storeApplicationGroup(group, cb) {
  HTTP.POST(HTTP.PathToCall(Endpoints.store_application_group), group, {}).then(
    response => {
      cb(response);
    },
    response => {
      HTTP.handleError(response, cb);
    }
  );
}

export function listAppsGroups(cb) {
  HTTP.GET(HTTP.PathToCall(Endpoints.list_application_group), {}).then(
    response => {
      cb(response);
    },
    response => {
      HTTP.handleError(response, cb);
    }
  );
}

export function updateAppsGroup(group, cb) {
  HTTP.PUT(
    HTTP.PathToCall(
      Endpoints.update_application_group.replace("<group_id>", group.Id)
    ),
    group,
    {}
  ).then(
    response => {
      cb(response);
    },
    response => {
      HTTP.handleError(response, cb);
    }
  );
}
export function deleteAppsGroup(groupId, cb) {
  HTTP.DELETE(
    HTTP.PathToCall(
      Endpoints.delete_application_group.replace("<group_id>", groupId)
    ),
    {}
  ).then(
    response => {
      cb(response);
    },
    response => {
      HTTP.handleError(response, cb);
    }
  );
}

export function addServiceToAppsGroup(groupId, serviceId, cb) {
  let obj = {
    service_id: serviceId
  };
  HTTP.POST(
    HTTP.PathToCall(
      Endpoints.associate_to_application_group
        .replace("<group_id>", groupId)
        .replace("<service_id>", serviceId)
    ),
    obj,
    {}
  ).then(
    response => {
      cb(response);
    },
    response => {
      HTTP.handleError(response, cb);
    }
  );
}

export function deassociateServiceFromAppsGroup(groupId, serviceId, cb) {
  HTTP.DELETE(
    HTTP.PathToCall(
      Endpoints.deassociate_from_application_group
        .replace("<group_id>", groupId)
        .replace("<service_id>", serviceId)
    ),
    {}
  ).then(
    response => {
      cb(response);
    },
    response => {
      HTTP.handleError(response, cb);
    }
  );
}

export function AppsGroupForService(serviceId, cb) {
  HTTP.GET(
    HTTP.PathToCall(
      Endpoints.application_group_for_service.replace("<service_id>", serviceId)
    ),
    {}
  ).then(
    response => {
      cb(response);
    },
    response => {
      HTTP.handleError(response, cb);
    }
  );
}

export function listUngroupedApps(cb) {
  HTTP.GET(HTTP.PathToCall(Endpoints.list_ungrouped_applications), {}).then(
    response => {
      cb(response);
    },
    response => {
      HTTP.handleError(response, cb);
    }
  );
}

export function findPossibleMatches(name, cb) {
  HTTP.GET(
    HTTP.PathToCall(Endpoints.find_possible_group_matches).replace(
      "<group_name>",
      name
    ),
    {}
  ).then(
    response => {
      cb(response);
    },
    response => {
      HTTP.handleError(response, cb);
    }
  );
}

export function applicationGroupById(groupId, cb) {
  HTTP.GET(
    HTTP.PathToCall(
      Endpoints.get_application_group_by_id.replace("<group_id>", groupId)
    ),
    {}
  ).then(
    response => {
      cb(response);
    },
    response => {
      HTTP.handleError(response, cb);
    }
  );
}
