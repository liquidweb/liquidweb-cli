# lw (liquidweb-cli)
Official command line interface for the LiquidWeb API
```
CLI interface for LiquidWeb.

Command line interface for interacting with LiquidWeb services via
LiquidWeb's Public API.

If this is your first time running, you will need to setup at least
one auth context. An auth context contains authentication data for
accessing your LiquidWeb account. As such one auth context represents
one LiquidWeb account. You can have multiple auth contexts defined.

To setup your first auth context, you run 'auth init'. For further
information on auth contexts, be sure to checkout 'help auth' for a
list of capabilities.

As always, consult the various subcommands for specific features and
capabilities.

Usage:
  lw [command]

Available Commands:
  auth        authentication actions
  cloud       Interact with LiquidWeb's Cloud platform.
  help        Help about any command
  network     network actions
  version     show build information

Flags:
      --config string   config file (default is $HOME/.liquidweb-cli.yaml)
  -h, --help            help for lw

Use "lw [command] --help" for more information about a command.
```
## Obtaining prebuilt binaries

Head on over to the [releases page](https://github.com/liquidweb/liquidweb-cli/releases)  to get prebuilt binaries for your platform.

## Building from source

You can build lw from source by running `make build` from the root of this repository. The resulting program will be located at `./_exe/lw`.
You can also build+install lw onto your system in the ordinary `go install` way. To do this, either just run `go install` from the root of this repository,
or `make install`. If you run `make` with no arguments, this will be the default action.

## First Time Setup
The first time you use lw, you will need to setup an auth context. An auth context holds authentication related data for a specific LiquidWeb account. You can follow a guided questionnaire to add your auth contexts if you pass arguments `auth init` to lw. By default contexts are stored in `~/.liquidweb-cli.yaml` or `%APPDATA%/.liquidweb-cli.yaml` on Windows.

## Adding auth contexts later
If you end up wanting to add an auth context later on, you can do so with `auth add-context`. You can find the usage documentation in `help auth add-context`.

## Removing auth contexts later
If you end up wanting to remove an auth context later on, you can do so with `auth remove-context`. You can find the usage documentation in `help auth remove-context`.

## Modifying auth contexts later
If you end up wanting to modify an auth context later on, you can do so with `auth update-context`. You can find the usage documentation in `help auth update-context`.

## LiquidWeb Cloud
The Cloud features you can use in manage.liquidweb.com on your Cloud Servers you can do with this command line tool. See `help cloud` for a full list of features and capabilities.
