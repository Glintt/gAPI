# [gAPIManagement](https://glintt.github.io/gAPIManagement/)


gAPIManagement allows to easily manage all your microservices by providing a simple and intuitive dashboard.

It also provides analytics which are useful to take value from request logs.


## Requirements

1. Golang - https://golang.org/
2. Node - https://nodejs.org/en/
3. RabbitMQ - Optional
4. Elasticsearch - Optional

### Configuration

Before installing, some configurations are required on the API. Inside the *api/* folder, there is a folder called *configs-example/*. This folder contains an example of all configurations required.

Copy all configuration files to a new directory called *config/*. This new directory must be placed on the root of the api project.


Configuration files explanation:

1. ***oauth.json*** - oauth server configuration (domain, port, endpoint)
2. ***services.json*** - All microservices registered on the api management. (this can be replaced with mongodb)
3. ***urls.json*** - base path for service-discovery and analytics api endpoints
4. ***gAPI.json*** - general api configurations (Authentication; Request Logs) 

## Installation

#### Using docker

Run the following command:

```
docker-compose up -d
```

#### API

In order to run the API, you need Golang installed.

Inside *api/* folder, run the following command:

```
sh start.sh
```


#### Dashboard

In order to start the dashboard, run the following command inside *dashboard/* folder:

```
npm run serve
```


