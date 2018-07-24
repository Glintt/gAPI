## HTTPS


### API

In order to activate HTTPS, you only need to add the configuration to gAPI configuration file (located at *configs/gAPI.json*) and an environment variable.

The configuration for the HTTPS is as follows:

```
"Protocol": {
    "CertificateFile": "certificates/certificate.crt",
    "CertificateKey": "certificates/privatekey.key"
}
```

* CertificateFile - certificate file location (*required*)
* CertificateKey - certificate key file location (*required*)

The environment variable is called *API_PROTOCOL* and can assume these two values: *http* or *https*

### Dashboard


You also need to define the protocol for the dashboard. If you enable HTTPS on the server side, you need to configure the dashboard to comunicate over https with the API.

In order to enable HTTPS on the frontend, you need to specify the following env vars:

1. *FRONTEND_PROTOCOL* - the protocol of the dashboard (*https* or *http*)
1. *FRONTEND_CERT_PRIVATE_KEY* - https certificate key file location
1. *FRONTEND_CERT_FILE* - https certificate file location
1. *API_PROTOCOL* - the same as the API environment variable. it is used to communicate with the API


