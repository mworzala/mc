# mc-cli
todo write more here


## Command Spec
`<>` indicates required arg
`[]` indicates optional arg
otherwise it is a literal string

### Account management

- `mc account default|use` - Show the current default account
- `mc account default|use <username|uuid>` - Set the default account
- `mc account login` - Login to a new account with Microsoft
- ! `mc account list` - Show a list of all known accounts

### Java management

- `mc java default|use` - Show the current default java installation to be used
- `mc java default|use <name>` - Set the default installation
- `mc java list` - List all discovered java installations
- `mc java discover <path-to-exe> [--set-default]` - Discover a preinstalled java version
- !! `mc java install <version> [name]` - Install a new java dedicated to Minecraft

### Profile management

- `mc install <mc-version> [name]` - Install vanilla minecraft (profile will be named as version if unspecified)
- `mc install <mc-version> [name] --fabric [--loader loader-version]` - Install fabric
- `mc run <profile>` - Launch a profile immediately

- ! `mc profile list`
- ! `mc profile delete <name>`

### Mod management (todo)
All commands would have a `--profile`, `-p` arg to specify which profile is being edited.

- `mc mod install`

### Config management
maybe

## Configuration

All subject to change, none of it is implemented

```
# Autolink options describe whether to symlink those resources between different instances
# todo is this really viable?
autolink.config = bool
autolink.servers = bool
autolink.saves = bool
autolink.resourcepacks = bool
```

## Building from source

`go build -o mc cmd/mc/*.go`
