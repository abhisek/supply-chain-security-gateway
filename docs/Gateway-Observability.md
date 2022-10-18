# Gateway Observability

## Goal

Adopt OpenTelemetry for MELT export from gateway services. Keep collectors and APM tools decoupled from Gateway.

## Service Configuration for MELT

Following environment variables can be used to configure MELT for each microservice

| Environment Name                | Purpose                                |
| ------------------------------- | -------------------------------------- |
| APP_SERVICE_OBS_ENABLED         | True / False                           |
| APP_SERVICE_NAME                | The service name to include traces     |
| APP_SERVICE_ENV                 | The service environment name           |
| APP_SERVICE_LABELS              | Command separate key-value pairs       |
| APP_OTEL_EXPORTER_OTLP_ENDPOINT | OTLP exporter GRPC endpoint (insecure) |
