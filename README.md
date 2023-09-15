# Minecraft CLI (mc)
A simple command line interface for managing and launching Minecraft instances, friendly to both users and automation/tooling.

(insert fancy image here)

## Features
todo

Planned:
- Mod management/auto update
- Modpack installation
- MultiMC/Prism integration (either autodetecting installations and using them, or importing instances)
- Vanilla launcher integration (either autodetecting installations and using them, or importing instances)
- Java installation based on requested version in manifest
- Support for legacy launcher metadata formats (eg the ability to launch older Minecraft versions)
- Forge support
- Automatic synchronization of saves/resource packs/configs/servers between instances

## Installation

> ![WARNING]
> Windows is not officially supported yet, but it may work. If you encounter any issues, please report them as an issue.

A prebuilt binary is available for macOS and Linux in the [releases](https://github.com/mworzala/mc/releases) tab.

### Go Install
It is possible to install the latest version of `mc` using `go install`:
```shell
go install github.com/mworzala/mc@latest
```

### Building From Source
It is possible to build `mc` from source as long as you have Make and Go installed.

```shell
make build
```

## Usage
todo

## Automation

## Contributing
todo

## License
This project is licensed under the [MIT](../LICENSE).
