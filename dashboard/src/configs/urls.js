export const config = {
    "API" : {
        "HOST" :  "localhost",
        "PORT" : "8080",
        "SOCKET_HOST" : process.env.SOCKET_HOST ||  "localhost",
        "SOCKET_PORT" : process.env.SOCKET_PORT || "5000",
        "BASE_PATH" : "",
        "SERVICE_DISCOVERY_BASEPATH" : "/service-discovery",
        "ANALYTICS_BASEPATH" : "/analytics"
    }
}