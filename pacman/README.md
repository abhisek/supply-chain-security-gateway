# PacMan
Tooling to configure build tools to use security gateway

`pacman` aka. `Package Manager` inspired by the `pacman` is a tool for easily configuring various package managers such as Gradle, Maven etc. to use the security gateway.

## Setup

```bash
export GATEWAY_URL="https://<Your-Gateway-Base-URL>"
export GATEWAY_USERNAME="<Your-Gateway-Username>"
export GATEWAY_PASSWORD="<Your-Gateway-Password>"
```

> Refer to [gateway authentication](../docs/Gateway-Authentication.md) for more details on how to create gateway users.

### Configure Gradle

```bash
./pacman.sh setup-gradle
```

## Cleanup

Remove any configuration file added by `pacman`

```bash
./pacman clean
```

# Reference

* https://www.google.com/logos/2010/pacman10-i.html
