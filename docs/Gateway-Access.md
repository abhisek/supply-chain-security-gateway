# Gateway Access

The security gateway can be accessed as any HTTP based package repository such as `maven2`, `npm`, `rubygems` etc. However package managers can be configured to send additional metadata to enrich the generated events for auditing and traceability purpose.

## Request Metadata

Every artifact download request can have additional metadata for auditing or traceability purpose. Following metadata can be attached to a request to the gateway

1. Project ID
2. Project Environment Name
3. Labels (Generic key value pairs)

Since different package managers have different capabilities of *decorating a request*, we have to support different channels through which additional metadata can be included in a request to the gateway.

### Using Headers

| Name                 | Description                                              |
| -------------------- | -------------------------------------------------------- |
| X-SGW-Project-Id     | Name of the project using this artifact                  |
| X-SGW-Project-Env    | Environment name for which the project is built          |
| X-SGW-Project-Labels | Additional metadata in `key1=value1, key2=value2` format |
