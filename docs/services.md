## Services

A service contains the following information:

* _Id_ - unique identifier
* _Name_ - Service name
* _Hosts_ - List of hosts in which the service is hosted. Multiple hosts are allowed for load balancing
* _Domain_ - Default domain
* _Port_ - Default port
* _MatchingURI_ - Unique uri to represent the service
* _MatchingURIRegex_ - Regex to search by MatchingURI
* _ToURI_ - Service base uri
* _Protected_ - If service is protected with OAuth (Boolean)
* _IsCachingActive_ - If caching is enabled (Boolean)
* _IsActive_ - Is service running (Boolean)
* _APIDocumentation_ - API Documentation location
* _HealthcheckUrl_ - URL to call to check if service is running
* _LastActiveTime_ - Last time the service was running (in milliseconds)
* _ServiceManagementHost_ - Service management service host
* _ServiceManagementPort_ - Service management service port
* _ServiceManagementEndpoints_ - Service management service available endpoints
* _RateLimit_ - Global rate limit (rps)
* _RateLimitExpirationTime_ - Global rate limit expiration time (in minutes)
* _IsReachable_ - If micro service is reachable from external requests
* _GroupId_ - Group Id to which the micro service belongs. null if no group
* _GroupVisibility_ - Group Visibility to external requests
* _UseGroupAttributes_ - If the service must use group attributes (currently: visibility)

In order to get more information on ServiceManagement part, check [here](./service-management.md)

In order to get more information on Rate Limiting part, check [here](./rate-limiting.md)

In order to get more information on Reachability, check [here](./reachability.md)

#### Matching URI

The matching uri is used to call the correct service. Currently, it is supported complex URIs with more than one subroute.

**Example:**

_/gapi/modules/service-discovery_ calls gapi service discovery component located at _/service-discovery_. Everything that goes after _/service-discovery_ is considered a service endpoint.

With this configuration, when you call _/gapi/modules/service-discovery/services_ it is called _/service-discovery/services_ on the service host.

#### Hosts

Currently, gAPI has a simple load balancing.

In order to register multiple hosts for a service, you must add an entry with a new host to **Hosts** array.
