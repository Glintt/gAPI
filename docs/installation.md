## Configuration

Before installing, some configurations are required on the API. Inside the _api/_ folder, there is a folder called _configs-example/_. This folder contains an example of all configurations required.

Copy all configuration files to a new directory called _config/_. This new directory must be placed on the root of the api project.

Configuration files explanation:

1. **_oauth.json_** - oauth server configuration (domain, port, endpoint)
2. **_services.json_** - All microservices registered on the api management. (this can be replaced with mongodb)
3. **_gAPI.json_**
   - **Authentication** - user and password configuration to access admin area;
     - TokenExpirationTime - time for the token to expire. Min value = 30s
     - TokenSigningKey - Signing key for the token. Min length = 10
     - LDAP - Ldap related configuration
       - Active - If LDAP is enabled on gAPI or not
       - Domain - Ldap server
       - Port - Ldap Port (default: 389)
   - **Logs** - activate/deactivate logs and type of logging (available: _Elastic_ or _Rabbit_)
   - **CORS** - AllowedOrigins (array with all allowed origins) and AllowCredentials
   - **ServiceDiscovery** - configuration storage type (available: _file_ and _mongo_) - note that when it is being used mongo, you MUST specify MONGO_HOST and MONGO_DB environment variables.
   - **Urls** - Base urls for some of the services available on gAPI api:
     - SERVICE_DISCOVERY_GROUP - service discovery base url
     - ANALYTICS_GROUP - analytics base url
   - **Healthcheck** - Healthcheck configuration
     - Active - boolean to activate or deactivate healthcheck monitoring
     - Frequency - Frequency in seconds at which monitor is done
     - Notification - Enable or disable notifications when service goes down.
   - **Notifications**
     - Type - Notification type. Available options:
       - Slack
     - Slack - Slack notifications configuration.
       - WebhookUrl - URL to POST notifications to
   - **Protocol**
     - Https - boolean to active or deactivate HTTPS
     - CertificateFile - certificate file location
     - CertificateKey - certificate key file location
   - **Plugins**
     - Location - where plugins are stored
     - BeforeRequest - list of BeforeRequestPlugin type plugins
4. **_users.json_** - Users registered on the system to access dashboard

## Installation

gAPI is composed by six parts:

1. [gAPI Server](#gapi-server "gAPI Server")
   1. Environment Variables
   2. Run
2. [gAPI Dashboard](#gapi-dashboard "gAPI Dashboard")
   1. Environment Variables
   2. Run
3. [gAPI rabbit listener](#gapi-rabbit-listener "gAPI rabbit listener") - only required when using RabbitMQ for queueing logs storage
   1. Environment Variables
   2. Run
4. Elasticsearch - logs storage
5. RabbitMQ - used as queue for logs (_optional_)
6. MongoDB - used as service discovery storage engine (_optional_)

gAPI also can be run using docker:

1. [Docker](#docker "gAPI Docker")
   1. Environment Variables

### gAPI Server

##### Environment Variables

1. Specify gAPI Server port:

```
API_MANAGEMENT_PORT=<new port>   (default: 8080)
```

2. Enable live analytics:

```
SOCKET_PORT=<socket port>
```

3. Elasticsearch is required for logging requests:

```
ELASTICSEARCH_HOST=<elastic host>
ELASTICSEARCH_PORT=<elastic port>
```

4. To use RabbitMQ as queueing system for logging:

```
RABBITMQ_HOST=<rabbit host>
RABBITMQ_PORT=<rabbit port>     (default: 5601)
RABBITMQ_USER=<rabbit user>
RABBITMQ_PASSWORD=<rabbit password>
RABBITMQ_QUEUE=<rabbit gapi queue name>
```

5. Use MongoDB as service discovery storage engine:

```
MONGO_HOST=<mongodb host>
MONGO_DB=<mongodb database name>
```

6. Service discovery is a separate service:

```
SERVICEDISCOVERY_HOST=<custom SD host>
SERVICEDISCOVERY_PORT=<custom SD port>
```

##### Run

To run gAPI Server, follow this steps:

1. Copy the project to _go/src_ folder
2. Compile the code using the command:

```
go build -o server ./server.go
```

3. Start the server using the following command:

```
./server
```

### gAPI Rabbit Listener

##### Environment Variables

gAPI Rabbit Listener requires the following environement variables:

```
RABBITMQ_HOST=<rabbit host>
RABBITMQ_PORT=<rabbit port>     (default: 5601)
RABBITMQ_USER=<rabbit user>
RABBITMQ_PASSWORD=<rabbit password>
RABBITMQ_QUEUE=<rabbit gapi queue name>
ELASTICSEARCH_HOST=<elastic host>
ELASTICSEARCH_PORT=<elastic port>
```

These environment variables must go along with the ones specified on gAPI server.

##### Run

1. Copy the project to _go/src_ folder
2. Compile the code using the command:

```
go build -o rabbit-listener ./rabbit-listener.go
```

3. Start the listener using the following command:

```
./rabbit-listener
```

### gAPI Dashboard

All commands regarding the dashboard, must be run inside _dashboard/_ folder.

##### Environment Variables

Some environment variables are required to build the dashboard. Env vars are located in _.env.{ENV_NAME}_ files.

These are the required env vars:

```
API_HOST=<gAPI host>
API_PORT=<gAPI port>
SOCKET_HOST=<gAPI socket host>
SOCKET_PORT=<gAPI socket port>
```

##### Run

To start the dashboard, follow this steps:

1. Install all dependencies:

```
npm install
```

2. Build the sources:

```
npm run build
```

3. Start the service

```
node index.js
```

## Docker

gAPI can also be run using [docker compose](https://docs.docker.com/compose/).

To run all gAPI dependencies, just run the following command on the root of the project:

```
docker-compose up -d
```

##### Environment Variables

When using docker, all environment variables have default values already which allow to start all services without any configuration.
If you want to customize it, you can override the following environment variables:

```
- API_MANAGEMENT_PORT=${API_MANAGEMENT_PORT:-8080}
- RABBITMQ_HOST=${RABBITMQ_HOST:-rabbit}
- RABBITMQ_PORT=${RABBITMQ_PORT:-5672}
- RABBITMQ_USER=${RABBITMQ_USER:-gapi}
- RABBITMQ_PASSWORD=${RABBITMQ_PASSWORD:-gapi}
- RABBITMQ_QUEUE=${RABBITMQ_QUEUE:-gAPI-logqueue}
- ELASTICSEARCH_HOST=${ELASTICSEARCH_HOST:-elastic}
- ELASTICSEARCH_PORT=${ELASTICSEARCH_URL:-9200}
- SERVICEDISCOVERY_HOST=${SERVICEDISCOVERY_HOST:-localhost}
- SERVICEDISCOVERY_PORT=${SERVICEDISCOVERY_PORT:-8080}
- MONGO_HOST=${MONGO_HOST:-mongodb}
- MONGO_DB=${MONGO_DB:-gapi}
- SOCKET_PORT=${SOCKET_PORT:-5000}
```
