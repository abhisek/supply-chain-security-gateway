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
