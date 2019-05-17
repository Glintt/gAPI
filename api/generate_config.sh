#!/bin/sh
for ARGUMENT in "$@"
do
  KEY=$(echo $ARGUMENT | cut -f1 -d=)
  VALUE=$(echo $ARGUMENT | cut -f2 -d=)
  case "$KEY" in
          token_signing_key)              token_signing_key=${VALUE} ;;
          token_expiration_time)    token_expiration_time=${VALUE} ;;     
          logs_type)    logs_type=${VALUE} ;;     
          db)    db=${VALUE} ;;     
          matching_regex)    matching_regex=${VALUE} ;;     
          healthcheck_freq)    healthcheck_freq=${VALUE} ;;     
      *)   
  esac
done

mkdir -p configs

if [ -z "$token_expiration_time" ]; then
    token_expiration_time=86400000
fi

if [ -z "$token_signing_key" ]; then
    token_signing_key="GAPI_SIGN_KEY"
fi

if [ -z "$matching_regex" ]; then
    matching_regex="((/([\\\\\\w?\\\\\\-=.&+#])*)*$)"
fi

if [ -z "$db" ]; then
    db="mongo"
fi

if [ -z "$logs_type" ]; then
    logs_type="Elastic"
fi

if [ -z "$healthcheck_freq" ]; then
    healthcheck_freq=60
fi

rm configs/gAPI.json

echo "
{
  \"Authentication\": {
    \"TokenExpirationTime\": $token_expiration_time,
    \"TokenSigningKey\": \"$token_signing_key\",
    \"LDAP\": {
      \"Active\": false,
      \"Domain\": \"ldap.example.com\",
      \"Port\": \"389\"
    }
  },
  \"Logs\": {
    \"Active\": true,
    \"Type\": \"$logs_type\"
  },
  \"CORS\": {
    \"AllowedOrigins\": [\"http://localhost:8080\"],
    \"AllowCredentials\": true
  },
  \"ServiceDiscovery\": {
    \"Type\": \"$db\"
  },
  \"MatchingUriRegex\": \"$matching_regex\",
  \"Healthcheck\": {
    \"Active\": false,
    \"Frequency\": $healthcheck_freq,
    \"Notification\": true
  },
  \"Notifications\": {
    \"Type\": \"Slack\",
    \"Slack\": {
      \"WebhookUrl\": \"https://hooks.slack.com/services/asld/lak/la\"
    }
  },
  \"RateLimiting\": {
    \"Active\": false,
    \"Limit\": 20,
    \"Period\": 1,
    \"Metrics\": [\"RemoteAddr\", \"MatchingUri\"]
  },
  \"ManagementTypes\": {
      
        \"logs\" : {
            \"action\": \"logs\",
            \"method\": \"GET\",
            \"icon\": \"fas fa-file\",
            \"background\": \"\",
            \"description\": \"View service logs\"
        }

  },
  \"Protocol\": {
    \"Https\": false,
    \"CertificateFile\": \"certificates/certificate.crt\",
    \"CertificateKey\": \"certificates/privatekey.key\"
  },
  \"ThirdPartyOAuth\": {
    \"Host\": \"http://localhost\",
    \"Port\": \"8084\",
    \"AuthorizeEndpoint\": \"/api/experience/oauth/authorize\",
    \"UserTokenInformation\": {
      \"Active\": false,
      \"Source\": \"header\",
      \"Name\": \"CallData\"
    }
  }
}" >> configs/gAPI.json