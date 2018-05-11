
## Service Management Module

gAPI has features that allow to manage each service:

1.  Stop a service
2.  Restart a service
3.  Redeploy a service
4.  Backup a service version
5.  View service internal logs


### How it works

Currently, in order to manage each service, gAPI needs to interact with another services which are responsible for implementing each of the features. These services must be implemented by service manager.

Each feature is linked to an endpoint which will be called when the feature is used:

1.  Stop a service - UndeployEndpoint
2.  Restart a service - RestartEndpoint
3.  Redeploy a service - RedeployEndpoint
4.  Backup a service version - BackupEndpoint
5.  View service internal logs - LogsEndpoint

All the features share the same host and port and, as consequence, for a single service, these endpoints must be exposed on the same host and port:

1. Host - ServiceManagementHost
2. Port - ServiceManagementPort


### Implementation

#### Adding more managing features

In order to add new features, you need to follow these steps:

1. API
    1. Add a new management type to the management types map (*ManagementTypes*). This map is located inside */api/servicediscovery/service.go*
    2. Add a new field to Service struct to save the new endpoint
    3. Add new rule to *GetManagementEndpoint()*. This method will return the endpoint to call for the new action.
2. Dashboard
    1. Add a new input to the NewService form (located at: */dashboard/src/views/ServiceDiscovery/NewService.vue*)
    2. Add a new input to the EditService form (located at: */dashboard/src/views/Service/EditService.vue*)
    3. Add a new button to call the new action on Services list view (located at: */dashboard/src/views/ServiceDiscovery/ListServices.vue*)
    