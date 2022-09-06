# PacMan
Utility to configure build tools to use security gateway as package repository.

`pacman` aka. `Package Manager` inspired by the `pacman` is a tool for easily configuring various package managers such as Gradle, Maven etc. to use the security gateway for downloading required dependencies.

## Setup

Run `pacman` configuration wizard

```bash
./pacman.sh configure
```

> Refer to [gateway authentication]([../README.md#authentication)) for more details on how to create gateway users.

### Configure Gradle

```bash
./pacman.sh setup-gradle
```

### Configure Maven

```bash
./pacman.sh setup-maven
```

> **Note:** This script overwrite `$HOME/.m2/settings.xml`

## Cleanup

Remove any configuration file added by `pacman`

```bash
./pacman clean
```

# Reference

* https://www.google.com/logos/2010/pacman10-i.html
