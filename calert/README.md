# Calert Helm chart

calert pushes Alertmanager notifications to Google Chat via webhook integration.


To install the chart with the release name `calert`, in the namespace `clu-inf-all`, using the Google Hangouts Chat webhooks listed in the file `calert_values.yaml`:

```console
$ helm install incubator/calert --values=calert_values.yaml --name=calert --namespace=clu-inf-all
```

calert_values.yaml

```yaml
secret:
  rooms: |

    [app.chat.cluster001-alerts]
    notification_url = "https://chat.googleapis.com/v1/spaces/xxx/messages?key=abc-xyz&token=token-unique-key%3D"

    [app.chat.cluster002-alerts]
    notification_url = "https://chat.googleapis.com/v1/spaces/xxx/messages?key=abc-xyz&token=token-unique-key%3D"
```

## Configuration

| Parameter                                   | Description                                                      | Default                         |
|:--------------------------------------------|:-----------------------------------------------------------------|:--------------------------------|
| `image.repository`                          | Docker image repository                                          | mrkaran/calert                  |
| `image.tag`                                 | Docker image tag                                                 | 1.0.0-stable                    |
| `image.pullPolicy`                          | Docker image pull policy                                         | Always                          |
| `replicaCount`                              | Number of pod replicas                                           | 1                               |
| `secret.server.address`                     | Port that the app listens to in the pod                          | ":6000"                         |
| `secret.server.socket`                      | Socket that the app listens to in the pod                        | "/tmp/calert.sock"              |
| `secret.server.name`                        | Name for the server instance                                     | "calert"                        |
| `secret.server.read_timeout`                | Read timeout in milliseconds                                     | "8000"                          |
| `secret.server.write_timeout`               | Write timeout in milliseconds                                    | "8000"                          |
| `secret.server.keepalive_timeout`           | Keepalive timeout in milliseconds                                | "300000"                        |
| `secret.app.template_file`                  | Application template file                                        | "message.tmpl"                  |
| `secret.app.http_client.max_idle_conns`     | Client max idele connections                                     | "100"                           |
| `secret.app.http_client.request_timeout`    | Client request timeout in milliseconds                           | "8000"                          |
| `secret.rooms`                              | List of webhooks to send to. See `calert_values.yaml` above      | [app.chat.alertManagerTestRoom] |
| `service.type`                              | Service type                                                     | ClusterIP                       |
| `service.port`                              | Should be same as `secret.server.address` but without the `:` | 6000                            |


## Alertmanager configuration
NOTE: Currently this chart has only been tested with the `prometheus-operator` chart, so that is the only configuration described below:

The following changes should be made in the Alertmanager section of the `prometheus-operator` values file:

```yaml
alertmanager:
  config:
    route:
      routes:
      - match:
          severity: critical
        receiver: google-chat
        group_by: [alertname]
    receivers:
    - name: 'google-chat'
      webhook_configs:
      - url: "http://calert.clu-inf-all.svc.cluster.local:6000/create?room_name=<room>"
```
