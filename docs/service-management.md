
## Service Management Module

gAPI has a module that allows to integrate external service management services. 
This module is highly scalable and flexible allowing to easily add new management feature in an easy way.

Currently, gAPI configuration example contains the following features to manage each service:

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

In order to add new features, you just need to add a new entry to management types map (*ManagementTypes*). This map is located inside the configuration file (*/api/configs/gAPI.json*)

Example:

```
"feature_name" : {
		"action": "feature_name",
		"method": "POST",
		"icon": "fas fa-sync",
		"background": "info",
		"description": "feature_name service"}
```


If you need to handle the new feature in a custom way inside the dashboard, you need to add an to *CustomManagementActions* list with the new feature name.  This list is located at: */dashboard/src/api/service-discovery.js*

