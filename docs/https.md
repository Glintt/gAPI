## HTTPS

In order to activate HTTPS, you only need to add the configuration to gAPI configuration file (located at *configs/gAPI.json*).

The configuration for the HTTPS is as follows:

```
"Protocol": {
    "Https": true,
    "CertificateFile": "certificates/certificate.crt",
    "CertificateKey": "certificates/privatekey.key"
}
```

* Https - boolean to active or deactivate HTTPS
* CertificateFile - certificate file location
* CertificateKey - certificate key file location