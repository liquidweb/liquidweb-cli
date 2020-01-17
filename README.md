# liquidweb-cli
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
  liquidweb-cli [command]

Available Commands:
  auth        authentication actions
  cloud       Interact with LiquidWeb's Cloud platform.
  help        Help about any command
  network     network actions

Flags:
      --config string   config file (default is $HOME/.liquidweb-cli.yaml)
  -h, --help            help for liquidweb-cli

Use "liquidweb-cli [command] --help" for more information about a command.
```
## Building from source

You can build liquidweb-cli from source by running `make build` from the root of this repository. The resulting program will be located at `./_exe/liquidweb-cli`.

## First Time Setup
The first time you use liquidweb-cli, you will need to setup an auth context. An auth context holds authentication related data for a specific LiquidWeb account. You can follow a guided questionnaire to add your auth contexts if you pass arguments `auth init` to liquidweb-cli. By default contexts are stored in `~/.liquidweb-cli.yaml` or `%APPDATA%/.liquidweb-cli.yaml` on Windows.

## Adding auth contexts later
If you end up wanting to add a auth context later on, you can do so with `auth add-context`. You can find the usage documentation in `help auth add-context`.

## Removing auth contexts later
If you end up wanting to remove a auth context later on, you can do so with `auth remove-context`. You can find the usage documentation in `help auth remove-context`.

## LiquidWeb Cloud
The Cloud features you can use in manage.liquidweb.com on your Cloud Servers you can do with this command line tool. See `help cloud` for a full list of features and capabilities.
