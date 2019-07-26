# Healthcheck

Healthcheck is a module to check if a service is up and running.

This module can be enabled on the configuration file. When a service status changes, a notifications can be sent to a channel.
Currently, the only supported channel to send notifications is Slack.

### Enable/Disable

On _gAPI.json_ configuration, there is the following section:

```
"Healthcheck": {
    "Active": true,
    "Frequency": 30,
    "Notification": false
}
```

Info:

1. Active - boolean to activate/deactivate healthcheck
2. Frequency - healthcheck frequency in seconds
3. Notification - boolean to activate/deactivate status change notifications

### Notifications

Notifications can also be enabled and disabled. In order to enable notifications, healthcheck must be enabled too.
To enable notifications, there is a section on _gAPI.json_ configuration file:

```
"Notifications": {
    "Type": "Slack",
    "Slack": {
      "WebhookUrl": "https://hooks.slack.com/services/asld/lak/la"
    }
}
```

Info:

1. Type - type of notifications to be sent (available: _Slack_)
2. Slack - for notification type _Slack_
   1. WebhookUrl - slack notifications are sent by calling an webhook. Here we can configure the webhook url.
