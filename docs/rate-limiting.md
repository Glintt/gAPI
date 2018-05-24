## Rate Limiting

gAPI allows to add rate limiting to micro services.

It has two options:

1. Global rate limiting - affects all registered micro-services
2. Service specific rate limiting - only affects a specific micro-service


### Configuration

In order to activate rate limiting add the following configuration to gAPI.json configuration:

```json
"RateLimiting":{
    "Active": true,
    "Limit": 20,
    "Period": 1,
    "Metrics": ["RemoteAddr", "MatchingUri"]
}
```	

> **Active** - enable or disable rate limiting; <br />
> **Limit** - max number of requests;<br />
> **Period** - time until number requests available are reset (in minutes);<br />
> **Metrics** - Metrics to use to limit requests rate (allowed values: *"RemoteAddr"*, *"MatchingUri"*)

### Service specific configuration

Each service has the following fields:

1. RateLimit - max number of request
2. RateLimitExpirationTime - time until number requests is reset (in minutes)

By default, these fields are stored with value 0. When value is 0, custom rate limiting is disabled.

In order to add rate limiting to a specific service you need to override these values.
