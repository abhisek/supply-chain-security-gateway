# Configuration Management

## Goal

Adopt a spec based approach for configuration model definition that can be represented in formats such as JSON, Protobuf and can be persisted as files or in databases using an appropriate repository. Ensure configuration has a <u>Single Source of Truth</u>

## Who are the users of configuration management

1. Environment administrators who provision gateways
2. Gateway administrators who configure upstream, authentication, secrets etc.
3. Envoy proxy to act as the gateway and route request to upstream
4. Microservices runtime configuration

## Configuration Model

The configuration model is based on the [domain model](images/domain.png) but with additional detail for service and environment related configuration.

## How to define configuration schema

Look at `services/spec/config.proto`

## How to use configuration schema

Use `protoc` to compile the spec and generate code as required. The underlying spec has associated validator.

## How to store the configuration

Storage and service layer for configuration is required and is within the boundary of the service implementing it. The service layer exposes the configuration to management and operations (data) plane.

## Dynamic Configuration

To be able to configure the gateway and associated service at runtime, each service must support dynamic configuration i.e. monitor for change in its configuration state and periodically re-configure itself if the SSOT for its configuration has changed.

To do so, every service must:

* Have an instance identity of its own representing itself in a deployed environment
* Query the configuration repository with its instance identifier
* Update its in-memory configuration

> **Note:** There are services that are not dynamically configurable such as queue/topic listeners.

### Bootstrap Configuration

Every service has a bootstrap configuration which is minimal and allows it to discover and access dynamic configuration source. Bootstrap configuration can be passed through environment variables:

```
BOOTSTRAP_CONFIGURATION_REPOSITORY_TYPE=file
BOOTSTRAP_CONFIGURATION_REPOSITORY_PATH=/path/to/config.yml
```

### Configuration API

The `config` common module should be initialized and subsequently can be used for obtaining currently loaded configurations:

```go
conf := config.CurrentConfiguration()
if conf.PdpService.MonitorMode() {
  // Do something
}
```
